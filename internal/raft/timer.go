package raft

import (
	"math/rand"
	"time"
)

type TimerManner struct {
	TimeOut        time.Duration
	HeartBeat      time.Duration
	XuanJuTimer    *time.Timer //两计时器
	HeartBeatTimer *time.Timer
}

func NewTimerManner(heartbeat, xuanjumin, xuanjumax time.Duration) *TimerManner {
	t := &TimerManner{
		HeartBeat: heartbeat,
		TimeOut: func() time.Duration {
			temp := int64(xuanjumax - xuanjumin)
			random := rand.Int63n(temp)
			return xuanjumin + time.Duration(random)
		}(),
	}
	t.XuanJuTimer = time.NewTimer(t.TimeOut)
	t.HeartBeatTimer = time.NewTimer(t.HeartBeat)

	return t
}

func (t *TimerManner) ResetXuanJuTimer() {
	t.XuanJuTimer.Stop()
	t.XuanJuTimer.Reset(t.TimeOut)
}

func (t *TimerManner) ResetHeartBeatTimer() {
	t.HeartBeatTimer.Stop()
	t.HeartBeatTimer.Reset(t.HeartBeat)
}

// 并发爽
func (t *TimerManner) SetHeartBeatTimer() <-chan time.Time {
	return t.HeartBeatTimer.C
}
func (t *TimerManner) SetXuanJuTimer() <-chan time.Time {
	return t.XuanJuTimer.C
}
