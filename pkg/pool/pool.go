package pool

import "sync"

// WorkerFn is func for concurrently execution.
type WorkerFn func()

// Pool represents a control component for WorkerFn.
type Pool struct {
	size    int
	workers chan WorkerFn
	wg      *sync.WaitGroup
}

// start launches required number of goroutines.
func (wp *Pool) start() {
	wp.wg.Add(wp.size)

	for i := 0; i < wp.size; i++ {
		go wp.worker()
	}
}

// worker get function from workers queue and executes it.
func (wp *Pool) worker() {
	defer wp.wg.Done()

	for fn := range wp.workers {
		fn()
	}
}

// Add appends worker func to execution queue.
func (wp *Pool) Add(fn WorkerFn) {
	wp.workers <- fn
}

// Wait locks execution until all workers complete.
func (wp *Pool) Wait() {
	close(wp.workers)
	wp.wg.Wait()
}

// NewWorkersPool creates *WorkersPool and starts workers pool.
func New(size int) *Pool {
	wp := &Pool{
		size:    size,
		workers: make(chan WorkerFn),
		wg:      &sync.WaitGroup{},
	}

	go wp.start()

	return wp
}
