package worker

import "github.com/enricozb/pho/shared/pkg/effects/daos/jobs"

// TODO(enricozb): maybe switch to https://github.com/brianvoe/gofakeit
type MockWorker struct {
	callback func(importID jobs.ImportID) error
}

var _ Worker = &MockWorker{}

func NewMockWorker(callback func(importID jobs.ImportID) error) *MockWorker {
	return &MockWorker{callback: callback}
}

func (w *MockWorker) Work(importID jobs.ImportID) error {
	return w.callback(importID)
}
