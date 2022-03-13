package hw05parallelexecution

import (
	"errors"
	"sync"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrEmptyTasks          = errors.New("passed empty tasks list")
	ErrZeroWorkers         = errors.New("passed invalid workers number")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, workerNum, maxErrCount int) error {
	if maxErrCount <= 0 {
		return ErrErrorsLimitExceeded
	}
	if len(tasks) == 0 {
		return ErrEmptyTasks
	}
	if workerNum < 1 {
		return ErrZeroWorkers
	}

	var (
		taskChan = make(chan Task, len(tasks))
		stopChan = make(chan struct{})
	)

	for i := range tasks {
		taskChan <- tasks[i]
	}
	close(taskChan)

	resChan := run(stopChan, taskChan, workerNum)
	return process(resChan, stopChan, maxErrCount)
}

func process(errCh chan error, stopCh chan struct{}, maxErrNum int) error {
	var (
		failedTasksNum   int
		errRes           error
		isstopChanClosed bool

		stopChanCheckerCloser = func() {
			if !isstopChanClosed {
				close(stopCh)
				isstopChanClosed = true
			}
		}
	)

	for err := range errCh {
		if err != nil {
			failedTasksNum++
		}

		if failedTasksNum >= maxErrNum {
			errRes = ErrErrorsLimitExceeded
			stopChanCheckerCloser()
		}
	}

	stopChanCheckerCloser()
	return errRes
}

func run(stopCh chan struct{}, taskCh chan Task, workerNum int) chan error {
	errCh := make(chan error, workerNum)

	go func() {
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
	}()

	return errCh
}
