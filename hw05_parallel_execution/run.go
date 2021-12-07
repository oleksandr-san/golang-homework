package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var wg sync.WaitGroup
	wg.Add(n)
	defer wg.Wait()

	tasksCh := make(chan Task)
	defer close(tasksCh)

	var errorCount int32

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for task := range tasksCh {
				if err := task(); err != nil {
					atomic.AddInt32(&errorCount, 1)
				}
			}
		}()
	}

	for _, task := range tasks {
		if int(atomic.LoadInt32(&errorCount)) >= m {
			return ErrErrorsLimitExceeded
		}

		tasksCh <- task
	}

	return nil
}
