package jobs

import "gorm.io/gorm"

type ImportStatus string
type JobStatus string

const (
	ImportStatusNotStarted ImportStatus = "NOT_STARTED"
	ImportStatusScan       ImportStatus = "SCAN"
	ImportStatusMetadata   ImportStatus = "METADATA"
	ImportStatusDedupe     ImportStatus = "DEDUPE"
	ImportStatusConvert    ImportStatus = "CONVERT"
	ImportStatusCleanup    ImportStatus = "CLEANUP"
	ImportStatusDone       ImportStatus = "DONE"
	ImportStatusFailed     ImportStatus = "FAILED"
)

const (
	JobStatusNotStarted JobStatus = "NOT_STARTED"
	JobStatusStarted    JobStatus = "STARTED"
	JobStatusDone       JobStatus = "DONE"
	JobStatusFailed     JobStatus = "FAILED"
)

func GetJobStatus(db *gorm.DB, jobID JobID) (status JobStatus, err error) {
	var job Job
	return job.Status, db.Model(&job).Where("id = ?", jobID).Error
}
