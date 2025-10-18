package types

import "time"

type State int

const (
	Follower State = iota
	Candidate
	Leader
)

// 一条日志条目包含任期号和具体命令
type LogEntry struct {
	Term    int
	Command any
}

// RPC请求和响应
type VoteRequest struct {
	Term         int
	CandidateId  int
	LastLogIndex int //最后一条日志的索引
	LastLogTerm  int //最后一条任期
}

type VoteResponse struct {
	Term        int
	VoteGranted bool //是否投票
}

// 日志追加要用的
type AppendEntriesRequest struct {
	Term         int
	LeaderId     int
	PreLogIndex  int //追加条目之前的索引，方便检查
	PreLogTerm   int
	LeaderCommit int        //已提交的日志索引
	Entries      []LogEntry //追加条目列表
}

type AppendEntriesResponse struct {
	Term int
	OK   bool
}

type NodeConfig struct {
	ID        int
	RPCAddr   string
	Friends   []int //其他节点的id列表
	HeartBeat time.Duration
	TimeOut   time.Duration
}
