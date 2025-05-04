package pkg_repository

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/riad/banksystemendtoend/db/sqlc"
	"github.com/riad/banksystemendtoend/pkg/cache"
	pkg_interface "github.com/riad/banksystemendtoend/pkg/interface"
	logger "github.com/riad/banksystemendtoend/pkg/log"
	"github.com/riad/banksystemendtoend/pkg/model"
	pkg_utils "github.com/riad/banksystemendtoend/pkg/utils"
	"go.uber.org/zap"
)

const (
	uploadJobCacheKey = "upload_jobs:%s"
)

// SQLUploadRepository is a PostgreSQL implementation of UploadRepository
type SQLUploadRepository struct {
	store       db.Store
	cacheClient *cache.Service
}

// NewSQLUploadRepository creates a new SQLUploadRepository
func NewSQLUploadRepository(store db.Store, cacheClient *cache.Service) pkg_interface.UploadRepository {
	return &SQLUploadRepository{
		store:       store,
		cacheClient: cacheClient,
	}
}

// CreateUploadJob stores a new upload job in the database
func (r *SQLUploadRepository) CreateUploadJob(ctx context.Context, job *model.UploadJob) error {
	sqlStore, err := db.GetSQLStore(r.store)
	if err != nil {
		logger.GetLogger().Error("failed to get SQL store", zap.Error(err))
		return fmt.Errorf("failed to get SQL store: %w", err)
	}

	status, err := pkg_utils.ConvertModelStatusToDBStatus(job.Status)
	if err != nil {
		logger.GetLogger().Error("failed to convert model status to DB status", zap.Error(err))
		return fmt.Errorf("failed to convert model status to DB status: %w", err)
	}

	arg := db.CreateUploadJobParams{
		ID:           job.ID,
		FileName:     job.FileName,
		FileSize:     job.FileSize,
		ContentType:  job.ContentType,
		TempPath:     job.TempPath,
		TargetPath:   job.TargetPath,
		UserID:       job.UserID,
		Status:       status,
		ErrorMessage: sql.NullString{String: job.Error, Valid: job.Error != ""},
		CreatedAt:    job.CreatedAt,
		UpdatedAt:    job.UpdatedAt,
	}

	_, err = sqlStore.CreateUploadJob(ctx, arg)
	if err != nil {
		return fmt.Errorf("failed to create upload job: %w", err)
	}

	// Cache the job
	cacheKey := fmt.Sprintf(uploadJobCacheKey, job.ID)
	if err := r.cacheClient.Set(ctx, cacheKey, job); err != nil {
		logger.GetLogger().Error("failed to cache upload job", zap.String("cacheKey", cacheKey), zap.Error(err))
		return fmt.Errorf("failed to cache upload job: %w", err)
	}
	return nil
}

// GetUploadJob retrieves an upload job by its ID
func (r *SQLUploadRepository) GetUploadJob(ctx context.Context, id string) (*model.UploadJob, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf(uploadJobCacheKey, id)

	var job model.UploadJob

	err := r.cacheClient.Get(ctx, cacheKey, &job)
	if err == nil {
		return &job, nil
	}

	// If not in cache, get from database using SQLC-generated function
	sqlStore, err := db.GetSQLStore(r.store)
	if err != nil {
		return nil, fmt.Errorf("failed to get SQL store: %w", err)
	}

	dbJob, err := sqlStore.GetUploadJob(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("upload job not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get upload job: %w", err)
	}

	job = model.UploadJob{
		ID:          dbJob.ID,
		FileName:    dbJob.FileName,
		FileSize:    dbJob.FileSize,
		ContentType: dbJob.ContentType,
		TempPath:    dbJob.TempPath,
		TargetPath:  dbJob.TargetPath,
		UserID:      dbJob.UserID,
		Status:      pkg_utils.ConvertDBStatusToModelStatus(dbJob.Status),
		CreatedAt:   dbJob.CreatedAt,
		UpdatedAt:   dbJob.UpdatedAt,
	}

	if dbJob.ErrorMessage.Valid {
		job.Error = dbJob.ErrorMessage.String
	}

	if dbJob.CompletedAt.Valid {
		job.CompletedAt = &dbJob.CompletedAt.Time
	}

	if err := r.cacheClient.Set(ctx, cacheKey, &job); err != nil {
		logger.GetLogger().Error("failed to cache upload job", zap.String("cacheKey", cacheKey), zap.Error(err))
		return nil, fmt.Errorf("failed to cache upload job: %w", err)
	}
	return &job, nil
}

func (r *SQLUploadRepository) UpdateUploadJob(ctx context.Context, job *model.UploadJob) error {
	return nil
}

func (r *SQLUploadRepository) ListPendingJobs(ctx context.Context, limit int) (*[]model.UploadJob, error) {
	return nil, nil
}

func (r *SQLUploadRepository) ListFailedJobs(ctx context.Context, limit int) (*[]model.UploadJob, error) {
	return nil, nil
}

func (r *SQLUploadRepository) ListProcessingJobs(ctx context.Context, limit int) (*[]model.UploadJob, error) {
	return nil, nil
}

func (r *SQLUploadRepository) ListUserJobs(ctx context.Context, userID int64, limit, offset int) (*[]model.UploadJob, error) {
	return nil, nil
}

func (r *SQLUploadRepository) CreateFileMetadata(ctx context.Context, metadata *model.FileMetadata) error {
	return nil
}

func (r *SQLUploadRepository) GetFileMetadata(ctx context.Context, id int64) (*model.FileMetadata, error) {
	return nil, nil
}
