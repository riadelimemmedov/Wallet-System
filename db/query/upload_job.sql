-- name: CreateUploadJob :one
INSERT INTO upload_jobs (
    id,
    file_name,
    file_size,
    content_type,
    temp_path,
    target_path,
    user_id,
    status,
    error_message,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: GetUploadJob :one
SELECT * FROM upload_jobs
WHERE id = $1;

-- name: UpdateUploadJobStatus :one
UPDATE upload_jobs SET 
    status = $1,
    error_message = $2,
    updated_at = $3,
    completed_at = CASE WHEN $1 IN ('COMPLETED', 'FAILED') THEN $4 ELSE completed_at END
WHERE id = $5
RETURNING *;

-- name: ListCompletedUploadJobs :many
SELECT * FROM upload_jobs
WHERE status = 'COMPLETED'
ORDER BY completed_at DESC
LIMIT $1;

-- name: ListPendingUploadJobs :many
SELECT * FROM upload_jobs
WHERE status = 'PENDING'
ORDER BY created_at DESC
LIMIT $1;

-- name: ListProcessingUploadJobs :many
SELECT * FROM upload_jobs
WHERE status = 'PROCESSING'
ORDER BY created_at DESC
LIMIT $1;

-- name: ListFailedUploadJobs :many
SELECT * FROM upload_jobs
WHERE status = 'FAILED'
ORDER BY updated_at DESC
LIMIT $1;

-- name: ListUseUrploadJobs :many
SELECT * FROM upload_jobs
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: DeleteUploadJob :exec
DELETE FROM upload_jobs
WHERE id = $1;

-- name: CreateFileMetadata :one
INSERT INTO file_metadata (
    id,
    s3_url,
    checksum,
    mime_type,
    width,
    height,
    additional_data,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetFileMetadata :one
SELECT * FROM file_metadata
WHERE id = $1;

-- name: ListFilesByMimeType :many
SELECT fm.*
FROM file_metadata fm
JOIN upload_jobs uj ON fm.upload_jobs_id = uj.id
WHERE fm.mime_type LIKE $1 || '%'
AND uj.user_id = $2
ORDER BY fm.created_at DESC
LIMIT $3 OFFSET $4;

-- name: CountUserUploads :one
SELECT COUNT(*) FROM upload_jobs
WHERE user_id = $1 AND status = 'COMPLETED';