// Code generated by protoc-gen-go. DO NOT EDIT.
// source: xraft.proto

/*
Package xraftpb is a generated protocol buffer package.

It is generated from these files:
	xraft.proto

It has these top-level messages:
	RequestVoteArgs
	RequestVoteReply
	LogEntry
	AppendEntriesArgs
	AppendEntriesReply
	InstallSnapshotArgs
*/
package xraftpb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type RequestVoteArgs struct {
	Term         int32 `protobuf:"varint,1,opt,name=term" json:"term,omitempty"`
	CandidateId  int32 `protobuf:"varint,2,opt,name=candidateId" json:"candidateId,omitempty"`
	LastLogTerm  int32 `protobuf:"varint,3,opt,name=lastLogTerm" json:"lastLogTerm,omitempty"`
	LastLogIndex int32 `protobuf:"varint,4,opt,name=lastLogIndex" json:"lastLogIndex,omitempty"`
}

func (m *RequestVoteArgs) Reset()                    { *m = RequestVoteArgs{} }
func (m *RequestVoteArgs) String() string            { return proto.CompactTextString(m) }
func (*RequestVoteArgs) ProtoMessage()               {}
func (*RequestVoteArgs) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *RequestVoteArgs) GetTerm() int32 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *RequestVoteArgs) GetCandidateId() int32 {
	if m != nil {
		return m.CandidateId
	}
	return 0
}

func (m *RequestVoteArgs) GetLastLogTerm() int32 {
	if m != nil {
		return m.LastLogTerm
	}
	return 0
}

func (m *RequestVoteArgs) GetLastLogIndex() int32 {
	if m != nil {
		return m.LastLogIndex
	}
	return 0
}

type RequestVoteReply struct {
	Term        int32 `protobuf:"varint,1,opt,name=term" json:"term,omitempty"`
	VoteGranted bool  `protobuf:"varint,2,opt,name=voteGranted" json:"voteGranted,omitempty"`
}

func (m *RequestVoteReply) Reset()                    { *m = RequestVoteReply{} }
func (m *RequestVoteReply) String() string            { return proto.CompactTextString(m) }
func (*RequestVoteReply) ProtoMessage()               {}
func (*RequestVoteReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *RequestVoteReply) GetTerm() int32 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *RequestVoteReply) GetVoteGranted() bool {
	if m != nil {
		return m.VoteGranted
	}
	return false
}

type LogEntry struct {
	LogIndex   int32  `protobuf:"varint,1,opt,name=logIndex" json:"logIndex,omitempty"`
	LogTerm    int32  `protobuf:"varint,2,opt,name=logTerm" json:"logTerm,omitempty"`
	LogCommand string `protobuf:"bytes,3,opt,name=logCommand" json:"logCommand,omitempty"`
}

func (m *LogEntry) Reset()                    { *m = LogEntry{} }
func (m *LogEntry) String() string            { return proto.CompactTextString(m) }
func (*LogEntry) ProtoMessage()               {}
func (*LogEntry) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *LogEntry) GetLogIndex() int32 {
	if m != nil {
		return m.LogIndex
	}
	return 0
}

func (m *LogEntry) GetLogTerm() int32 {
	if m != nil {
		return m.LogTerm
	}
	return 0
}

func (m *LogEntry) GetLogCommand() string {
	if m != nil {
		return m.LogCommand
	}
	return ""
}

type AppendEntriesArgs struct {
	Term         int32       `protobuf:"varint,1,opt,name=term" json:"term,omitempty"`
	LeaderId     int32       `protobuf:"varint,2,opt,name=leaderId" json:"leaderId,omitempty"`
	PrevLogTerm  int32       `protobuf:"varint,3,opt,name=prevLogTerm" json:"prevLogTerm,omitempty"`
	PrevLogIndex int32       `protobuf:"varint,4,opt,name=prevLogIndex" json:"prevLogIndex,omitempty"`
	Entries      []*LogEntry `protobuf:"bytes,5,rep,name=entries" json:"entries,omitempty"`
	LeaderCommit int32       `protobuf:"varint,6,opt,name=leaderCommit" json:"leaderCommit,omitempty"`
}

