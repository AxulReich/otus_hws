package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

		workersCount := 10
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, int32(tasksCount), runTasksCount, "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})
}

func TestRun_Extra(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("pass invalid parameters", func(t *testing.T) {
		tasks := []Task{func() error { return nil }}

		for _, tc := range []struct {
			name      string
			tasks     []Task
			workerNum int
			maxErrNum int
			expErr    error
		}{
			{
				name:      "invalid maxErrNum == 0; expect ErrErrorsLimitExceeded",
				tasks:     tasks,
				workerNum: 1,
				maxErrNum: 0,
				expErr:    ErrErrorsLimitExceeded,
			},
			{
				name:      "invalid workerNum == 0; expect ErrZeroWorkers",
				tasks:     tasks,
				workerNum: 0,
				maxErrNum: len(tasks),
				expErr:    ErrZeroWorkers,
			},
			{
				name:      "invalid workerNum == -1; expect ErrZeroWorkers",
				tasks:     tasks,
				workerNum: -1,
				maxErrNum: len(tasks),
				expErr:    ErrZeroWorkers,
			},
			{
				name:      "invalid tasks with no tasks; expect ErrEmptyTasks",
				tasks:     []Task{},
				workerNum: 1,
				maxErrNum: 1,
				expErr:    ErrEmptyTasks,
			},
			{
				name:      "invalid nil tasks; expect ErrEmptyTasks",
				tasks:     nil,
				workerNum: 1,
				maxErrNum: 1,
				expErr:    ErrEmptyTasks,
			},
		} {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				err := Run(tc.tasks, tc.workerNum, tc.maxErrNum)
				assert.Truef(t, errors.Is(err, tc.expErr), "actual err - %v", err)
			})
		}
	})

	t.Run("case with maxErrorsCount==len(tasksCount); expect error", func(t *testing.T) {
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
		maxErrorsCount := 50
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.Equal(t, int32(tasksCount), runTasksCount, "not all tasks were completed")
	})

	t.Run("maxErrorsCount==len(tasksCount)==maxErrorsCount==1, task error occurred: expect error", func(t *testing.T) {
		tasksCount := 1
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 1
		maxErrorsCount := 1

		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.Equal(t, int32(tasksCount), runTasksCount, "not all tasks were completed")
	})

	t.Run("maxErrorsCount==len(tasksCount)==maxErrorsCount==1, no task error occurred: expect no error", func(t *testing.T) {
		tasksCount := 1
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 1
		maxErrorsCount := 1

		err := Run(tasks, workersCount, maxErrorsCount)

		require.NoError(t, err)
		require.Equal(t, int32(tasksCount), runTasksCount, "not all tasks were completed")
	})
}
