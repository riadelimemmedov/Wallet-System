package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"go.uber.org/zap"

	logger "github.com/riad/banksystemendtoend/pkg/log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// ! S3ServiceConfig holds the configuration for the S3 service
type S3ServiceConfig struct {
	Region          string
	Bucket          string
	Prefix          string
	MaxRetries      int
	Timeout         time.Duration
	UploadChunkSize int64
}

// ! DefaultS3Config returns default configuration for S3 service
func DefaultS3Config() S3ServiceConfig {
	return S3ServiceConfig{
		Region:          "us-east-1",
		Bucket:          "default-config-bucket",
		Prefix:          "default-prefix",
		MaxRetries:      3,
		Timeout:         2 * time.Minute,
		UploadChunkSize: 5 * 1024 * 1024,
	}
}

// ! S3Service provides operations with AWS S3
type S3Service struct {
	client      *s3.Client
	config      S3ServiceConfig
	mu          sync.RWMutex
	initialized bool
}

var (
	instance *S3Service
	once     sync.Once
)

// ! GetS3Service returns a singleton instance of S3Service
func GetS3Service(cfg S3ServiceConfig) (*S3Service, error) {
	once.Do(func() {
		instance = &S3Service{
			config: cfg,
		}
	})

	if !instance.initialized {
		if err := instance.initialize(); err != nil {
			return nil, err
		}
	}
	logger.GetLogger().Info("S3 service initialized successfully")
	return instance, nil
}

// ! initialize creates the S3 client
func (s *S3Service) initialize() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.initialized {
		return nil
	}

	awsCfg, awsCfgErr := aws_config.LoadDefaultConfig(context.Background(),
		aws_config.WithRegion(s.config.Region),
		aws_config.WithRetryMode(aws.RetryModeStandard),
		aws_config.WithRetryMaxAttempts(s.config.MaxRetries),
	)

	if awsCfgErr != nil {
		logger.GetLogger().Error("Failed to load AWS config", zap.Error(awsCfgErr))
		return fmt.Errorf("failed to load AWS config: %2w", awsCfgErr)
	}

	s.client = s3.NewFromConfig(awsCfg)
	s.initialized = true

	if err := s.testConnection(); err != nil {
		s.initialized = false
		logger.GetLogger().Error("Failed to connect to S3", zap.Error(err))
		return fmt.Errorf("failed to connect to S3: %w", err)
	}
	return nil
}

// prepareRequest ensures service is initialized and formats the key with proper prefix
func (s *S3Service) prepareRequest(key string) (string, error) {
	if !s.initialized {
		if err := s.initialize(); err != nil {
			return "", fmt.Errorf("failed to initialize S3 service: %w", err)
		}
	}
	formattedKey := key
	if s.config.Prefix != "" && !strings.HasPrefix(key, s.config.Prefix) {
		formattedKey = s.config.Prefix + key
	}

	return formattedKey, nil
}

// ! testConnection checks if the S3 bucket is accessible
func (s *S3Service) testConnection() error {
	if s.config.Bucket == "" {
		logger.GetLogger().Error("Please set bucket before check connection")
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.config.Timeout)
	defer cancel()

	_, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(s.config.Bucket),
	})
	return err
}

// UploadFile uploads a file to s3
func (s *S3Service) UploadFile(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	processedKey, err := s.prepareRequest(key)
	if err != nil {
		return "", err
	}

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.config.Bucket),
		Key:         aws.String(processedKey),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})

	if err != nil {
		logger.Error("Failed to upload file to S3",
			zap.String("bucket", s.config.Bucket),
			zap.String("key", key),
			zap.Error(err))
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.config.Bucket, s.config.Region, key), nil
}

// !IMPLEMENT Later
func (s *S3Service) UploadLargeFile(ctx context.Context, key string, reader io.Reader, contentType string) (string, error) {
	return "", nil
}

// DownloadFile downloads a file from s3
func (s *S3Service) DownloadFile(ctx context.Context, key string) ([]byte, error) {
	processedKey, err := s.prepareRequest(key)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(processedKey),
	})

	if err != nil {
		logger.Error("Failed to download file from s3",
			zap.String("bucket", s.config.Bucket),
			zap.String("key", processedKey),
			zap.Error(err))
		return nil, fmt.Errorf("failed to download file from s3 %w", err)
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content %w", err)
	}
	return content, nil
}

// ! ListFiles lists files in a specific path
func (s *S3Service) ListFiles(ctx context.Context, prefix string) ([]string, error) {
	processedKey, err := s.prepareRequest(prefix)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.config.Bucket),
		Prefix: aws.String(processedKey),
	})

	if err != nil {
		logger.Error("failed to list files from s3",
			zap.String("bucket", s.config.Bucket),
			zap.String("prefix", processedKey),
			zap.Error(err))
		return nil, fmt.Errorf("failed to list files from s3 %w", err)
	}

	keys := make([]string, 0, len(resp.Contents))
	for _, obj := range resp.Contents {
		keys = append(keys, *obj.Key)
	}
	return keys, nil
}

// DeleteFile deletes a file from S3
func (s *S3Service) DeleteFile(ctx context.Context, key string) error {
	processedKey, err := s.prepareRequest(key)
	if err != nil {
		return err
	}

	_, err = s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(processedKey),
	})

	if err != nil {
		logger.Error("Failed to delete file from s3",
			zap.String("bucket", s.config.Bucket),
			zap.String("key", key),
			zap.Error(err))
		return fmt.Errorf("failed to delete file from s3 %w", err)
	}
	return nil
}

// FileExists checks if a file exists in S3
func (s *S3Service) FileExists(ctx context.Context, key string) (bool, error) {
	processedKey, err := s.prepareRequest(key)
	if err != nil {
		return false, err
	}

	_, err = s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(processedKey),
	})

	if err != nil {
		if strings.Contains(err.Error(), "NotFound") || strings.Contains(err.Error(), "404") {
			return false, nil
		}
		logger.Error("Failed to check if file exists",
			zap.String("bucket", s.config.Bucket),
			zap.String("key", key),
			zap.Error(err))
		return false, fmt.Errorf("failed to check if file exists %w", err)
	}

	return true, nil
}

// Close close s3 connection
func (s *S3Service) Close() error {
	return nil
}
