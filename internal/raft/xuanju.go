package raft

import (
	"log"
	"raft/internal/rpcs"
	"raft/pkg/types"
)

type XuanJuManner struct {
	NodeID    int
	State     *State
	LogManner *LogManner
	Client    rpcs.RPCClient
	Friends   []int
}

func NewXuanJuManner(nodeID int, state *State, manner *LogManner, client rpcs.RPCClient, friends []int) *XuanJuManner {
	return &XuanJuManner{
		NodeID:    nodeID,
		State:     state,
		LogManner: manner,
		Client:    client,
		Friends:   friends,
	}
}

// 启动选举流程
func (xm *XuanJuManner) XuanJu() {
	xm.State.SetState(types.Candidate)       //先变候选人
	xm.State.SetTerm(xm.State.GetTerm() + 1) //递增任期
	xm.State.SetVote(xm.NodeID)              //投自己一票

	//检查日志一致性和任期
	LogIndex, Term := xm.State.GetLogAndTerm()
	voteRequest := types.VoteRequest{
		Term:         xm.State.GetTerm(),
		CandidateId:  xm.NodeID,
		LastLogIndex: LogIndex,
		LastLogTerm:  Term,
	}
	//票数
	votes := 1
	voteChan := make(chan bool, len(xm.Friends)) //收集投票结果

	for _, fri := range xm.Friends {
		go func(fri int) {
			OK, err := xm.Client.RequestVote(fri, voteRequest)
			if err != nil {
				log.Printf("fail to get vote from node:%d.%v", fri, err)
				voteChan <- false
				return
			}
			voteChan <- OK
		}(fri)
	}

	//统计投票结果
	go func() {
		for i := 0; i < len(xm.Friends); i++ {
			if <-voteChan {
				votes++
			}
			//上任鹅城！
			if votes > len(xm.Friends)/2 && xm.State.GetState() == types.Candidate {
				xm.State.SetState(types.Leader)
				return
			}
		}
	}()
}

// 处理投票请求
func (xm *XuanJuManner) Handle(req types.VoteRequest) types.VoteResponse {
	res := types.VoteResponse{
		Term:        xm.State.GetTerm(),
		VoteGranted: false,
	}

	if req.Term < xm.State.GetTerm() {
		return res
	}

	//老弟不行啊
	if req.Term > xm.State.GetTerm() {
		xm.State.SetTerm(req.Term)
		xm.State.SetState(types.Follower)
		xm.State.SetVote(-1)
	}
	log := xm.LogManner.Check(req.LastLogIndex, req.LastLogTerm)

	//判断能投吗
	//话说这个狗屎自动补全怎么老补出来的和我选的不一样
	flag := xm.State.GetVote() == -1 || xm.State.GetVote() == req.CandidateId

	if log && flag {
		res.VoteGranted = true
		res.Term = req.Term
		xm.State.SetVote(req.CandidateId)
	}
	return res
}
