package rpcs

import (
	"net/rpc"
	"raft/pkg/types"
)

type RPCClient interface {
	RequestVote(friendID int, req types.VoteRequest) (bool, error)
	AppendEntries(friendID int, req types.AppendEntriesRequest) (*types.AppendEntriesResponse, error)
}

type RaftClient struct {
	FriendAddr map[int]string
}

func NewRaftClient(addr map[int]string) *RaftClient {
	return &RaftClient{
		FriendAddr: addr,
	}
}

func (c *RaftClient) RequestVote(friendID int, req types.VoteRequest) (bool, error) {
	client, err := rpc.DialHTTP("tcp", c.FriendAddr[friendID])
	if err != nil {
		return false, err
	}
	defer client.Close() //神经啊这也要处理错误处理毛线,莫名其妙这提示

	var res types.VoteResponse
	err = client.Call("RaftNode.RequestVote", req, &res)
	if err != nil {
		return false, err
	}
	return res.VoteGranted, nil

}

func (c *RaftClient) AppendEntries(friendID int, req types.AppendEntriesRequest) (*types.AppendEntriesResponse, error) {
	client, err := rpc.DialHTTP("tcp", c.FriendAddr[friendID])
	if err != nil {
		return nil, err
	}
	defer client.Close()

	var res types.AppendEntriesResponse
	err = client.Call("RaftNode.AppendEntries", req, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
