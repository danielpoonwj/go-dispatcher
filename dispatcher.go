package dispatcher

import (
	"sync"
)

// Dispatcher : Manage dispatching Jobs to Workers
type Dispatcher struct {
	jobQueue     chan Job
	workerPool   chan chan Job
	jobWaitGroup *sync.WaitGroup
}

// Start : Register and start workers
func (d *Dispatcher) Start() {
	for i := 0; i < cap(d.workerPool); i++ {
		w := newWorker(d.workerPool, d.jobWaitGroup)
		w.start()
	}

	go func() {
		for {
			select {
			case job := <-d.jobQueue:
				// try to obtain workerQueue from pool, will block until one is free
				workerQueue := <-d.workerPool

				// dispatch job to workerQueue
				workerQueue <- job
			}
		}
	}()
}

// AddJob : Add Job to job queue
func (d *Dispatcher) AddJob(job Job) {
	// Increment first - jobs may be processed faster than incr, avoid panic on negative count
	d.jobWaitGroup.Add(1)

	d.jobQueue <- job
}

// QueuedJobCount : Get number of jobs in queue. Note: some may already dispatched to the worker pool
func (d *Dispatcher) QueuedJobCount() int {
	return len(d.jobQueue)
}

// Stop : Ensure all jobs registered are processed
func (d *Dispatcher) Stop() {
	d.jobWaitGroup.Wait()
}

// NewDispatcher : Initialize new Dispatcher
func NewDispatcher(maxWorkers int, queueSize int) *Dispatcher {
	return &Dispatcher{
		jobQueue:     make(chan Job, queueSize),
		workerPool:   make(chan chan Job, maxWorkers),
		jobWaitGroup: &sync.WaitGroup{},
	}
}
