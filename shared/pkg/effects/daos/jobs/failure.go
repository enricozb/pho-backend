package jobs

import (
	"fmt"

	"gorm.io/gorm"
)

func RecordJobFailure(db *gorm.DB, job Job, err error) error {
	job.Status = JobStatusFailed
	if err := db.Save(&job).Error; err != nil {
		return fmt.Errorf("update job status: %v", err)
	}

	if err := db.Model(&Import{}).Where("id", job.ImportID).Update("status", ImportStatusFailed).Error; err != nil {
		return fmt.Errorf("update import status: %v", err)
	}

	if err := db.Create(&ImportFailure{ImportID: job.ImportID, Message: err.Error()}).Error; err != nil {
		return fmt.Errorf("create import failure: %v", err)
	}

	return nil
}
