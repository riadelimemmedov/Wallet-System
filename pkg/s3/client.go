package s3

import (
	"context"
	"fmt"
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
		return fmt.Errorf("Failed to load AWS config: %2w", awsCfgErr)
	}

	s.client = s3.NewFromConfig(awsCfg)
	s.initialized = true

	if err := s.testConnection(); err != nil {
		s.initialized = false
		logger.GetLogger().Error("Failed to connect to S3", zap.Error(err))
		return fmt.Errorf("Failed to connect to S3: %w", err)
	}
	return nil
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
