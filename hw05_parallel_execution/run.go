package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	var (
		taskChan = make(chan Task, len(tasks))
		resChan  = make(chan error)
		stopChan = make(chan struct{})
	)

	for i := range tasks {
		taskChan <- tasks[i]
	}
	close(taskChan)

	go run(resChan, stopChan, taskChan, n)

	return process(resChan, stopChan, m)
}

func process(errCh chan error, stopCh chan struct{}, maxErrNum int) error {
	var (
		failedTasksNum int
		errRes         error
	)
	defer close(stopCh)

	for err := range errCh {
		if err != nil {
			failedTasksNum++
		}
		if failedTasksNum >= maxErrNum {
			errRes = ErrErrorsLimitExceeded
			break
		}
	}
	return errRes
}

func run(errCh chan error, stopCh chan struct{}, taskCh chan Task, workerNum int) {
	var wg sync.WaitGroup
	defer func() {
		wg.Wait()
		close(errCh)
	}()

	wg.Add(workerNum)
	for j := 0; j < workerNum; j++ {

		go func() {
			defer wg.Done()

			for task := range taskCh {
			LOOP:
				for {
					select {
					case errCh <- task():
						break LOOP
					case <-stopCh:
						return
					default:
					}
				}

			}
		}()
	}
}
