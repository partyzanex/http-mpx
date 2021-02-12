package pool

import (
	"context"
	"sync/atomic"
)

// WorkerFn is func for concurrently execution.
type Worker func()

// Pool represents a control component for WorkerFn.
type Pool struct {
	// count of executed WorkerFn at one time
	size int32
	// atomic counter completed workers
	counter int32
	// workers queue
	workers chan Worker
	// wait channel
	wait chan struct{}
	// context for exit
	ctx context.Context
}

// done decrements cnt and check it
// if cnt less or equal than 0 send to wait for unlock.
func (p *Pool) done() {
	atomic.AddInt32(&p.counter, -1)

	if p.counter <= 0 {
		// unlock wait
		p.wait <- struct{}{}
	}
}

// worker get function from workers queue and executes it.
func (p *Pool) worker() {
	defer p.done()

	for {
		select {
		case <-p.ctx.Done():
			return
		case fn, ok := <-p.workers:
			if fn != nil {
				fn()
			}

			if ok {
				continue
			}

			return
		}
	}
}

// Add appends worker func to execution queue.
func (p *Pool) Add(fn Worker) {
	select {
	case <-p.ctx.Done():
		return
	case p.workers <- fn:
	}
}

// Wait locks execution until all workers complete.
func (p *Pool) Wait() {
	close(p.workers)
	<-p.wait
}

// NewWorkersPool creates *WorkersPool and starts workers pool.
func New(size int, options ...Option) *Pool {
	p := &Pool{
		size: int32(size),
		wait: make(chan struct{}, 1),
		ctx:  context.Background(),
	}

	// apply options
	for _, option := range options {
		option(p)
	}

	if p.workers == nil {
		p.workers = make(chan Worker)
	}

	go func() {
		atomic.AddInt32(&p.counter, p.size)

		for i := int32(0); i < p.size; i++ {
			go p.worker()
		}
	}()

	return p
}
