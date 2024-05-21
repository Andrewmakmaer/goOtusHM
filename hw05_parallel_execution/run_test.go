package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("Try run without sleep", func(t *testing.T) {
		tasksCount := 500
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			var result error
			if i%10 == 0 {
				result = fmt.Errorf("error from task %d", i)
			} else {
				result = nil
			}
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return result
			})
		}

		workersCount := 20
		maxErrorsCount := 51
		var countGorutines int
		var wg sync.WaitGroup
		var err error
		wg.Add(2)
		go func() {
			err = Run(tasks, workersCount, maxErrorsCount)
			wg.Done()
		}()
		go func() {
			time.Sleep(time.Millisecond * time.Duration(100))
			countGorutines = runtime.NumGoroutine() // надо как-то сделать чтобы показатель снимался во время рана основного
			wg.Done()
		}()
		wg.Wait()

		require.Equal(t, workersCount, countGorutines-5, "Max gorutines count")
		// countGorutines-5 т.к. есть горутины запускаемые не только в Run функции
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.Falsef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
	})
}
