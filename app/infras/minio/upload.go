package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"io"
	"strings"
)

// UploadFile upload file local lên MinIO
func (s *MinioStore) UploadFile(bucketName, objectName string, reader io.Reader, size int64) error {
	ctx := context.Background()
	info, err := s.Client.PutObject(
		ctx,
		bucketName,
		objectName,
		reader,
		size,
		minio.PutObjectOptions{},
	)
	if err != nil {
		fmt.Printf("Direct upload failed: %v\n", err)
		if strings.Contains(err.Error(), "bucket does not exist") || strings.Contains(err.Error(), "NoSuchBucket") {
			fmt.Printf("Trying to create bucket %s\n", bucketName)
			if createErr := s.EnsureBucket(bucketName); createErr != nil {
				return fmt.Errorf("failed to create bucket and upload file: create bucket error: %v, upload error: %v", createErr, err)
			}

			fmt.Printf("Retrying upload after creating bucket\n")
			info, err = s.Client.PutObject(
				ctx,
				bucketName,
				objectName,
				reader,
				size,
				minio.PutObjectOptions{},
			)
			if err != nil {
				return fmt.Errorf("failed to upload file to bucket %s after creating bucket: %v", bucketName, err)
			}
		} else {
			return fmt.Errorf("failed to upload file to bucket %s: %v", bucketName, err)
		}
	}
	fmt.Printf("Uploaded %s (%d bytes)\n", objectName, info.Size)
	return nil
}
