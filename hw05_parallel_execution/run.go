package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrNoWorkersRequested  = errors.New("no workers requested")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return ErrNoWorkersRequested
	}

	var wg sync.WaitGroup
	defer wg.Wait()

	tasksCh := make(chan Task)
	var errorCount int32

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tasksCh {
				if err := task(); err != nil {
					atomic.AddInt32(&errorCount, 1)
				}
			}
		}()
	}

	defer close(tasksCh)
	for _, task := range tasks {
		if int(atomic.LoadInt32(&errorCount)) >= m {
			return ErrErrorsLimitExceeded
		}

		tasksCh <- task
	}

	return nil
}
