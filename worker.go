package dispatcher

import (
	"sync"
)

// Worker : Generic worker to manage processing Jobs
type Worker struct {
	workerQueue  chan Job
	workerPool   chan chan Job
	jobWaitGroup *sync.WaitGroup
}

// start : Register WorkQueue to WorkerPool and start processing jobs
func (w *Worker) start() {
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

func (w *Worker) handleJob(job Job) {
	// Decrement after job has been handled
	defer w.jobWaitGroup.Done()

	job.Process()
}

// NewWorker : Initialize worker
func NewWorker(workerPool chan chan Job, wg *sync.WaitGroup) *Worker {
	return &Worker{
		workerQueue:  make(chan Job),
		workerPool:   workerPool,
		jobWaitGroup: wg,
	}
}
