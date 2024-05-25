package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var wg sync.WaitGroup
	var mute sync.RWMutex
	ch := make(chan Task, len(tasks))

	for _, task := range tasks {
		wg.Add(1)
		ch <- task
	}
	close(ch)

	var finish sync.Once
	finishBody := func() {
		for range ch {
			wg.Done()
		}
	}

	for i := 0; i < n; i++ {
		go func() {
			for task := range ch {
				err := task()
				if err != nil {
					mute.Lock()
					m--
					mute.Unlock()
				}
				wg.Done()
				mute.RLock()
				mVal := m
				mute.RUnlock()
				if mVal <= 0 {
					finish.Do(finishBody)
					return
				}
			}
		}()
	}
	wg.Wait()

	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}
