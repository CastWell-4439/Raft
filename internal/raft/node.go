package raft

import (
	"raft/internal/rpcs"
	"raft/pkg/types"
	"time"
)

type Node struct {
	Config types.NodeConfig
	State  *State
	Log    *LogManner
	Timer  *TimerManner
	XuanJU *XuanJuManner
	Client rpcs.RPCClient
	FuZhi  *FuZhi
	Stop   chan struct{}
}

func NewNode(config types.NodeConfig, addr map[int]string) *Node {

	state := NewState()
	log := NewLog(state)
	client := rpcs.NewRaftClient(addr)
	node := &Node{
		Config: config,
		State:  state,
		Log:    log,
		Client: client,
		Stop:   make(chan struct{}),
		Timer: NewTimerManner(
			config.HeartBeat,
			100*time.Millisecond,
			200*time.Millisecond,
		),
		XuanJU: NewXuanJuManner(
			config.ID,
			state,
			log,
			client,
			config.Friends,
		),
		FuZhi: NewFuZhi(
			config.ID,
			state,
			log,
			client,
			config.Friends,
		),
	}
	return node
}

func (n *Node) XuanJuTimeout() {
	if n.State.GetState() != types.Leader {
		n.XuanJU.XuanJu() //妹心跳就是妹leader
	}
	n.Timer.ResetXuanJuTimer()
}

func (n *Node) HeartBeatTimeout() {
	if n.State.GetState() != types.Leader {
		n.XuanJU.XuanJu()
	}
	n.Timer.ResetHeartBeatTimer()
}

func (n *Node) run() {
	for {
		select {
		case <-n.Stop:
			return
		case <-n.Timer.SetXuanJuTimer():
			n.XuanJU.XuanJu()
		case <-n.Timer.SetHeartBeatTimer():
			n.HeartBeatTimeout()
		}
	}
}

func (n *Node) Run() {
	go n.run()
}

func (n *Node) Close() {
	close(n.Stop)
	n.Timer.XuanJuTimer.Stop()
	n.Timer.HeartBeatTimer.Stop()
}

func (n *Node) RequestVote(req types.VoteRequest, res *types.VoteResponse) error {
	*res = n.XuanJU.Handle(req)
	if req.Term >= n.State.GetTerm() || res.VoteGranted {
		n.Timer.ResetXuanJuTimer()
	}
	return nil
}

func (n *Node) AppendEntries(req types.AppendEntriesRequest, res *types.AppendEntriesResponse) error {
	*res = n.FuZhi.HandleAppend(req)
	if res.OK {
		n.Timer.ResetXuanJuTimer() //重置防止触发选举（因为已经有leader了
	}
	return nil
}

func (n *Node) Submit(info interface{}, res *string) error {
	if n.State.GetState() != types.Leader {
		*res = "without leader"
		return nil
	}
	entry := types.LogEntry{
		Term:    n.State.GetTerm(),
		Command: info,
	}
	n.State.AppendLog([]types.LogEntry{entry})
	n.FuZhi.startFuzhi()
	*res = "ok"
	return nil
}
