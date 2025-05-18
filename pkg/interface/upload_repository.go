package pkg_interface

import (
	"context"

	"github.com/riad/banksystemendtoend/pkg/model"
)

// UploadRepository provides operations for storing and retrieving upload jobs
type UploadRepository interface {
	// CreateUploadJob stores a new upload job
	CreateUploadJob(ctx context.Context, job *model.UploadJob) error

	// GetUploadJob retrieves an upload job by ID
	GetUploadJob(ctx context.Context, id int64) (*model.UploadJob, error)

	// UpdateUploadJob updates an existing upload job
	UpdateUploadJob(ctx context.Context, job *model.UploadJob) error

	// ListPendingJobs retrieves a list of pending upload jobs
	ListPendingJobs(ctx context.Context, limit int) (*[]model.UploadJob, error)

	// ListFailedJobs retrieves a list of failed upload jobs
	ListFailedJobs(ctx context.Context, limit int) (*[]model.UploadJob, error)

	// ListProcessingJobs retrieves a list of currently processing upload jobs
	ListProcessingJobs(ctx context.Context, limit int) (*[]model.UploadJob, error)

	// ListUserJobs retrieves a user's upload jobs with pagination
	ListUserJobs(ctx context.Context, userID int64, limit, offset int) (*[]model.UploadJob, error)

	// CreateFileMetadata stores metadata for a successfully uploaded file
	CreateFileMetadata(ctx context.Context, metadata *model.FileMetadata) error

	// GetFileMetadata retrieves file metadata by ID
	GetFileMetadata(ctx context.Context, id int64) (*model.FileMetadata, error)
}
