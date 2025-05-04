package upload_service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	pkg_interface "github.com/riad/banksystemendtoend/pkg/interface"
	logger "github.com/riad/banksystemendtoend/pkg/log"
	"github.com/riad/banksystemendtoend/pkg/model"
	"github.com/riad/banksystemendtoend/pkg/rabbitmq"
	"go.uber.org/zap"
)

/* The File Upload Flow

User uploads a file to the API
HandleFileUpload validates the file and saves it to temporary storage
An UploadJob is created with PENDING status and stored in the repository
The job is published to RabbitMQ exchange, which routes it to the upload queue
Worker(s) consume jobs from the queue and process them
When processing completes, the job status is updated to COMPLETED or FAILED
Users can query job status via GetUploadStatus

*/

// UploadService handles file upload operations
type UploadService struct {
	tempDir        string
	maxUploadSize  int64
	rmqClient      *rabbitmq.Client
	uploadRepo     pkg_interface.UploadRepository
	uploadExchange string
	uploadQueue    string
	mu             sync.Mutex
}

func NewUploadService(
	tempDir string,
	maxUploadSize int64,
	rmqClient *rabbitmq.Client,
	uploadRepo pkg_interface.UploadRepository,
	uploadExchange string,
	uploadQueue string,
) (*UploadService, error) {
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		logger.GetLogger().Error("failed to create temp directory", zap.Error(err))
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	if err := rmqClient.DeclareExchange(uploadExchange, "direct", true, false, false); err != nil {
		logger.GetLogger().Error("failed to declare exchange", zap.Error(err))
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	if _, err := rmqClient.DeclareQueue(uploadQueue, true, false, false); err != nil {
		logger.GetLogger().Error("failed to declare queue", zap.Error(err))
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	if err := rmqClient.BindQueue(uploadQueue, uploadQueue, uploadExchange); err != nil {
		logger.GetLogger().Error("failed to bind queue", zap.Error(err))
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	return &UploadService{
		tempDir:        tempDir,
		maxUploadSize:  maxUploadSize,
		rmqClient:      rmqClient,
		uploadRepo:     uploadRepo,
		uploadExchange: uploadExchange,
		uploadQueue:    uploadQueue,
	}, nil
}

/*
HandleFileUpload handles the file upload process:
1. Saves file to temporary storage
2. Creates an upload job record
3. Publishes the job to RabbitMQ queue
4. Returns the job ID for status tracking
*/
func (s *UploadService) HandleFileUpload(ctx context.Context, file *multipart.FileHeader, userID int32, targetPath string) (string, error) {
	if file.Size > s.maxUploadSize {
		logger.GetLogger().Error("file size exceeds limit", zap.Int64("size", file.Size), zap.Int64("limit", s.maxUploadSize))
		return "", fmt.Errorf("file size exceeds limit: %d > %d", file.Size, s.maxUploadSize)
	}

	uploadID := uuid.New().String()

	job := &model.UploadJob{
		ID:          uploadID,
		FileName:    file.Filename,
		FileSize:    file.Size,
		ContentType: file.Header.Get("Content-Type"),
		UserID:      userID,
		TempPath:    filepath.Join(s.tempDir, uploadID),
		TargetPath:  targetPath,
		Status:      model.UploadStatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.saveToTempStorage(file, job.TempPath); err != nil {
		logger.GetLogger().Error("failed to save file to temp storage", zap.Error(err))
		return "", fmt.Errorf("failed to save file to temp storage: %w", err)
	}

	if err := s.uploadRepo.CreateUploadJob(ctx, job); err != nil {
		os.Remove(job.TempPath)
		logger.GetLogger().Error("failed to create upload job", zap.Error(err))
		return "", fmt.Errorf("failed to create upload job: %w", err)
	}

	if err := s.rmqClient.PublishJSON(ctx, s.uploadExchange, s.uploadQueue, job); err != nil {
		job.Status = model.UploadStatusFailed
		job.Error = fmt.Sprintf("failed to publish job to RabbitMQ: %v", err)
		s.uploadRepo.UpdateUploadJob(ctx, job)
		logger.GetLogger().Error("failed to publish job to RabbitMQ", zap.Error(err))
		return "", fmt.Errorf("failed to publish job to RabbitMQ: %w", err)
	}

	logger.GetLogger().Info("file upload job created", zap.String("job_id", uploadID), zap.String("file_name", file.Filename), zap.Int64("file_size", file.Size))
	return uploadID, nil
}

// saveToTempStorage saves the uploaded file to temporary storage
func (s *UploadService) saveToTempStorage(file *multipart.FileHeader, tempPath string) error {
	src, err := file.Open()
	if err != nil {
		logger.GetLogger().Error("failed to open file", zap.Error(err))
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(tempPath)
	if err != nil {
		logger.GetLogger().Error("failed to create temp file", zap.Error(err))
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		logger.GetLogger().Error("failed to copy file", zap.Error(err))
		return fmt.Errorf("failed to copy file: %w", err)
	}
	return nil
}

// GetUploadStatus retrieves the status of an upload job
func (s *UploadService) GetUploadStatus(ctx context.Context, jobID int64) (*model.UploadJob, error) {
	job, err := s.uploadRepo.GetUploadJob(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("failed to get upload job: %w", err)
	}
	return job, nil
}

// CalculateFileChecksum calculates the MD5 checksum of a file
func (s *UploadService) CalculateFileChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to calculate file hash: %w", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// CleanupExpiredFiles removes temporary files older than specified duration
func (s *UploadService) CleanupExpiredFiles(ctx context.Context, maxAge time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entries, err := os.ReadDir(s.tempDir)
	if err != nil {
		return fmt.Errorf("failed to read temp directory: %w", err)
	}

	cutoff := time.Now().Add(-maxAge)
	var removed int

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(s.tempDir, entry.Name())
		info, err := entry.Info()
		if err != nil {
			logger.GetLogger().Error("Failed to get file info",
				zap.String("path", filePath),
				zap.Error(err))
			continue
		}

		if info.ModTime().Before(cutoff) {
			if err := os.Remove(filePath); err != nil {
				logger.GetLogger().Error("Failed to remove expired file",
					zap.String("path", filePath),
					zap.Error(err))
				continue
			}
			removed++
		}
	}

	logger.GetLogger().Info("Cleaned up expired temporary files",
		zap.Int("removed", removed),
		zap.Duration("maxAge", maxAge))

	return nil
}
