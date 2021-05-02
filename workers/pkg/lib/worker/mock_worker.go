package worker

import "github.com/enricozb/pho/shared/pkg/effects/daos/jobs"

// TODO(enricozb): maybe switch to https://github.com/brianvoe/gofakeit
type MockWorker struct {
	callback func(job jobs.Job) error
}

var _ Worker = &MockWorker{}

func NewMockWorker(callback func(job jobs.Job) error) *MockWorker {
	return &MockWorker{callback: callback}
}

func (w *MockWorker) Work(job jobs.Job) error {
	return w.callback(job)
}
