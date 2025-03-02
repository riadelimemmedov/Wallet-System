package service_interface

import (
	"context"
	"mime/multipart"
)

// S3Service defines the interface for S3 storage operations
type S3Service interface {
	// UploadFile uploads a file to S3 and returns the URL
	UploadFile(ctx context.Context, file *multipart.FileHeader) (string, error)

	// DeleteFile removes a file from S3
	DeleteFile(ctx context.Context, fileURL string) error
}
