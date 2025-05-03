package pkg_repository

import (
	"context"

	db "github.com/riad/banksystemendtoend/db/sqlc"
	"github.com/riad/banksystemendtoend/pkg/cache"
	pkg_interface "github.com/riad/banksystemendtoend/pkg/interface"
	"github.com/riad/banksystemendtoend/pkg/model"
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
	return nil
}

// GetUploadJob retrieves an upload job by its ID
func (r *SQLUploadRepository) GetUploadJob(ctx context.Context, id int64) (*model.UploadJob, error) {
	return nil, nil
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
