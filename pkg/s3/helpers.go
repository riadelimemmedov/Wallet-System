package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3_types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	logger "github.com/riad/banksystemendtoend/pkg/log"
	"go.uber.org/zap"
)

// createMultipartUpload initiates a multipart upload and returns the upload ID
func (s *S3Service) createMultipartUpload(ctx context.Context, key string, contentType string) (string, error) {
	createResp, err := s.client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
		Bucket:      aws.String(s.config.Bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		logger.Error("Failed to create multipart upload",
			zap.String("bucket", s.config.Bucket),
			zap.String("key", key),
			zap.Error(err))
		return "", fmt.Errorf("failed to create multipart upload: %w", err)
	}
	return *createResp.UploadId, nil
}

// uploadParts uploads file chunks as individual parts and returns completed parts
func (s *S3Service) uploadParts(ctx context.Context, key string, uploadID string, reader io.Reader) ([]s3_types.CompletedPart, error) {
	var parts []s3_types.CompletedPart
	var partNumber int32 = 1
	buffer := make([]byte, s.config.UploadChunkSize)

	for {
		n, err := io.ReadFull(reader, buffer)
		if err != nil && !isExpectedEOFError(err) {
			return nil, fmt.Errorf("failed to read file content: %w", err)
		}
		if n == 0 {
			break
		}

		part, err := s.uploadSinglePart(ctx, key, uploadID, partNumber, buffer[:n])
		if err != nil {
			return nil, fmt.Errorf("failed to upload part: %w", err)
		}

		parts = append(parts, part)
		partNumber++

		if n < int(s.config.UploadChunkSize) {
			break
		}
	}
	return parts, nil
}

// uploadSinglePart uploads a single part of the multipart upload
func (s *S3Service) uploadSinglePart(ctx context.Context, key string, uploadID string, partNumber int32, data []byte) (s3_types.CompletedPart, error) {
	partResp, err := s.client.UploadPart(ctx, &s3.UploadPartInput{
		Bucket:     aws.String(s.config.Bucket),
		Key:        aws.String(key),
		PartNumber: &partNumber,
		UploadId:   aws.String(uploadID),
		Body:       bytes.NewReader(data),
	})

	if err != nil {
		logger.Error("Failed to Single upload part",
			zap.String("bucket", s.config.Bucket),
			zap.String("key", key),
			zap.Int32("partNumber", partNumber),
			zap.Error(err))
		return s3_types.CompletedPart{}, nil
	}

	return s3_types.CompletedPart{
		ETag:       partResp.ETag,
		PartNumber: &partNumber,
	}, nil
}

func (s *S3Service) abortMultipartUpload(ctx context.Context, key string, uploadID string) {
	_, err := s.client.AbortMultipartUpload(ctx, &s3.AbortMultipartUploadInput{
		Bucket:   aws.String(s.config.Bucket),
		Key:      aws.String(key),
		UploadId: aws.String(uploadID),
	})

	if err != nil {
		logger.Error("Failed to abort multipart upload",
			zap.String("bucket", s.config.Bucket),
			zap.String("key", key),
			zap.String("uploadID", uploadID),
			zap.Error(err))
	}
}

// completeMultipartUpload finalizes a multipart upload with all uploaded parts
func (s *S3Service) completeMultipartUpload(ctx context.Context, key string, uploadID string, parts []s3_types.CompletedPart) error {
	_, err := s.client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(s.config.Bucket),
		Key:      aws.String(key),
		UploadId: aws.String(uploadID),
		MultipartUpload: &s3_types.CompletedMultipartUpload{
			Parts: parts,
		},
	})
	if err != nil {
		logger.Error("Failed to complete multipart upload",
			zap.String("bucket", s.config.Bucket),
			zap.String("key", key),
			zap.Error(err))
		return fmt.Errorf("failed to complete multipart upload: %w", err)
	}
	return nil
}

// generateFileURL creates the public URL for the uploaded file
func (s *S3Service) generateFileURL(key string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.config.Bucket, s.config.Region, key)
}

// isNotFoundError checks if an error is a "not found" error
func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "NotFound") || strings.Contains(err.Error(), "404")
}

// isExpectedEOFError checks if an error is an expected EOF condition
func isExpectedEOFError(err error) bool {
	return err == io.EOF || err == io.ErrUnexpectedEOF
}
