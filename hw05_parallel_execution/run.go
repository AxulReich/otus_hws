package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
var ErrEmptyTasks = errors.New("passed empty tasks list")
var ErrZeroWorkers = errors.New("passed invalid workers number")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}
	if len(tasks) == 0 {
		return ErrEmptyTasks
	}
	if n < 1 {
		return ErrZeroWorkers
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

	err := process(resChan, stopChan, m)
	return err
}

func process(errCh chan error, stopCh chan struct{}, maxErrNum int) error {
	var (
		failedTasksNum int
		errRes         error
	)

	for err := range errCh {
		if err != nil {
			failedTasksNum++
		}
		if failedTasksNum >= maxErrNum {
			errRes = ErrErrorsLimitExceeded
			break
		}
	}

	close(stopCh)
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
					case <-stopCh:
						return
					default:
					}

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
