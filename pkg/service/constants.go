package upload_service

// UploadStatus represents the status of a file upload
type UploadStatus string

const (
	UploadStatusPending    UploadStatus = "PENDING"
	UploadStatusProcessing UploadStatus = "PROCESSING"
	UploadStatusCompleted  UploadStatus = "COMPLETED"
	UploadStatusFailed     UploadStatus = "FAILED"
)