func (m *AppendEntriesArgs) Reset()                    { *m = AppendEntriesArgs{} }
func (m *AppendEntriesArgs) String() string            { return proto.CompactTextString(m) }
func (*AppendEntriesArgs) ProtoMessage()               {}
func (*AppendEntriesArgs) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *AppendEntriesArgs) GetTerm() int32 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *AppendEntriesArgs) GetLeaderId() int32 {
	if m != nil {
		return m.LeaderId
	}
	return 0
}

func (m *AppendEntriesArgs) GetPrevLogTerm() int32 {
	if m != nil {
		return m.PrevLogTerm
	}
	return 0
}

func (m *AppendEntriesArgs) GetPrevLogIndex() int32 {
	if m != nil {
		return m.PrevLogIndex
	}
	return 0
}

func (m *AppendEntriesArgs) GetEntries() []*LogEntry {
	if m != nil {
		return m.Entries
	}
	return nil
}

func (m *AppendEntriesArgs) GetLeaderCommit() int32 {
	if m != nil {
		return m.LeaderCommit
	}
	return 0
}

type AppendEntriesReply struct {
	Term      int32 `protobuf:"varint,1,opt,name=term" json:"term,omitempty"`
	NextIndex int32 `protobuf:"varint,2,opt,name=nextIndex" json:"nextIndex,omitempty"`
	Success   bool  `protobuf:"varint,3,opt,name=success" json:"success,omitempty"`
}

func (m *AppendEntriesReply) Reset()                    { *m = AppendEntriesReply{} }
func (m *AppendEntriesReply) String() string            { return proto.CompactTextString(m) }
func (*AppendEntriesReply) ProtoMessage()               {}
func (*AppendEntriesReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *AppendEntriesReply) GetTerm() int32 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *AppendEntriesReply) GetNextIndex() int32 {
	if m != nil {
		return m.NextIndex
	}
	return 0
}

func (m *AppendEntriesReply) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

