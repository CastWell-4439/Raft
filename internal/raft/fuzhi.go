package raft

import (
	"log"
	"raft/internal/rpcs"
	"raft/pkg/types"
	"sort"
	"time"
)

type FuZhi struct {
	Node    int
	State   *State
	Log     *LogManner
	Client  rpcs.RPCClient
	Friends []int
	Stop    chan struct{}
}

// 重置日志相关信息
func (fz *FuZhi) reset() {
	logs := fz.State.GetLog()
	index := len(logs)

	for _, fri := range fz.Friends {
		fz.State.NextIndex[fri] = index //和leader相同
		fz.State.MatchIndex[fri] = -1   //未同步
	}
}

func NewFuZhi(node int, state *State, log *LogManner, client rpcs.RPCClient, friend []int) *FuZhi {
	fz := &FuZhi{
		Node:    node,
		State:   state,
		Log:     log,
		Client:  client,
		Friends: friend,
		Stop:    make(chan struct{}),
	}
	fz.reset()
	return fz
}

func (fz *FuZhi) sendHeartBeat(fri int) {
	req := types.AppendEntriesRequest{
		Term:         fz.State.GetTerm(),
		LeaderId:     fz.Node,
		PreLogIndex:  -1,
		PreLogTerm:   -1,
		Entries:      []types.LogEntry{},
		LeaderCommit: fz.State.GetCommit(),
	}
	res, err := fz.Client.AppendEntries(fri, req)
	if err != nil {
		log.Printf("fail to sen heartbeat,%v", err)
		return
	}

	if !res.OK && res.Term > fz.State.GetTerm() {
		fz.State.SetTerm(res.Term)
		fz.State.SetState(types.Follower)
		fz.State.SetVote(-1)
	}
}

func (fz *FuZhi) sendHeartBeats() {
	for _, fri := range fz.Friends {
		go fz.sendHeartBeat(fri)
	}
}

func (fz *FuZhi) StartHeartBeat(times time.Duration) {
	go func() {
		heart := time.NewTicker(times)
		defer heart.Stop()

		for {
			select {
			case <-fz.Stop:
				return
			case <-heart.C:
				if fz.State.GetState() == types.Leader {
					fz.sendHeartBeats()
				}
			}
		}
	}()
}

// 根据多数的进度更新提交索引
func (fz *FuZhi) UpdateIndex() {
	match := make([]int, 0, len(fz.Friends)+1)
	match = append(match, len(fz.State.GetLog())-1)

	for _, fri := range fz.Friends {
		if idx, ok := fz.State.NextIndex[fri]; ok {
			match = append(match, idx) //前index个已复制惹
		}
	}
	sort.Ints(match)
	commitIndex := match[len(match)/2] //占到多数派就行啦

	//只能提交索引大于当前，且任期正确的
	logs := fz.State.GetLog()
	if commitIndex > fz.State.GetTerm() && commitIndex < len(logs) && logs[commitIndex].Term == fz.State.GetTerm() {
		fz.State.SetCommit(commitIndex)
	}
}

// 向跟随着发请求
func (fz *FuZhi) fuzhiToFriend(fri int) {
	index := fz.State.NextIndex[fri]
	preLogIndex := index - 1
	preLogTerm := -1

	if preLogIndex >= 0 {
		if entry, ok := fz.Log.GetLogEntry(preLogIndex); ok {
			preLogTerm = entry.Term
		}
	}

	//获取要复制的
	entries := fz.Log.GetLogSlice(index, len(fz.State.GetLog()))
	req := types.AppendEntriesRequest{
		Term:         fz.State.GetTerm(),
		LeaderId:     fz.Node,
		PreLogIndex:  preLogIndex,
		PreLogTerm:   preLogTerm,
		Entries:      entries,
		LeaderCommit: fz.State.GetCommit(),
	}

	res, err := fz.Client.AppendEntries(fri, req)
	if err != nil {
		log.Printf("fail to add log to:%d ,%v", fri, err)
		return
	}

	if res.OK {
		fz.State.NextIndex[fri] = index + len(entries)
		fz.State.MatchIndex[fri] = fz.State.NextIndex[fri] - 1
		fz.UpdateIndex()
	} else {
		if res.Term > fz.State.GetTerm() {
			fz.State.SetTerm(res.Term)
			fz.State.SetState(types.Follower)
		} else {
			fz.State.NextIndex[fri]--
			go fz.fuzhiToFriend(fri)
		}
	}

}

// goroutine太伟大了我说
func (fz *FuZhi) startFuzhi() {
	for _, fri := range fz.Friends {
		go fz.fuzhiToFriend(fri)
	}
}

// follower处理请求
func (fz *FuZhi) HandleAppend(req types.AppendEntriesRequest) types.AppendEntriesResponse {
	res := types.AppendEntriesResponse{
		Term: fz.State.GetTerm(),
		OK:   false,
	}
	if req.Term < fz.State.GetTerm() {
		return res
	}

	if req.Term > fz.State.GetTerm() {
		fz.State.SetTerm(req.Term)
		fz.State.SetState(types.Follower)
		fz.State.SetVote(-1)
	}

	if req.PreLogIndex >= 0 {
		logs := fz.State.GetLog()
		if req.PreLogIndex >= len(logs) && logs[req.PreLogIndex].Term != req.PreLogTerm {
			return res
		}
	}

	if len(req.Entries) > 0 {
		fz.Log.AppendEntries(req.PreLogIndex, req.PreLogTerm, req.Entries)
	}

	if req.LeaderCommit > fz.State.GetCommit() {
		logs := fz.State.GetLog()
		commitIndex := req.LeaderCommit
		if commitIndex >= len(logs) {
			commitIndex = len(logs) - 1
		}
		fz.State.SetCommit(commitIndex)
	}
	res.OK = true
	return res

}
