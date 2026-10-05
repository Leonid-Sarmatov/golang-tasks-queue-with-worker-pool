package workerpool

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"worker_pool/internal/app/conf"
	"worker_pool/internal/app/execute"
	"worker_pool/internal/core/domain"
)

var (
	ErrQueueIsLocked = errors.New("worker pool has stopped, queue is locked")
)

type WorkerPool struct {
	ctxShutdown context.Context
	workersNum  domain.WorkersNumber
	queue       chan domain.Task
	semaphore   chan struct{}
	wg          sync.WaitGroup
	cfg         conf.IConfig
}

func NewWorkerPool(
	ctxShutdown context.Context,
	cfg conf.IConfig,
) *WorkerPool {
	var wp WorkerPool

	wp.ctxShutdown = ctxShutdown

	wp.workersNum = cfg.GetWorkersNumber()
	wp.queue = make(chan domain.Task, cfg.GetTaskQueueSize())

	wp.cfg = cfg

	wp.wg = sync.WaitGroup{}

	return &wp
}

func (wp *WorkerPool) AddTaskToQueue(taskIn domain.Task) error {
	if wp.ctxShutdown.Err() != nil {
		return ErrQueueIsLocked
	}

	// Transition task to queued state
	taskOut, err := domain.TaskToQueue(taskIn)
	if err != nil {
		return err
	}

	// Put task into queue, unless shutdown signaled
	select {
	case <-wp.ctxShutdown.Done():
		return ErrQueueIsLocked
	case wp.queue <- taskOut:
		return nil
	}
}

func (wp *WorkerPool) Run() {
	// Start worker goroutines
	for i := 0; i < int(wp.workersNum); i += 1 {
		wp.wg.Add(1)
		go func() {
			defer wp.wg.Done()
			for {
				// Pull task from queue, unless shutdown signaled
				select {
				case <-wp.ctxShutdown.Done():
					return
				case t := <-wp.queue:
					wp.processor(t)
				}
			}
		}()
	}
}

func (wp *WorkerPool) Stop() {
	wp.wg.Wait()
}

func (wp *WorkerPool) processor(t domain.Task) {
	select {
	case <-wp.ctxShutdown.Done():
		return

	default:
		// Transition task to executing state
		t, err := domain.ExecuteTask(t)
		if err != nil {
			fmt.Printf("cannot execute task, task=%v, err=%v\n", t, err)
			return
		}

		// Execute task
		t, errExec := execute.Execute(t, wp.cfg)
		if errExec != nil {
			// Transition task to fail state if execution failed
			t, err = domain.FailTask(t)
			if err != nil {
				fmt.Printf("failed to mark task as failed, task=%v, err=%v\n", t, err)
				return
			}

			fmt.Printf("failed to execute task, task=%v, err=%v\n", t, errExec)

			// Prepare task for next attempt
			t, err = domain.TaskToQueue(t)
			if err != nil {
				fmt.Printf("failed to put task in queue, task=%v, err=%v\n", t, err)
				return
			}

			wp.retry(t)

			return
		}

		// Mark task as completed
		t, err = domain.DoneTask(t)
		if err != nil {
			fmt.Printf("failed to mark task as completed, task=%v, err=%v\n", t, err)
		}
	}

}

func (wp *WorkerPool) retry(t domain.Task) {
	// Calculate exponential backoff delay
	x := time.Duration(1<<uint(t.CurrentRetries)) * time.Second

	// Add random jitter to delay to spread retries
	x += time.Duration(rand.Intn(1000)-500) * time.Millisecond

	// Include pending retries in the shutdown wait
	wp.wg.Add(1)
	go func() {
		defer wp.wg.Done()

		timer := time.NewTimer(x)
		defer timer.Stop()

		// Wait for the delay, unless shutdown signaled
		select {
		case <-wp.ctxShutdown.Done():
			return
		case <-timer.C:
		}

		// Put task into queue, unless shutdown signaled
		select {
		case <-wp.ctxShutdown.Done():
			return
		case wp.queue <- t:
		}
	}()
}
