package pkg_repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
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

// ! CreateUploadJob stores a new upload job in the database
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
	if err := r.updateCaches(ctx, job); err != nil {
		logger.GetLogger().Error("failed to update caches",
			zap.String("id", job.ID),
			zap.Error(err))
	}
	return nil
}

// ! GetUploadJob retrieves an upload job by its ID
func (r *SQLUploadRepository) GetUploadJob(ctx context.Context, id int64) (*model.UploadJob, error) {
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

	// Cache the job
	if err := r.updateCaches(ctx, &job); err != nil {
		logger.GetLogger().Error("failed to update caches",
			zap.String("id", job.ID),
			zap.Error(err))
	}
	return &job, nil
}

// ! UpdateUploadJob updates an existing upload job
func (r *SQLUploadRepository) UpdateUploadJob(ctx context.Context, job *model.UploadJob) error {
	sqlStore, err := db.GetSQLStore(r.store)
	if err != nil {
		return fmt.Errorf("failed to get SQL store: %w", err)
	}

	tx, err := sqlStore.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		logger.GetLogger().Error("failed to begin transaction for update upload job", zap.Error(err))
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	status, err := pkg_utils.ConvertModelStatusToDBStatus(job.Status)
	if err != nil {
		return fmt.Errorf("failed to convert model status to DB status: %w", err)
	}

	job.UpdatedAt = time.Now()

	var completedAt sql.NullTime
	if (job.Status == model.UploadStatusCompleted || job.Status == model.UploadStatusFailed) && job.CompletedAt == nil {
		now := time.Now()
		job.CompletedAt = &now
		completedAt = sql.NullTime{Time: now, Valid: true}
	} else if job.CompletedAt != nil {
		completedAt = sql.NullTime{Time: *job.CompletedAt, Valid: true}
	}

	arg := db.UpdateUploadJobStatusParams{
		ID:           job.ID,
		Status:       status,
		ErrorMessage: sql.NullString{String: job.Error, Valid: job.Error != ""},
		UpdatedAt:    job.UpdatedAt,
		CompletedAt:  completedAt,
	}

	_, err = sqlStore.UpdateUploadJobStatus(ctx, arg)
	if err != nil {
		logger.GetLogger().Error("failed to update upload job", zap.Error(err))
		return fmt.Errorf("failed to update upload job: %w", err)
	}

	// Cache the job
	if err := r.updateCaches(ctx, job); err != nil {
		logger.GetLogger().Error("failed to update caches",
			zap.String("id", job.ID),
			zap.Error(err))
	}
	return nil
}

// ! ListPendingJobs retrieves a list of pending upload jobs
func (r *SQLUploadRepository) ListPendingJobs(ctx context.Context, limit int) (*[]model.UploadJob, error) {
	sqlStore, err := db.GetSQLStore(r.store)
	if err != nil {
		return nil, fmt.Errorf("failed to get SQL store: %w", err)
	}

	dbJobs, err := sqlStore.ListPendingUploadJobs(ctx, int32(limit))
	if err != nil {
		logger.GetLogger().Error("failed to list pending upload jobs", zap.Error(err))
		return nil, fmt.Errorf("failed to list pending upload jobs: %w", err)
	}

	jobs := make([]model.UploadJob, len(dbJobs))
	for i, dbJob := range dbJobs {
		jobs[i] = model.UploadJob{
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
			jobs[i].Error = dbJob.ErrorMessage.String
		}
		if dbJob.CompletedAt.Valid {
			completedAt := dbJob.CompletedAt.Time
			jobs[i].CompletedAt = &completedAt
		}
	}
	return &jobs, nil

}

// ! ListFailedJobs retrieves a list of failed upload jobs
func (r *SQLUploadRepository) ListFailedJobs(ctx context.Context, limit int) (*[]model.UploadJob, error) {
	sqlStore, err := db.GetSQLStore(r.store)
	if err != nil {
		return nil, fmt.Errorf("failed to get SQL store: %w", err)
	}

	dbJobs, err := sqlStore.ListFailedUploadJobs(ctx, int32(limit))
	if err != nil {
		logger.GetLogger().Error("failed to list failed upload jobs", zap.Error(err))
		return nil, fmt.Errorf("failed to list failed upload jobs: %w", err)
	}

	jobs := make([]model.UploadJob, len(dbJobs))
	for i, dbJob := range dbJobs {
		jobs[i] = model.UploadJob{
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
			jobs[i].Error = dbJob.ErrorMessage.String
		}
		if dbJob.CompletedAt.Valid {
			completedAt := dbJob.CompletedAt.Time
			jobs[i].CompletedAt = &completedAt
		}
	}
	return &jobs, nil
}

