package XRaft

import(
	pb "XRaft/xraftpb"
	"time"
	"golang.org/x/net/context"
	"sync"
	"bytes"
	"encoding/gob"
	"math/rand"
	_ "fmt"
)


const(
	LEADER=iota
	CANDIDATE
	FLLOWER

	HBINTERVAL= 50*time.Millisecond
)

//
// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make().
//
type ApplyMsg struct {
	Index       int32
	Command     string
	UseSnapshot bool   // ignore for lab2; only used in lab3
	Snapshot    []byte // ignore for lab2; only used in lab3
}


type XRaft struct{
	mu        sync.Mutex
	peers     []*pb.XRaftClient
	persister *Persister
	me        int32 // index into peers[]

	// Your data here.
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.

	//channel
	state         int
	voteCount     int
	chanCommit    chan bool
	chanHeartbeat chan bool
	chanGrantVote chan bool
	chanLeader    chan bool
	chanApply     chan ApplyMsg

	//persistent state on all server
	currentTerm int32
	votedFor    int32
	log         [] *pb.LogEntry

	//volatile state on all servers
	commitIndex int32
	lastApplied int32

	//volatile state on leader
	nextIndex  []int32
	matchIndex []int32
}


// return currentTerm and whether this server
// believes it is the leader.
func (xrf *XRaft) GetState() (int32, bool) {
	return xrf.currentTerm, xrf.state == LEADER
}

func (xrf *XRaft) getLastIndex() int32 {
	return xrf.log[len(xrf.log)-1].LogIndex
}
func (xrf *XRaft) getLastTerm() int32 {
	return xrf.log[len(xrf.log)-1].LogTerm
}
func (xrf *XRaft) IsLeader() bool {
	return xrf.state == LEADER
}




// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
//
func (xrf *XRaft) persist() {
	// Your code here.
	// Example:
	w := new(bytes.Buffer)
	e := gob.NewEncoder(w)
	e.Encode(xrf.currentTerm)
	e.Encode(xrf.votedFor)
	e.Encode(xrf.log)
	data := w.Bytes()
	xrf.persister.SaveRaftState(data)
	//fmt.Println("data:",xrf.log) //
}

func (xrf *XRaft) readSnapshot(data []byte) {

	xrf.readPersist(xrf.persister.ReadRaftState())

	if len(data) == 0 {
		return
	}

	r := bytes.NewBuffer(data)
	d := gob.NewDecoder(r)

	var LastIncludedIndex int32
	var LastIncludedTerm int32

	d.Decode(&LastIncludedIndex)
	d.Decode(&LastIncludedTerm)

	xrf.commitIndex = LastIncludedIndex
	xrf.lastApplied = LastIncludedIndex

	xrf.log = truncateLog(LastIncludedIndex, LastIncludedTerm, xrf.log)

	msg := ApplyMsg{UseSnapshot: true, Snapshot: data}

	go func() {
		xrf.chanApply <- msg
	}()
}



func truncateLog(lastIncludedIndex int32, lastIncludedTerm int32, log [] *pb.LogEntry) []*pb.LogEntry {

	var newLogEntries []*pb.LogEntry
	newLogEntries = append(newLogEntries, &pb.LogEntry{LogIndex: lastIncludedIndex, LogTerm: lastIncludedTerm})

	for index := len(log) - 1; index >= 0; index-- {
		if log[index].LogIndex == lastIncludedIndex && log[index].LogTerm == lastIncludedTerm {
			newLogEntries = append(newLogEntries, log[index+1:]...)
			break
		}
	}

	return newLogEntries
}

//
// restore previously persisted state.
//
func (xrf *XRaft) readPersist(data []byte) {
	// Your code here.
	// Example:
	r := bytes.NewBuffer(data)
	d := gob.NewDecoder(r)
	d.Decode(&xrf.currentTerm)
	d.Decode(&xrf.votedFor)
	d.Decode(&xrf.log)
}


