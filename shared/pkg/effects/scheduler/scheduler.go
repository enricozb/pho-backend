package scheduler

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/workers/pkg/lib/worker"
)

type Scheduler struct {
	db      *gorm.DB
	workers map[jobs.JobKind]worker.Worker

	SchedulerOptions
}

type SchedulerOptions struct {
	Concurrency     int
	PollingInterval time.Duration
}

func NewScheduler(db *gorm.DB, workers map[jobs.JobKind]worker.Worker, opts SchedulerOptions) *Scheduler {
	return &Scheduler{
		db:               db,
		workers:          workers,
		SchedulerOptions: opts,
	}
}

func (s *Scheduler) Run() error {
	g, ctx := errgroup.WithContext(context.Background())

	for proc := 0; proc < s.Concurrency; proc++ {
		g.Go(func() error {
			for {
				// error occured elsewhere...
				if err := ctx.Err(); err != nil {
					return nil
				}

				if err := s.ProcessNext(); err != nil {
					return err
				}

				time.Sleep(s.PollingInterval)
			}
		})
	}

	return g.Wait()
}

func (s *Scheduler) ProcessNext() error {
	job, jobExists, err := jobs.PopJob(s.db)
	if err != nil {
		return fmt.Errorf("pop job: %v", err)
	}

	if jobExists {
		if worker, workerExists := s.workers[job.Kind]; !workerExists {
			return fmt.Errorf("no worker for job kind: %s", job.Kind)
		} else if err := worker.Work(job); err != nil {
			return jobs.RecordJobFailure(s.db, job, err)
		}
	}

	return nil
}
