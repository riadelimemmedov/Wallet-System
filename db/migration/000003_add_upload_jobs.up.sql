-- Migration to add upload_jobs table
-- db/migration/000003_add_upload_jobs.up.sql

-- Create upload status enum type
CREATE TYPE upload_status AS ENUM (
    'PENDING',
    'PROCESSING',
    'COMPLETED',
    'FAILED',
    'CANCELLED'
);

-- Create upload_jobs table
CREATE TABLE IF NOT EXISTS upload_jobs (
    id VARCHAR(36) PRIMARY KEY,
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    content_type VARCHAR(255) NOT NULL,
    temp_path TEXT NOT NULL,
    target_path TEXT NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(user_id),
    status upload_status NOT NULL DEFAULT 'PENDING',
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMPTZ
);

-- Create indexes for common queries
CREATE INDEX idx_upload_jobs_user_id ON upload_jobs(user_id);
CREATE INDEX idx_upload_jobs_status ON upload_jobs(status);
CREATE INDEX idx_upload_jobs_created_at ON upload_jobs(created_at);

-- Create file_metadata table for additional information about successfully uploaded files
CREATE TABLE IF NOT EXISTS file_metadata(
    id SERIAL PRIMARY KEY NOT NULL,
    upload_jobs_id VARCHAR(36) NOT NULL REFERENCES upload_jobs(id) ON DELETE CASCADE,
    s3_url TEXT NOT NULL,
    checksum VARCHAR(64) NOT NULL,
    mime_type VARCHAR(255) NOT NULL,
    width INTEGER,
    height INTEGER,
    additional_data JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create index on mime_type for filtering by file type
CREATE INDEX idx_file_metadata_mime_type ON file_metadata(mime_type);

-- Add triggers for updating the updated_at timestamp automatically
CREATE TRIGGER trigger_update_upload_jobs_updated_at
BEFORE UPDATE ON upload_jobs
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_update_file_metadata_updated_at
BEFORE UPDATE ON file_metadata
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();