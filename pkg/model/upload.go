package model

import "time"

// UploadStatus represents the status of a file upload
type UploadStatus string

const (
	UploadStatusPending    UploadStatus = "PENDING"
	UploadStatusProcessing UploadStatus = "PROCESSING"
	UploadStatusCompleted  UploadStatus = "COMPLETED"
	UploadStatusFailed     UploadStatus = "FAILED"
	UploadStatusCancelled  UploadStatus = "CANCELLED"
)

// UploadJob represents a file upload job
type UploadJob struct {
	ID          string       `json:"id"`
	FileName    string       `json:"file_name"`
	FileSize    int64        `json:"file_size"`
	ContentType string       `json:"content_type"`
	TempPath    string       `json:"temp_path"`
	TargetPath  string       `json:"target_path"`
	UserID      int64        `json:"user_id"`
	Status      UploadStatus `json:"status"`
	Error       string       `json:"error,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	CompletedAt *time.Time   `json:"completed_at,omitempty"`
}

// FileMetadata represents additional information about an uploaded file
type FileMetadata struct {
	ID             string    `json:"id"`
	S3URL          string    `json:"s3_url"`
	Checksum       string    `json:"checksum,omitempty"`
	MimeType       string    `json:"mime_type"`
	Width          int       `json:"width,omitempty"`
	Height         int       `json:"height,omitempty"`
	AdditionalData []byte    `json:"additional_data,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
