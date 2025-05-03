-- Migration to remove upload_jobs and related tables
-- db/migration/000003_add_upload_jobs.down.sql

-- Remove triggers first
DROP TRIGGER IF EXISTS trigger_update_upload_jobs_updated_at ON upload_jobs;
DROP TRIGGER IF EXISTS trigger_update_file_metadata_updated_at ON file_metadata;

-- Remove indexes
DROP INDEX IF EXISTS idx_upload_jobs_user_id;
DROP INDEX IF EXISTS idx_upload_jobs_status;
DROP INDEX IF EXISTS idx_upload_jobs_created_at;
DROP INDEX IF EXISTS idx_file_metadata_mime_type;

-- Remove file_metadata table (has foreign key to upload_jobs)
DROP TABLE IF EXISTS file_metadata;

-- Remove upload_jobs table
DROP TABLE IF EXISTS upload_jobs;