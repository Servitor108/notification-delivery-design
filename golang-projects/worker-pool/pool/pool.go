// Package pool provides a small bounded-concurrency worker pool.
package pool

import (
	"context"
	"errors"
	"sync"
)

var ErrClosed = errors.New("worker pool is closed")

// Job is a unit of work executed by a worker.
type Job func(context.Context) error

// Pool runs submitted jobs with a fixed number of workers.
type Pool struct {
	ctx    context.Context
	cancel context.CancelFunc
	jobs   chan Job
	wg     sync.WaitGroup

	mu     sync.RWMutex
	closed bool
}

// New creates a pool with workers concurrent workers and a queue of queueSize.
// A non-positive worker count is normalized to one; a negative queue size is zero.
func New(parent context.Context, workers, queueSize int) *Pool {
	if workers < 1 {
		workers = 1
	}
	if queueSize < 0 {
		queueSize = 0
	}
	ctx, cancel := context.WithCancel(parent)
	p := &Pool{ctx: ctx, cancel: cancel, jobs: make(chan Job, queueSize)}
	p.wg.Add(workers)
	for i := 0; i < workers; i++ {
		go p.worker()
	}
	return p
}

func (p *Pool) worker() {
	defer p.wg.Done()
	for {
		select {
		case <-p.ctx.Done():
			return
		case job, ok := <-p.jobs:
			if !ok {
				return
			}
			if job != nil {
				_ = job(p.ctx)
			}
		}
	}
}

// Submit queues a job, blocking while the queue is full. It returns when the
// job is accepted, the caller is cancelled, or the pool is closed.
func (p *Pool) Submit(ctx context.Context, job Job) error {
	if job == nil {
		return errors.New("job cannot be nil")
	}
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.closed {
		return ErrClosed
	}
	select {
	case p.jobs <- job:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-p.ctx.Done():
		return ErrClosed
	}
}

// Shutdown stops accepting new work, drains queued jobs, and waits for workers.
// If ctx expires, it cancels in-flight jobs and returns the context error.
func (p *Pool) Shutdown(ctx context.Context) error {
	p.mu.Lock()
	if !p.closed {
		p.closed = true
		close(p.jobs)
	}
	p.mu.Unlock()

	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		p.cancel()
		return nil
	case <-ctx.Done():
		p.cancel()
		return ctx.Err()
	}
}
