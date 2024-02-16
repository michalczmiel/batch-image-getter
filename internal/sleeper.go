package internal

import (
	"math/rand"
	"time"
)

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

type Sleeper interface {
	Sleep()
}

type realSleeper struct {
	MinInterval int
	MaxInterval int
}

func (s *realSleeper) duration() int {
	if s.MinInterval == s.MaxInterval {
		return s.MinInterval
	}

	return RandomInt(s.MinInterval, s.MaxInterval)
}

func (s *realSleeper) Sleep() {
	duration := s.duration()

	if duration == 0 {
		return
	}

	time.Sleep(time.Duration(duration) * time.Second)
}

func NewSleeper(interval, maxInterval int) Sleeper {
	return &realSleeper{MinInterval: interval, MaxInterval: maxInterval}
}
