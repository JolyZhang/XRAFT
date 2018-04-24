package XRaft

//
// support for Raft and kvraft to save persistent
// Raft state (log &c) and k/v server snapshots.
//
// we will use the original persister.go to test your code for grading.
// so, while you can modify this code to help you debug, please
// test with the original before submitting.
//

import (
	"sync"
	"os"
	"log"
	"bufio"
)

type Persister struct {
	mu        sync.Mutex
	raftstate []byte
	snapshot  []byte
	logfile  *os.File
}

func MakePersister(filepath string) *Persister {
	persister := &Persister{}
	var err error

	persister.logfile, err= os.OpenFile(filepath + "//xraft.log" ,os.O_RDWR|os.O_CREATE, 0600)
	if err != nil{
		log.Fatal("[Persist Error] could not open log file.")
	}
	return persister
}

func (ps *Persister) Copy(filepath string) *Persister {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	np := MakePersister(filepath)
	np.raftstate = ps.raftstate
	np.snapshot = ps.snapshot
	return np
}

func (ps *Persister) SaveRaftState(data []byte) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	//log.Println("save raft state.")
	writer := bufio.NewWriter(ps.logfile)
	writer.Write(data)
	writer.Flush()
	ps.raftstate = data
}

func (ps *Persister) ReadRaftState() []byte {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	//log.Println("read raft state.")
	reader := bufio.NewReader(ps.logfile)
	reader.Read(ps.raftstate)
	return ps.raftstate
}

func (ps *Persister) RaftStateSize() int {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	return len(ps.raftstate)
}

func (ps *Persister) SaveSnapshot(snapshot []byte) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.snapshot = snapshot
}

func (ps *Persister) ReadSnapshot() []byte {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	return ps.snapshot
}
