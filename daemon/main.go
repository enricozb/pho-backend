package main

import (
	"context"
	"flag"
	"time"

	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/effects/db"
	"github.com/enricozb/pho/shared/pkg/effects/scheduler"
	"github.com/enricozb/pho/workers/pkg/effects/workers"
	"github.com/enricozb/pho/workers/pkg/lib/worker"
)

type daemonOpts struct {
	concurrency     int
	pollingInterval time.Duration
}

func main() {
	d := daemonOpts{}

	flag.IntVar(&d.concurrency, "concurrency", 5, "The number of workers used for the job scheduler. Must be at least 2.")
	flag.DurationVar(&d.pollingInterval, "interval", 1*time.Second, "The amount of time to wait between polling for new jobs.")
	flag.Parse()

	db := db.MustDB()
	s := scheduler.NewScheduler(
		db,
		buildWorkers(db),
		scheduler.SchedulerOptions{
			Concurrency:     d.concurrency,
			PollingInterval: d.pollingInterval,
		},
	)

	g, _ := errgroup.WithContext(context.Background())
	g.Go(s.Run)
	// g.Go(api.run)

	if err := g.Wait(); err != nil {
		panic(err)
	}
}

func buildWorkers(db *gorm.DB) map[jobs.JobKind]worker.Worker {
	return map[jobs.JobKind]worker.Worker{
		jobs.JobScan:            workers.NewScanWorker(db),
		jobs.JobMetadata:        workers.NewMetadataWorker(db),
		jobs.JobMetadataHash:    workers.NewHashWorker(db),
		jobs.JobMetadataEXIF:    workers.NewEXIFWorker(db),
		jobs.JobMetadataMonitor: workers.NewMonitorWorker(db, jobs.JobDedupe),
		jobs.JobDedupe:          workers.NewDedupeWorker(db),
		jobs.JobConvert:         workers.NewConvertWorker(db),
		jobs.JobCleanup:         workers.NewCleanupWorker(db),
	}
}