func (xrf *XRaft)RequestVote(ctx context.Context,args *pb.RequestVoteArgs) (*pb.RequestVoteReply, error){
	xrf.mu.Lock()
	defer xrf.mu.Unlock()
	defer xrf.persist()

	reply := & pb.RequestVoteReply{}
	reply.VoteGranted =false
	if args.Term < xrf.currentTerm {
		reply.Term = xrf.currentTerm
		return reply,nil
	}

	if args.Term > xrf.currentTerm {
		xrf.currentTerm = args.Term
		xrf.state = FLLOWER
		xrf.votedFor = -1
	}
	reply.Term = xrf.currentTerm

	term := xrf.getLastTerm()
	index := xrf.getLastIndex()
	uptoDate := false

	//If votedFor is null or candidateId, and candidate’s log is at least as up-to-date as receiver’s log, grant vote
	//Raft determines which of two logs is more up-to-date by comparing the index and term of the last entries in the logs.
	// If the logs have last entries with different terms,then the log with the later term is more up-to-date.
	// If the logs end with the same term, then whichever log is longer is more up-to-date
	if args.LastLogTerm > term {
		uptoDate = true
	}

	if args.LastLogTerm == term && args.LastLogIndex >= index {
		// at least up to date
		uptoDate = true
	}

	if (xrf.votedFor == -1 || xrf.votedFor == args.CandidateId) && uptoDate {
		xrf.chanGrantVote <- true
		xrf.state = FLLOWER
		reply.VoteGranted = true
		xrf.votedFor = args.CandidateId
	}
	return reply,nil
}



//
func (xrf *XRaft) sendRequestVote(server int32, args pb.RequestVoteArgs, reply *pb.RequestVoteReply) bool {
	//ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	var err error
	ok := false
	reply,err =  (* xrf.peers[server]).RequestVote(context.Background(),&args)
	if err == nil{
		ok = true
	}
	xrf.mu.Lock()
	defer xrf.mu.Unlock()
	if ok {
		term := xrf.currentTerm
		if xrf.state != CANDIDATE {
			return ok
		}
		if args.Term != term {
			return ok
		}
		if reply.Term > term {
			xrf.currentTerm = reply.Term
			xrf.state = FLLOWER
			xrf.votedFor = -1
			xrf.persist()
		}
		if reply.VoteGranted {
			xrf.voteCount++
			if xrf.state == CANDIDATE && xrf.voteCount > len(xrf.peers)/2 {
				xrf.state = FLLOWER//NO USE!
				xrf.chanLeader <- true
			}
		}
	}
	return ok
}




func (xrf *XRaft) broadcastRequestVote() {
	var args pb.RequestVoteArgs
	xrf.mu.Lock()
	args.Term = xrf.currentTerm
	args.CandidateId = xrf.me
	args.LastLogTerm = xrf.getLastTerm()
	args.LastLogIndex = xrf.getLastIndex()
	xrf.mu.Unlock()

	for i := range xrf.peers {
		if i != int(xrf.me) && xrf.state == CANDIDATE {
			go func(i int32) {
				var reply pb.RequestVoteReply
				xrf.sendRequestVote(i, args, &reply)
			}(int32(i))
		}
	}
}



