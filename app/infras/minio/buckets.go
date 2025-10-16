package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"time"
)

func (s *MinioStore) ListBuckets() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	buckets, err := s.Client.ListBuckets(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list buckets: %v", err)
	}

	var bucketNames []string
	for _, bucket := range buckets {
		bucketNames = append(bucketNames, bucket.Name)
	}

	return bucketNames, nil
}

func (s *MinioStore) EnsureBucket(bucketName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	exists, err := s.Client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to check if bucket %s exists: %v", bucketName, err)
	}
	if !exists {
		fmt.Printf("Creating bucket: %s\n", bucketName)
		err = s.Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket %s: %v", bucketName, err)
		}
		fmt.Printf("Created bucket: %s\n", bucketName)
	} else {
		fmt.Printf("Bucket already exists: %s\n", bucketName)
	}

	return nil
}

func (s *MinioStore) GetBucketInfo(bucketName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectCh := s.Client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Recursive: false,
		MaxKeys:   1,
	})

	hasError := false
	for object := range objectCh {
		if object.Err != nil {
			hasError = true
			return fmt.Errorf("failed to access bucket %s: %v", bucketName, object.Err)
		}
		fmt.Printf("Bucket %s is accessible\n", bucketName)
	}

	if !hasError {
		fmt.Printf("Bucket %s is accessible (empty)\n", bucketName)
	}
	return nil
}
