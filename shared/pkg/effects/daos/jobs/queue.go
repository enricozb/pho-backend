package jobs

import (
	"encoding/json"
	"fmt"
	"sync"

	"gorm.io/gorm"
)

var queueMutex sync.Mutex

func NumJobs(db *gorm.DB) (count int64, err error) {
	return count, db.Model(&Job{}).Where("status = ?", JobStatusNotStarted).Count(&count).Error
}

func PushJob(db *gorm.DB, importID ImportID, kind JobKind) (Job, error) {
	job := Job{ImportID: importID, Kind: kind}
	return job, db.Create(&job).Error
}

func PushJobWithArgs(db *gorm.DB, importID ImportID, kind JobKind, args interface{}) (Job, error) {
	data, err := json.Marshal(args)
	if err != nil {
		return Job{}, fmt.Errorf("marshal: %v", err)
	}

	job := Job{ImportID: importID, Kind: kind, Args: data}
	return job, db.Create(&job).Error
}

func PopJob(db *gorm.DB) (job Job, jobExists bool, err error) {
	queueMutex.Lock()
	defer queueMutex.Unlock()

	numJobs, err := NumJobs(db)
	if err != nil {
		return Job{}, false, fmt.Errorf("num jobs: %v", err)
	} else if numJobs == 0 {
		return Job{}, false, nil
	}

	if err := db.Where("status = ?", JobStatusNotStarted).First(&job).Error; err != nil {
		return Job{}, false, fmt.Errorf("first: %v", err)
	}

	if err := job.SetStatus(db, JobStatusStarted); err != nil {
		return Job{}, false, fmt.Errorf("save: %v", err)
	}

	return job, true, nil
}