//func (rf *Raft) AppendEntries(args AppendEntriesArgs, reply *AppendEntriesReply) {
//error 标记RPC是否成功返回
func (xrf *XRaft) AppendEntries(ctx context.Context, args *pb.AppendEntriesArgs) (*pb.AppendEntriesReply, error) {
	// Your code here.
	xrf.mu.Lock()
	defer xrf.mu.Unlock()
	defer xrf.persist()

	reply := &pb.AppendEntriesReply{}
	reply.Success = false
	//Reply false if term < currentTerm
	if args.Term < xrf.currentTerm {
		reply.Term = xrf.currentTerm
		reply.NextIndex = xrf.getLastIndex() + 1
		return reply, nil
		//errors.New("Error: term expired.")
	}
	xrf.chanHeartbeat <- true
	//If RPC request or response contains term T > currentTerm: set currentTerm = T, convert to follower
	if args.Term > xrf.currentTerm {
		xrf.currentTerm = args.Term
		xrf.state = FLLOWER
		xrf.votedFor = -1
	}
	reply.Term = args.Term

	if args.PrevLogIndex > xrf.getLastIndex() {
		reply.NextIndex = xrf.getLastIndex() + 1
		return reply, nil
		//errors.New("Error: missing log entries in follower")
	}

	baseIndex := xrf.log[0].LogIndex

	// If a follower’s log is inconsistent with the leader’s, the AppendEntries consistency check will fail in the next AppendEntries RPC.
	// After a rejection, the leader decrements nextIndex and retries the AppendEntries RPC
	//Eventually nextIndex will reach a point where the leader and follower logs match
	//which removes any conflicting entries in the follower’s log and appends entries from the leader’s log (if any).
	if args.PrevLogIndex > baseIndex {
		term := xrf.log[args.PrevLogIndex-baseIndex].LogTerm
		if args.PrevLogTerm != term {
			for i := args.PrevLogIndex - 1; i >= baseIndex; i-- {
				if xrf.log[i-baseIndex].LogTerm != term {
					reply.NextIndex = i + 1
					break
				}
			}
			return reply, nil
			// errors.New("Error: follower contains log entries conflict with leader's")
		}
	}
	if args.PrevLogIndex < baseIndex {

	} else {
		//Append any new entries not already in the log
		xrf.log = xrf.log[:args.PrevLogIndex+1-baseIndex]
		xrf.log = append(xrf.log, args.Entries...)
		reply.Success = true
		reply.NextIndex = xrf.getLastIndex() + 1
	}
	//If leaderCommit > commitIndex, set commitIndex =min(leaderCommit, index of last new entry)
	if args.LeaderCommit > xrf.commitIndex {
		last := xrf.getLastIndex()
		if args.LeaderCommit > last {
			xrf.commitIndex = last
		} else {
			xrf.commitIndex = args.LeaderCommit
		}
		xrf.chanCommit <- true
	}
	return reply, nil
}

func (xrf *XRaft) sendAppendEntries(server int32, args pb.AppendEntriesArgs, reply *pb.AppendEntriesReply) bool {
	//ok := rf.peers[server].Call("Raft.AppendEntries", args, reply)
	xrf.mu.Lock()
	var err error = nil
	var ok bool = false
	reply ,err = (*xrf.peers[server]).AppendEntries(context.Background(),&args)
	if err == nil{
		ok = true
	}
	defer xrf.mu.Unlock()
	if ok {
		if xrf.state != LEADER {
			return ok
		}
		//RPC调用过程中，当前节点的term可能被其他节点的RPC修改
		if args.Term != xrf.currentTerm {
			return ok
		}

		if reply.Term > xrf.currentTerm {
			xrf.currentTerm = reply.Term
			xrf.state = FLLOWER
			xrf.votedFor = -1
			xrf.persist()
			return ok
		}
		if reply.Success {
			if len(args.Entries) > 0 {
				xrf.nextIndex[server] = args.Entries[len(args.Entries)-1].LogIndex + 1
				//reply.NextIndex
				//rf.nextIndex[server] = reply.NextIndex
				xrf.matchIndex[server] = xrf.nextIndex[server] - 1
			}
		} else {
			xrf.nextIndex[server] = reply.NextIndex
		}
	}
	return ok
}


/**
 * Log replication
 */
