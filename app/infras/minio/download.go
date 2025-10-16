package minio

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
)

func (s *MinioStore) GetFileURL(bucketName, objectName string, expiry time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	objInfo, err := s.Client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("object %s not found in bucket %s: %v", objectName, bucketName, err)
	}
	fmt.Printf("Object found: %s, size: %d bytes\n", objectName, objInfo.Size)

	reqParams := make(url.Values)
	presignedURL, err := s.Client.PresignedGetObject(ctx, bucketName, objectName, expiry, reqParams)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned url for %s in bucket %s: %v", objectName, bucketName, err)
	}

	urlString := presignedURL.String()
	fmt.Printf("Generated presigned URL: %s\n", urlString)
	fmt.Printf("URL expires in: %v\n", expiry)

	return urlString, nil
}

func (s *MinioStore) DownloadFile(bucketName, objectName, localPath string) error {
	ctx := context.Background()
	err := s.Client.FGetObject(ctx, bucketName, objectName, localPath, minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}
	fmt.Printf("Downloaded %s -> %s\n", objectName, localPath)
	return nil
}