// ! ListProcessingJobs retrieves a list of jobs that are currently being processed
func (r *SQLUploadRepository) ListProcessingJobs(ctx context.Context, limit int) (*[]model.UploadJob, error) {
	sqlStore, err := db.GetSQLStore(r.store)
	if err != nil {
		return nil, fmt.Errorf("failed to get SQL store: %w", err)
	}

	dbJobs, err := sqlStore.ListProcessingUploadJobs(ctx, int32(limit))
	if err != nil {
		return nil, fmt.Errorf("failed to list processing upload jobs: %w", err)
	}

	jobs := make([]model.UploadJob, len(dbJobs))
	for i, dbJob := range dbJobs {
		jobs[i] = model.UploadJob{
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
			jobs[i].Error = dbJob.ErrorMessage.String
		}

		if dbJob.CompletedAt.Valid {
			completedAt := dbJob.CompletedAt.Time
			jobs[i].CompletedAt = &completedAt
		}
	}
	return &jobs, nil
}

// ! ListUserJobs retrieves a list of upload jobs for a specific user with pagination
func (r *SQLUploadRepository) ListUserJobs(ctx context.Context, userID int64, limit, offset int) (*[]model.UploadJob, error) {
	sqlStore, err := db.GetSQLStore(r.store)
	if err != nil {
		return nil, fmt.Errorf("failed to get SQL store: %w", err)
	}

	dbJobs, err := sqlStore.ListUseUrploadJobs(ctx, db.ListUseUrploadJobsParams{
		UserID: int32(userID),
		Limit:  int32(limit),
		Offset: int32(offset),
	})

	if err != nil {
		logger.GetLogger().Error("failed to list user upload jobs", zap.Error(err))
		return nil, fmt.Errorf("failed to list user upload jobs: %w", err)
	}

	jobs := make([]model.UploadJob, len(dbJobs))
	for i, dbJob := range dbJobs {
		jobs[i] = model.UploadJob{
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
			jobs[i].Error = dbJob.ErrorMessage.String
		}

		if dbJob.CompletedAt.Valid {
			completedAt := dbJob.CompletedAt.Time
			jobs[i].CompletedAt = &completedAt
		}
	}
	return &jobs, nil
}

// ! CreateFileMetadata stores metadata for a successfully uploaded file
func (r *SQLUploadRepository) CreateFileMetadata(ctx context.Context, metadata *model.FileMetadata) error {
	sqlStore, err := db.GetSQLStore(r.store)
	if err != nil {
		return fmt.Errorf("failed to get SQL store: %w", err)
	}

	arg := db.CreateFileMetadataParams{
		ID:             metadata.ID,
		S3Url:          metadata.S3URL,
		Checksum:       metadata.Checksum,
		MimeType:       metadata.MimeType,
		Width:          sql.NullInt32{Int32: int32(metadata.Width), Valid: metadata.Width > 0},
		Height:         sql.NullInt32{Int32: int32(metadata.Height), Valid: metadata.Height > 0},
		AdditionalData: pkg_utils.BytesToPgJSONB(metadata.AdditionalData),
		CreatedAt:      metadata.CreatedAt,
		UpdatedAt:      metadata.UpdatedAt,
	}
	_, err = sqlStore.CreateFileMetadata(ctx, arg)

	if err != nil {
		logger.GetLogger().Error("failed to create file metadata", zap.Error(err))
		return fmt.Errorf("failed to create file metadata: %w", err)
	}
	return nil
}

// ! GetFileMetadata retrieves file metadata by ID
func (r *SQLUploadRepository) GetFileMetadata(ctx context.Context, id int64) (*model.FileMetadata, error) {
	var metadata model.FileMetadata

	sqlStore, err := db.GetSQLStore(r.store)
	if err != nil {
		return nil, fmt.Errorf("failed to get SQL store: %w", err)
	}

	dbMetadata, err := sqlStore.GetFileMetadata(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("file metadata not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get file metadata: %w", err)
	}

	metadata = model.FileMetadata{
		ID:        dbMetadata.ID,
		S3URL:     dbMetadata.S3Url,
		MimeType:  dbMetadata.MimeType,
		CreatedAt: dbMetadata.CreatedAt,
		UpdatedAt: dbMetadata.UpdatedAt,
	}

	if dbMetadata.Width.Valid {
		metadata.Width = int(dbMetadata.Width.Int32)
	}
	if dbMetadata.Height.Valid {
		metadata.Height = int(dbMetadata.Height.Int32)
	}
	return &metadata, nil
}

// ! updateCaches updates the cache for the upload job
func (r *SQLUploadRepository) updateCaches(ctx context.Context, job *model.UploadJob) error {
	cacheKey := fmt.Sprintf(uploadJobCacheKey, job.ID)
	if err := r.cacheClient.Set(ctx, cacheKey, job); err != nil {
		logger.GetLogger().Error("failed to update cache for upload job", zap.String("cacheKey", cacheKey), zap.Error(err))
		return fmt.Errorf("failed to update cache for upload job: %w", err)
	}
	r.invalidateListCaches(ctx)
	return nil
}

// ! invalidateListCaches invalidates all list caches when a job status changes
func (r *SQLUploadRepository) invalidateListCaches(ctx context.Context) {
	r.cacheClient.DeleteByPattern(ctx, "upload_jobs:*")
}
