package jobs

import (
	"fmt"

	"gorm.io/gorm"
)

func (job *Job) RecordFailure(db *gorm.DB, err error) error {
	if err := db.Model(&Job{}).Where("id = ?", job.ID).Update("status", JobStatusFailed).Error; err != nil {
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
