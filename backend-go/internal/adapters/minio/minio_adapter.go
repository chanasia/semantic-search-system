package minio

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/chanasia/semantic-search-system/internal/core/domain"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioAdapter struct {
	client     *minio.Client
	endpoint   string
	bucketName string
}

func NewMinioAdapter(endpoint, accessKey, secretKey string, useSSL bool, bucketName string) (*MinioAdapter, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	// Check if bucket exists
	exists, err := client.BucketExists(context.Background(), bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}

	// Create bucket if it doesn't exist
	if !exists {
		opts := minio.MakeBucketOptions{
			Region: "us-east-1",
		}
		if err := client.MakeBucket(context.Background(), bucketName, opts); err != nil {
			return nil, fmt.Errorf("failed to create bucket %s: %w", bucketName, err)
		}
	}

	return &MinioAdapter{
		client:     client,
		endpoint:   endpoint,
		bucketName: bucketName,
	}, nil
}

// ปรับ UploadFile ให้ใช้ default bucketName
func (m *MinioAdapter) UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (*domain.FileMetadata, error) {
	// Upload file using default bucket
	_, err := m.client.PutObject(ctx, m.bucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	return &domain.FileMetadata{
		Filename:    objectName,
		Size:        size,
		ContentType: contentType,
		Bucket:      m.bucketName,
		Path:        objectName,
		UploadedAt:  time.Now(),
	}, nil
}

// ปรับ GetFileURL ให้ใช้ default bucketName
func (m *MinioAdapter) GetFileURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	url, err := m.client.PresignedGetObject(ctx, m.bucketName, objectName, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url.String(), nil
}

// ปรับ DeleteFile ให้ใช้ default bucketName
func (m *MinioAdapter) DeleteFile(ctx context.Context, objectName string) error {
	return m.client.RemoveObject(ctx, m.bucketName, objectName, minio.RemoveObjectOptions{})
}

// ปรับ FileExists ให้ใช้ default bucketName (ถ้าต้องการใช้ในอนาคต)
func (m *MinioAdapter) FileExists(ctx context.Context, objectName string) (bool, error) {
	_, err := m.client.StatObject(ctx, m.bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		var errResponse minio.ErrorResponse
		if errors.As(err, &errResponse) {
			if errResponse.Code == "NoSuchKey" {
				return false, nil
			}
		}
		return false, err
	}
	return true, nil
}
