package raft

import (
	"raft/pkg/types"
)

type LogManner struct {
	state *State
}

func NewLog(state *State) *LogManner {
	return &LogManner{state: state} //为什么这里加逗号就冗余，何意味
}

func (l *LogManner) GetLogEntry(index int) (types.LogEntry, bool) {
	log := l.state.GetLog() //先搞个完整日志
	if index < 0 || index >= len(log) {
		return types.LogEntry{}, false
	}
	return log[index], true
}

func (l *LogManner) GetLogSlice(left, right int) []types.LogEntry {
	log := l.state.GetLog()
	if left < 0 {
		left = 0
	}
	if right > len(log) {
		right = len(log)
	}
	if left > right {
		return []types.LogEntry{}
	}
	return log[left:right]
}

func (l *LogManner) AppendEntries(LogIndex int, LogTerm int, entries []types.LogEntry) bool {
	log := l.state.GetLog()
	if LogIndex >= 0 {
		if LogIndex >= len(log) {
			return false
		}
		//删掉不匹配的条目，完成同步
		if log[LogIndex].Term != LogTerm {
			//清空
			l.state.AppendLog([]types.LogEntry{})
			//搞回来到有效之前的
			l.state.AppendLog(log[:LogIndex])
			return false
		}
	}
	//新条目++
	if len(entries) > 0 {
		l.state.AppendLog(entries)
	}
	return true
}

// 检查candidate日志是否够新
func (l *LogManner) Check(LogIndex int, LogTerm int) bool {
	LastIndex, LastTerm := l.state.GetLogAndTerm()

	//首先任期大的更新，若一样更长的日志更新
	if LogIndex != LastIndex {
		return LogTerm > LastTerm
	}
	return LogIndex >= LastIndex
}
