package scheduler_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/effects/scheduler"
	"github.com/enricozb/pho/shared/pkg/lib/testutil"
	"github.com/enricozb/pho/workers/pkg/lib/worker"
)

func setup(t *testing.T) (*require.Assertions, *gorm.DB, func()) {
	assert := require.New(t)
	db, cleanup := testutil.MockDB(t)

	return assert, db, cleanup
}

func TestScheduler_MissingWorker(t *testing.T) {
	assert, db, cleanup := setup(t)
	defer cleanup()

	importEntry := testutil.MockImport(t, db)

	_, err := jobs.PushJob(db, importEntry.ID, jobs.JobScan)
	assert.NoError(err, "push job")

	opts := scheduler.SchedulerOptions{
		Concurrency:     5,
		PollingInterval: 1 * time.Second,
	}

	s := scheduler.NewScheduler(db, map[jobs.JobKind]worker.Worker{}, opts)

	assert.EqualError(s.Run(), fmt.Sprintf("no worker for job kind: %s", jobs.JobScan))
}

func TestScheduler_WorkerErrors(t *testing.T) {
	assert, db, cleanup := setup(t)
	defer cleanup()

	expectedMessages := []string{}

	importEntry := testutil.MockImport(t, db)
	for i := 0; i < 5; i++ {
		_, err := jobs.PushJob(db, importEntry.ID, jobs.JobScan)
		assert.NoError(err, "push job")

		expectedMessages = append(expectedMessages, "mock error")
	}

	opts := scheduler.SchedulerOptions{
		Concurrency:     5,
		PollingInterval: 1 * time.Second,
	}

	workers := map[jobs.JobKind]worker.Worker{
		jobs.JobScan: worker.NewMockWorker(
			func(job jobs.Job) error {
				return errors.New("mock error")
			},
		),
	}

	s := scheduler.NewScheduler(db, workers, opts)

	// TODO(enricozb): ensure that s.Run() hasn't returned with an error
	go s.Run()

	// wait for errors to be recorded
	time.Sleep(1 * time.Second)

	assert.NoError(db.First(&importEntry).Error)
	assert.Equal(jobs.ImportStatusFailed, importEntry.Status)

	failures := []jobs.ImportFailure{}
	assert.NoError(db.Where("import_id = ?", importEntry.ID).Find(&failures).Error)

	actualMessages := []string{}
	for _, failure := range failures {
		actualMessages = append(actualMessages, failure.Message)
	}

	assert.ElementsMatch(expectedMessages, actualMessages)
}
