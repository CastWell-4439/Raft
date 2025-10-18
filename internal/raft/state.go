package raft

import (
	"raft/pkg/types"
	"sync"
)

// 感觉都很无聊啊，没什么好注释的，看看得了
type State struct {
	Lock       sync.Mutex
	State      types.State
	LeaderID   int
	NextIndex  map[int]int //leader状态
	MatchIndex map[int]int
	Term       int
	Vote       int
	Commit     int
	LastApply  int
	Log        []types.LogEntry
}

func NewState() *State {
	return &State{
		State:      types.Follower,
		LeaderID:   -1,
		NextIndex:  make(map[int]int),
		MatchIndex: make(map[int]int),
		Commit:     -1,
		LastApply:  -1,
		Log:        make([]types.LogEntry, 0),
	}
}

func (s *State) GetTerm() int {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	return s.Term
}

func (s *State) SetTerm(term int) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.Term = term
	s.Vote = -1
}

func (s *State) GetVote() int {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	return s.Vote
}

func (s *State) SetVote(candidateID int) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.Vote = candidateID
}

func (s *State) GetState() types.State {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	return s.State
}

func (s *State) SetState(state types.State) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.State = state
}

func (s *State) GetCommit() int {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	return s.Commit
}

func (s *State) SetCommit(index int) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.Commit = index
}

func (s *State) GetLog() []types.LogEntry {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	return s.Log
}

func (s *State) GetLogAndTerm() (int, int) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	if len(s.Log) == 0 {
		return -1, -1
	}

	Log := s.Log[len(s.Log)-1]
	return len(s.Log) - 1, Log.Term
}

func (s *State) AppendLog(entries []types.LogEntry) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.Log = append(s.Log, entries...)
}
