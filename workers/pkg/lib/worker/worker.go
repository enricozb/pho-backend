package worker

import "github.com/enricozb/pho/shared/pkg/effects/daos/jobs"

type Worker interface {
	Work(jobs.ImportID) error
}
