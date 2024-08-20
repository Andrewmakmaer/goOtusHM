package scheduler

import (
	"context"
	"time"
)

type Scheduler struct {
	period time.Duration
	doFunc func()
}

func NewSchedule(everyTime time.Duration, doFunc func()) Scheduler {
	return Scheduler{period: everyTime, doFunc: doFunc}
}

func (s *Scheduler) Do(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				s.doFunc()
				time.Sleep(s.period)
			}
		}
	}()
}