func (xrf *XRaft) broadcastAppendEntries() {
	xrf.mu.Lock()
	defer xrf.mu.Unlock()
	N := xrf.commitIndex
	last := xrf.getLastIndex()
	baseIndex := xrf.log[0].LogIndex
	//If there exists an N such that N > commitIndex, a majority of matchIndex[i] ≥ N, and log[N].term == currentTerm: set commitIndex = N
	for i := xrf.commitIndex + 1; i <= last; i++ {
		num := 1
		for j := range xrf.peers {
			if j != int(xrf.me) && xrf.matchIndex[j] >= i && xrf.log[i-baseIndex].LogTerm == xrf.currentTerm {
				num++
			}
		}
		if 2*num > len(xrf.peers) {
			N = i
		}
	}
	if N != xrf.commitIndex {
		xrf.commitIndex = N
		xrf.chanCommit <- true
	}

	for i := range xrf.peers {
		if i != int(xrf.me) && xrf.state == LEADER {

			//copy(args.Entries, rf.log[args.PrevLogIndex + 1:])

			if xrf.nextIndex[i] > baseIndex {
				var args pb.AppendEntriesArgs
				args.Term = xrf.currentTerm
				args.LeaderId = xrf.me
				args.PrevLogIndex = xrf.nextIndex[i] - 1
				//	fmt.Printf("baseIndex:%d PrevLogIndex:%d\n",baseIndex,args.PrevLogIndex )
				args.PrevLogTerm = xrf.log[args.PrevLogIndex-baseIndex].LogTerm
				//args.Entries = make([]LogEntry, len(rf.log[args.PrevLogIndex + 1:]))
				args.Entries = make([]*pb.LogEntry, len(xrf.log[args.PrevLogIndex+1-baseIndex:]))
				copy(args.Entries, xrf.log[args.PrevLogIndex+1-baseIndex:])
				args.LeaderCommit = xrf.commitIndex
				go func(i int, args pb.AppendEntriesArgs) {
					var reply pb.AppendEntriesReply
					xrf.sendAppendEntries(int32(i), args, &reply)
				}(i, args)
			} else {
				/*var args InstallSnapshotArgs
				args.Term = rf.currentTerm
				args.LeaderId = rf.me
				args.LastIncludedIndex = rf.log[0].LogIndex
				args.LastIncludedTerm = rf.log[0].LogTerm
				args.Data = rf.persister.snapshot
				go func(server int, args InstallSnapshotArgs) {
					reply := &InstallSnapshotReply{}
					rf.sendInstallSnapshot(server, args, reply)
				}(i, args)*/
			}
		}
	}
}



/*
type InstallSnapshotArgs struct {
	Term              int
	LeaderId          int
	LastIncludedIndex int
	LastIncludedTerm  int
	Data              []byte
}

type InstallSnapshotReply struct {
	Term int
}

func (rf *XRaft) InstallSnapshot(args InstallSnapshotArgs, reply *InstallSnapshotReply) {
	// Your code here.
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		return
	}
	rf.chanHeartbeat <- true
	rf.state = FLLOWER
	rf.currentTerm = rf.currentTerm

	rf.persister.SaveSnapshot(args.Data)

	rf.log = truncateLog(args.LastIncludedIndex, args.LastIncludedTerm, rf.log)

	msg := ApplyMsg{UseSnapshot: true, Snapshot: args.Data}

	rf.lastApplied = args.LastIncludedIndex
	rf.commitIndex = args.LastIncludedIndex

	rf.persist()

	rf.chanApply <- msg
}

func (rf *XRaft) sendInstallSnapshot(server int, args InstallSnapshotArgs, reply *InstallSnapshotReply) bool {
	ok := rf.peers[server].Call("Raft.InstallSnapshot", args, reply)
	if ok {
		if reply.Term > rf.currentTerm {
			rf.currentTerm = reply.Term
			rf.state = FLLOWER
			rf.votedFor = -1
			return ok
		}

		rf.nextIndex[server] = args.LastIncludedIndex + 1
		rf.matchIndex[server] = args.LastIncludedIndex
	}
	return ok
}
*/


