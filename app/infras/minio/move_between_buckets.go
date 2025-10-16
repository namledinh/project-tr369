package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
)

// MoveFileBetweenBuckets moves a file from one bucket to another (actually copy + delete)
func (s *MinioStore) MoveFileBetweenBuckets(srcBucket, srcPath, dstBucket, dstPath string) error {
	ctx := context.Background()

	// Copy from source bucket to destination bucket
	src := minio.CopySrcOptions{
		Bucket: srcBucket,
		Object: srcPath,
	}
	dst := minio.CopyDestOptions{
		Bucket: dstBucket,
		Object: dstPath,
	}

	_, err := s.Client.CopyObject(ctx, dst, src)
	if err != nil {
		return fmt.Errorf("failed to copy object from %s/%s to %s/%s: %v", srcBucket, srcPath, dstBucket, dstPath, err)
	}

	// Delete source file
	err = s.Client.RemoveObject(ctx, srcBucket, srcPath, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to remove source object %s/%s: %v", srcBucket, srcPath, err)
	}

	fmt.Printf("Moved successfully from %s/%s to %s/%s\n", srcBucket, srcPath, dstBucket, dstPath)
	return nil
}
