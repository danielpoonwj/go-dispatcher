package dispatcher

import (
	"sync"
)

// Dispatcher : Manage dispatching Jobs to Workers
type Dispatcher struct {
	JobQueue     chan Job
	WorkerPool   chan chan Job
	jobWaitGroup *sync.WaitGroup
}

// Start : Register and start workers
func (d *Dispatcher) Start() {
	for i := 0; i < cap(d.WorkerPool); i++ {
		worker := NewWorker(d.WorkerPool, d.jobWaitGroup)
		worker.start()
	}

	go func() {
		for {
			select {
			case job := <-d.JobQueue:
				// try to obtain workerQueue from pool, will block until one is free
				workerQueue := <-d.WorkerPool

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

	d.JobQueue <- job
}

// Stop : Ensure all jobs registered are processed
func (d *Dispatcher) Stop() {
	d.jobWaitGroup.Wait()
}

// NewDispatcher : Initialize new Dispatcher
func NewDispatcher(maxWorkers int, queueSize int) *Dispatcher {
	return &Dispatcher{
		JobQueue:     make(chan Job, queueSize),
		WorkerPool:   make(chan chan Job, maxWorkers),
		jobWaitGroup: &sync.WaitGroup{},
	}
}
