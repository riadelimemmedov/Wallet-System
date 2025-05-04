package pkg_utils

import (
	"fmt"

	db "github.com/riad/banksystemendtoend/db/sqlc"
	"github.com/riad/banksystemendtoend/pkg/model"
)

func ConvertDBStatusToModelStatus(status db.UploadStatus) model.UploadStatus {
	switch status {
	case db.UploadStatusPENDING:
		return model.UploadStatusPending
	case db.UploadStatusPROCESSING:
		return model.UploadStatusProcessing
	case db.UploadStatusCOMPLETED:
		return model.UploadStatusCompleted
	case db.UploadStatusFAILED:
		return model.UploadStatusFailed
	case db.UploadStatusCANCELLED:
		return model.UploadStatusCancelled
	default:
		return model.UploadStatusFailed
	}
}

// convertModelStatusToDBStatus converts from model.UploadStatus to db.UploadStatus
func ConvertModelStatusToDBStatus(status model.UploadStatus) (db.UploadStatus, error) {
	switch status {
	case model.UploadStatusPending:
		return db.UploadStatusPENDING, nil
	case model.UploadStatusProcessing:
		return db.UploadStatusPROCESSING, nil
	case model.UploadStatusCompleted:
		return db.UploadStatusCOMPLETED, nil
	case model.UploadStatusFailed:
		return db.UploadStatusFAILED, nil
	case model.UploadStatusCancelled:
		return db.UploadStatusCANCELLED, nil
	default:
		return db.UploadStatusFAILED, fmt.Errorf("invalid upload status: %s", status)
	}
}
