package scheduler_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/enricozb/pho/shared/pkg/effects/daos/jobs"
	"github.com/enricozb/pho/shared/pkg/effects/scheduler"
	"github.com/enricozb/pho/shared/pkg/effects/worker"
	"github.com/enricozb/pho/shared/pkg/lib/testutil"
)

func setup(t *testing.T) (*require.Assertions, *sqlx.DB, jobs.Dao, func()) {
	assert := require.New(t)
	db, cleanup := testutil.MockDB(t)
	dao := jobs.NewDao(db)

	return assert, db, dao, cleanup
}

func TestScheduler_MissingWorker(t *testing.T) {
	assert, db, dao, cleanup := setup(t)
	defer cleanup()

	importID := testutil.MockImport(t, db)
	_, err := dao.PushJob(importID, jobs.JobScan)
	assert.NoError(err, "push job")

	opts := scheduler.SchedulerOptions{
		Concurrency:     5,
		PollingInterval: 1 * time.Second,
	}

	s := scheduler.NewScheduler(dao, map[jobs.JobKind]worker.Worker{}, opts)

	assert.EqualError(s.Run(), fmt.Sprintf("no worker for job kind: %s", jobs.JobScan))
}

func TestScheduler_WorkerErrors(t *testing.T) {
	assert, db, dao, cleanup := setup(t)
	defer cleanup()

	expectedMessages := []string{}

	importID := testutil.MockImport(t, db)
	for i := 0; i < 5; i++ {
		_, err := dao.PushJob(importID, jobs.JobScan)
		assert.NoError(err, "push job")

		expectedMessages = append(expectedMessages, "mock error")
	}

	opts := scheduler.SchedulerOptions{
		Concurrency:     5,
		PollingInterval: 1 * time.Second,
	}

	workers := map[jobs.JobKind]worker.Worker{
		jobs.JobScan: worker.NewMockWorker(
			func(importID jobs.ImportID) error {
				return errors.New("mock error")
			},
		),
	}

	s := scheduler.NewScheduler(dao, workers, opts)

	// TODO(enricozb): ensure that s.Run() hasn't returned with an error
	go s.Run()

	// wait for errors to be recorded
	time.Sleep(1 * time.Second)

	actualMessages, err := dao.GetImportFailureMessages(importID)
	assert.NoError(err)

	assert.ElementsMatch(expectedMessages, actualMessages)
}
