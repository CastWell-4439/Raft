package rpcs

import (
	"net"
	"net/rpc"
	"raft/pkg/types"
)

type RaftNode interface {
	RequestVote(req types.VoteRequest, res *types.VoteResponse) error
	AppendEntries(req types.AppendEntriesRequest, res *types.AppendEntriesResponse) error
}

// 实现一下
type RaftNodes struct {
	Node RaftNode
}

type Server struct {
	Node   RaftNode
	Server *rpc.Server
}

func NewRPCServer(node RaftNode) *Server {
	server := rpc.NewServer()
	raftNode := &RaftNodes{Node: node}
	server.RegisterName("RaftNode", raftNode)
	return &Server{
		Node:   node,
		Server: server,
	}
}

// 转发给node实例处理嘻嘻
func (w *RaftNodes) RequestVote(req types.VoteRequest, res *types.VoteResponse) error {
	return w.Node.RequestVote(req, res)
}

func (w *RaftNodes) AppendEntries(req types.AppendEntriesRequest, res *types.AppendEntriesResponse) error {
	return w.Node.AppendEntries(req, res)
}

func (s *Server) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	go s.Server.Accept(listener)
	return nil
}
