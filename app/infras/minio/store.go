package minio

import (
	"github.com/minio/minio-go/v7"
)

type MinioStore struct {
	Client *minio.Client
}

func NewMinioStore(Client *minio.Client) *MinioStore {
	return &MinioStore{
		Client: Client,
	}
}
