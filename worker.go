package dispatcher

import (
	"sync"
)

// worker : Generic worker to manage processing Jobs
type worker struct {
	workerQueue  chan Job
	workerPool   chan chan Job
	jobWaitGroup *sync.WaitGroup
}

// start : Register WorkQueue to WorkerPool and start processing jobs
func (w *worker) start() {
	go func() {
		for {
			// Adds itself to the worker pool
			w.workerPool <- w.workerQueue

			select {
			case job := <-w.workerQueue:
				w.handleJob(job)
			}
		}
	}()
}

func (w *worker) handleJob(job Job) {
	// Decrement after job has been handled
	defer w.jobWaitGroup.Done()

	job.Process()
}

// newWorker : Initialize worker
func newWorker(workerPool chan chan Job, wg *sync.WaitGroup) *worker {
	return &worker{
		workerQueue:  make(chan Job),
		workerPool:   workerPool,
		jobWaitGroup: wg,
	}
}
