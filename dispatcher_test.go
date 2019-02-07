package dispatcher

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockJob : Mock Job
type MockJob struct {
	mock.Mock
}

// Process : Job Process method
func (m *MockJob) Process() {
	m.Called()
}

// NewMockJob : Create new MockJob
func NewMockJob(delay time.Duration) *MockJob {
	job := new(MockJob)
	job.On("Process").WaitUntil(time.After(delay)).Return()

	return job
}

func TestDispatcher(t *testing.T) {
	dispatcher := NewDispatcher(1, 3)
	dispatcher.Start()

	// simulate time taken for job processing
	delayTime := 500 * time.Millisecond
	// short wait time for job to enter pipeline
	initWaitTime := 10 * time.Millisecond

	j1 := NewMockJob(delayTime)
	dispatcher.AddJob(j1)
	time.Sleep(initWaitTime)

	assert.Equal(t, 0, dispatcher.QueuedJobCount(), "No Jobs should be in queue")
	j1.AssertCalled(t, "Process")

	j2 := NewMockJob(delayTime)
	dispatcher.AddJob(j2)
	time.Sleep(initWaitTime)

	// quirk: number of jobs processing can be workers + 1
	// happens when job is received from JobQueue but no workers free yet
	// technically in limbo but will be considered processing Job
	assert.Equal(t, 0, dispatcher.QueuedJobCount(), "No Jobs should be in queue")
	j2.AssertNotCalled(t, "Process")

	j3 := NewMockJob(delayTime)
	dispatcher.AddJob(j3)
	time.Sleep(initWaitTime)

	assert.Equal(t, 1, dispatcher.QueuedJobCount(), "j3 should be in queue")
	j3.AssertNotCalled(t, "Process")

	dispatcher.Stop()

	// Ensure all jobs have been processed
	j1.AssertCalled(t, "Process")
	j2.AssertCalled(t, "Process")
	j3.AssertCalled(t, "Process")

	assert.Equal(t, 0, dispatcher.QueuedJobCount(), "No jobs should be in queue")
}
