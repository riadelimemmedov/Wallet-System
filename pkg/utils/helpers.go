package pkg_utils

import (
	db "github.com/riad/banksystemendtoend/db/sqlc"
	"github.com/riad/banksystemendtoend/pkg/model"
)

func convertDBStatusToModelStatus(status db.UploadStatus) model.UploadStatus {
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
