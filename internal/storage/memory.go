package storage

import "raft/pkg/types"

type Memory struct {
	Term int
	Vote int
	Log  []types.LogEntry
}

func NewMemory() *Memory {
	return &Memory{
		Term: 0,
		Vote: -1,
		Log:  make([]types.LogEntry, 0),
	}
}

func (m *Memory) GetTerm() int {
	return m.Term
}

func (m *Memory) SetTerm(term int) {
	m.Term = term
}

func (m *Memory) GetVote() int {
	return m.Vote
}

func (m *Memory) SetVote(vote int) {
	m.Vote = vote
}

func (m *Memory) GetLog() []types.LogEntry {
	return append([]types.LogEntry{}, m.Log...)
}

func (m *Memory) AppendLog(e []types.LogEntry) {
	m.Log = append(m.Log, e...)
}

func (m *Memory) GetIndex() int {
	return len(m.Log) - 1
}

func (m *Memory) GetLastLogTerm() int {
	if len(m.Log) == 0 {
		return -1
	}
	return m.Log[len(m.Log)-1].Term
}