//
// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
//
func (xrf *XRaft) Start(command string) (int32, int32, bool) {
	xrf.mu.Lock()
	defer xrf.mu.Unlock()
	var index int32 = -1
	term := xrf.currentTerm
	isLeader := xrf.state == LEADER
	if isLeader {
		index = xrf.getLastIndex() + 1
		//fmt.Printf("raft:%d start\n",rf.me)
		xrf.log = append(xrf.log, &pb.LogEntry{LogTerm: term, LogCommand: command, LogIndex: index}) // append new entry from client
		xrf.persist()
	}
	return index, term, isLeader
}

//
// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
//
func Make(peers [] *pb.XRaftClient, me int32,  persister *Persister, applyCh chan ApplyMsg) *XRaft {
	xrf := &XRaft{}
	xrf.peers = peers
	xrf.persister = persister
	xrf.me = me


	// Your initialization code here.
	xrf.state = FLLOWER
	xrf.votedFor = -1
	xrf.log = append(xrf.log, &pb.LogEntry{LogTerm: 0})
	xrf.currentTerm = 0
	xrf.chanCommit = make(chan bool, 100)
	xrf.chanHeartbeat = make(chan bool, 100)
	xrf.chanGrantVote = make(chan bool, 100)
	xrf.chanLeader = make(chan bool, 100)
	xrf.chanApply = applyCh

	// initialize from state persisted before a crash
	xrf.readPersist(persister.ReadRaftState())
	xrf.readSnapshot(persister.ReadSnapshot())


	go func() {
		for {
			switch xrf.state {
			case FLLOWER:
				select {
				case <-xrf.chanHeartbeat:
				case <-xrf.chanGrantVote:
				case <-time.After(time.Duration(rand.Int63()%333+550) * time.Millisecond):
					xrf.state = CANDIDATE
				}
			case LEADER:
				//fmt.Printf("Leader:%v %v\n",rf.me,"boatcastAppendEntries	")
				xrf.broadcastAppendEntries()
				time.Sleep(HBINTERVAL)
			case CANDIDATE:
				xrf.mu.Lock()
				//To begin an election, a follower increments its current term and transitions to candidate state
				xrf.currentTerm++
				//It then votes for itself and issues RequestVote RPCs in parallel to each of the other servers in the cluster.
				xrf.votedFor = xrf.me
				xrf.voteCount = 1
				xrf.persist()
				xrf.mu.Unlock()
				//(a) it wins the election, (b) another server establishes itself as leader, or (c) a period of time goes by with no winner
				go xrf.broadcastRequestVote()
				select {
				case <-time.After(time.Duration(rand.Int63()%300+510) * time.Millisecond):
				case <-xrf.chanHeartbeat:
					xrf.state = FLLOWER
					//	fmt.Printf("CANDIDATE %v reveive chanHeartbeat\n",rf.me)
				case <-xrf.chanLeader:
					xrf.mu.Lock()
					xrf.state = LEADER
					//fmt.Printf("%v is Leader\n",rf.me)//
					xrf.nextIndex = make([]int32, len(xrf.peers))
					xrf.matchIndex = make([]int32, len(xrf.peers))
					for i := range xrf.peers {
						//The leader maintains a nextIndex for each follower, which is the index of the next log entry the leader will send to that follower.
						// When a leader first comes to power, it initializes all nextIndex values to the index just after the last one in its log
						xrf.nextIndex[i] = xrf.getLastIndex() + 1
						xrf.matchIndex[i] = 0
					}
					xrf.mu.Unlock()
					//rf.boatcastAppendEntries()
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-xrf.chanCommit:
				xrf.mu.Lock()
				commitIndex := xrf.commitIndex
				baseIndex := xrf.log[0].LogIndex
				for i := xrf.lastApplied + 1; i <= commitIndex; i++ {
					msg := ApplyMsg{Index: i, Command: xrf.log[i-baseIndex].LogCommand}
					applyCh <- msg
					//fmt.Printf("me:%d %v\n",rf.me,msg)
					xrf.lastApplied = i
				}
				xrf.mu.Unlock()
			}
		}
	}()


	return xrf
}
