package jobs

import (
	"fmt"
	"sync"

	"gorm.io/gorm"
)

var queueMutex sync.Mutex

func NumJobs(db *gorm.DB) (count int64, err error) {
	return count, db.Model(&Job{}).Where("status = ?", JobStatusNotStarted).Count(&count).Error
}

func PushJob(db *gorm.DB, importID ImportID, kind JobKind) (Job, error) {
	job := Job{ImportID: importID, Kind: kind, Status: JobStatusNotStarted}
	return job, db.Create(&job).Error
}

func PopJob(db *gorm.DB) (job Job, jobExists bool, err error) {
	queueMutex.Lock()
	defer queueMutex.Unlock()

	numJobs, err := NumJobs(db)
	if err != nil {
		return Job{}, false, fmt.Errorf("num jobs: %v", err)
	} else if numJobs == 0 {
		return Job{}, false, fmt.Errorf("peek job: %v", err)
	}

	if err := db.Where("status = ?", JobStatusNotStarted).First(&job).Error; err != nil {
		return Job{}, false, fmt.Errorf("first: %v", err)
	}

	job.Status = JobStatusStarted
	if err := db.Save(&job).Error; err != nil {
		return Job{}, false, fmt.Errorf("save: %v", err)
	}

	return job, true, nil
}