type InstallSnapshotArgs struct {
	Term              int32  `protobuf:"varint,1,opt,name=term" json:"term,omitempty"`
	LeaderId          int32  `protobuf:"varint,2,opt,name=leaderId" json:"leaderId,omitempty"`
	LastIncludedIndex int32  `protobuf:"varint,3,opt,name=lastIncludedIndex" json:"lastIncludedIndex,omitempty"`
	LastIncludedTerm  int32  `protobuf:"varint,4,opt,name=lastIncludedTerm" json:"lastIncludedTerm,omitempty"`
	Data              []byte `protobuf:"bytes,5,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *InstallSnapshotArgs) Reset()                    { *m = InstallSnapshotArgs{} }
func (m *InstallSnapshotArgs) String() string            { return proto.CompactTextString(m) }
func (*InstallSnapshotArgs) ProtoMessage()               {}
func (*InstallSnapshotArgs) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *InstallSnapshotArgs) GetTerm() int32 {
	if m != nil {
		return m.Term
	}
	return 0
}

func (m *InstallSnapshotArgs) GetLeaderId() int32 {
	if m != nil {
		return m.LeaderId
	}
	return 0
}

func (m *InstallSnapshotArgs) GetLastIncludedIndex() int32 {
	if m != nil {
		return m.LastIncludedIndex
	}
	return 0
}

func (m *InstallSnapshotArgs) GetLastIncludedTerm() int32 {
	if m != nil {
		return m.LastIncludedTerm
	}
	return 0
}

func (m *InstallSnapshotArgs) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*RequestVoteArgs)(nil), "xraftpb.RequestVoteArgs")
	proto.RegisterType((*RequestVoteReply)(nil), "xraftpb.RequestVoteReply")
	proto.RegisterType((*LogEntry)(nil), "xraftpb.LogEntry")
	proto.RegisterType((*AppendEntriesArgs)(nil), "xraftpb.AppendEntriesArgs")
	proto.RegisterType((*AppendEntriesReply)(nil), "xraftpb.AppendEntriesReply")
	proto.RegisterType((*InstallSnapshotArgs)(nil), "xraftpb.installSnapshotArgs")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for XRaft service

type XRaftClient interface {
	RequestVote(ctx context.Context, in *RequestVoteArgs, opts ...grpc.CallOption) (*RequestVoteReply, error)
	AppendEntries(ctx context.Context, in *AppendEntriesArgs, opts ...grpc.CallOption) (*AppendEntriesReply, error)
}

type xRaftClient struct {
	cc *grpc.ClientConn
}

func NewXRaftClient(cc *grpc.ClientConn) XRaftClient {
	return &xRaftClient{cc}
}

func (c *xRaftClient) RequestVote(ctx context.Context, in *RequestVoteArgs, opts ...grpc.CallOption) (*RequestVoteReply, error) {
	out := new(RequestVoteReply)
	err := grpc.Invoke(ctx, "/xraftpb.XRaft/RequestVote", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *xRaftClient) AppendEntries(ctx context.Context, in *AppendEntriesArgs, opts ...grpc.CallOption) (*AppendEntriesReply, error) {
	out := new(AppendEntriesReply)
	err := grpc.Invoke(ctx, "/xraftpb.XRaft/AppendEntries", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for XRaft service

type XRaftServer interface {
	RequestVote(context.Context, *RequestVoteArgs) (*RequestVoteReply, error)
	AppendEntries(context.Context, *AppendEntriesArgs) (*AppendEntriesReply, error)
}

func RegisterXRaftServer(s *grpc.Server, srv XRaftServer) {
	s.RegisterService(&_XRaft_serviceDesc, srv)
}

func _XRaft_RequestVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestVoteArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(XRaftServer).RequestVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/xraftpb.XRaft/RequestVote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(XRaftServer).RequestVote(ctx, req.(*RequestVoteArgs))
	}
	return interceptor(ctx, in, info, handler)
}

func _XRaft_AppendEntries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppendEntriesArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(XRaftServer).AppendEntries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/xraftpb.XRaft/AppendEntries",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(XRaftServer).AppendEntries(ctx, req.(*AppendEntriesArgs))
	}
	return interceptor(ctx, in, info, handler)
}

var _XRaft_serviceDesc = grpc.ServiceDesc{
	ServiceName: "xraftpb.XRaft",
	HandlerType: (*XRaftServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RequestVote",
			Handler:    _XRaft_RequestVote_Handler,
		},
		{
			MethodName: "AppendEntries",
			Handler:    _XRaft_AppendEntries_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "xraft.proto",
}

func init() { proto.RegisterFile("xraft.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 442 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x93, 0xc1, 0x6e, 0x13, 0x31,
	0x10, 0x86, 0x59, 0xda, 0x34, 0xe9, 0xa4, 0x88, 0xc6, 0x5c, 0x96, 0x80, 0x50, 0xe4, 0x53, 0x04,
	0x28, 0x87, 0xf2, 0x04, 0x15, 0x20, 0x08, 0xea, 0xc9, 0x20, 0xc4, 0xb1, 0x6e, 0x3c, 0x0d, 0x91,
	0x1c, 0x7b, 0xb1, 0x27, 0x55, 0xfa, 0x10, 0x3c, 0x00, 0xcf, 0xc1, 0xdb, 0xf0, 0x34, 0xc8, 0xf6,
	0xee, 0xe2, 0x6d, 0x9a, 0x0b, 0xb7, 0x9d, 0x7f, 0x66, 0x3d, 0xff, 0x7c, 0x63, 0xc3, 0x70, 0xeb,
	0xe4, 0x35, 0xcd, 0x2a, 0x67, 0xc9, 0xb2, 0x7e, 0x0c, 0xaa, 0x2b, 0xfe, 0xb3, 0x80, 0xc7, 0x02,
	0x7f, 0x6c, 0xd0, 0xd3, 0x57, 0x4b, 0x78, 0xee, 0x96, 0x9e, 0x31, 0x38, 0x24, 0x74, 0xeb, 0xb2,
	0x98, 0x14, 0xd3, 0x9e, 0x88, 0xdf, 0x6c, 0x02, 0xc3, 0x85, 0x34, 0x6a, 0xa5, 0x24, 0xe1, 0x5c,
	0x95, 0x0f, 0x63, 0x2a, 0x97, 0x42, 0x85, 0x96, 0x9e, 0x2e, 0xec, 0xf2, 0x4b, 0xf8, 0xf9, 0x20,
	0x55, 0x64, 0x12, 0xe3, 0x70, 0x52, 0x87, 0x73, 0xa3, 0x70, 0x5b, 0x1e, 0xc6, 0x92, 0x8e, 0xc6,
	0x3f, 0xc2, 0x69, 0x66, 0x47, 0x60, 0xa5, 0x6f, 0xf7, 0xf9, 0xb9, 0xb1, 0x84, 0x1f, 0x9c, 0x34,
	0x84, 0xc9, 0xcf, 0x40, 0xe4, 0x12, 0xbf, 0x84, 0xc1, 0x85, 0x5d, 0xbe, 0x37, 0xe4, 0x6e, 0xd9,
	0x18, 0x06, 0xba, 0xe9, 0x9a, 0x4e, 0x69, 0x63, 0x56, 0x42, 0x5f, 0xd7, 0x9e, 0xd3, 0x54, 0x4d,
	0xc8, 0x5e, 0x00, 0x68, 0xbb, 0x7c, 0x6b, 0xd7, 0x6b, 0x69, 0x54, 0x1c, 0xe8, 0x58, 0x64, 0x0a,
	0xff, 0x53, 0xc0, 0xe8, 0xbc, 0xaa, 0xd0, 0xa8, 0xd0, 0x65, 0x85, 0x7e, 0x2f, 0xbd, 0xd0, 0x1f,
	0xa5, 0x42, 0xd7, 0xa2, 0x6b, 0xe3, 0x30, 0x49, 0xe5, 0xf0, 0xe6, 0x0e, 0xb7, 0x4c, 0x0a, 0xdc,
	0xea, 0xb0, 0xc3, 0x2d, 0xd7, 0xd8, 0x2b, 0xe8, 0x63, 0x32, 0x51, 0xf6, 0x26, 0x07, 0xd3, 0xe1,
	0xd9, 0x68, 0x56, 0xaf, 0x78, 0xd6, 0x50, 0x10, 0x4d, 0x45, 0x5c, 0x44, 0x6c, 0x1f, 0x26, 0x59,
	0x51, 0x79, 0x54, 0x2f, 0x22, 0xd3, 0xf8, 0x25, 0xb0, 0xce, 0x6c, 0xfb, 0x57, 0xf1, 0x1c, 0x8e,
	0x0d, 0x6e, 0x29, 0x79, 0x4b, 0xd3, 0xfd, 0x13, 0x02, 0x5e, 0xbf, 0x59, 0x2c, 0xd0, 0xfb, 0x38,
	0xda, 0x40, 0x34, 0x21, 0xff, 0x5d, 0xc0, 0x93, 0x95, 0xf1, 0x24, 0xb5, 0xfe, 0x6c, 0x64, 0xe5,
	0xbf, 0x5b, 0xfa, 0x2f, 0x80, 0xaf, 0x61, 0x14, 0xae, 0xd0, 0xdc, 0x2c, 0xf4, 0x46, 0xa1, 0x4a,
	0x3e, 0x12, 0xc6, 0xdd, 0x04, 0x7b, 0x09, 0xa7, 0xb9, 0x18, 0x99, 0x27, 0xa0, 0x3b, 0x7a, 0x70,
	0xa2, 0x24, 0xc9, 0xb2, 0x37, 0x29, 0xa6, 0x27, 0x22, 0x7e, 0x9f, 0xfd, 0x2a, 0xa0, 0xf7, 0x4d,
	0xc8, 0x6b, 0x62, 0xef, 0x60, 0x98, 0x5d, 0x55, 0x56, 0xb6, 0xc0, 0xef, 0xbc, 0xa7, 0xf1, 0xd3,
	0xfb, 0x32, 0x91, 0x27, 0x7f, 0xc0, 0x3e, 0xc1, 0xa3, 0x0e, 0x67, 0x36, 0x6e, 0xab, 0x77, 0xee,
	0xd6, 0xf8, 0xd9, 0xfd, 0xb9, 0xfa, 0xac, 0xab, 0xa3, 0xf8, 0xb8, 0xdf, 0xfc, 0x0d, 0x00, 0x00,
	0xff, 0xff, 0xa3, 0x51, 0xf7, 0x5d, 0xeb, 0x03, 0x00, 0x00,
}