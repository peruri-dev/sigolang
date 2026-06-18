package storage

import (
	"context"
	"io"
	"time"
)

type ListObjectOptions struct {
	Recursive bool
}

type StorageFileInfo struct {
	LastModified time.Time
	Size         int64
	ContentType  string
}

type PresignedGetObjectParam struct {
	Filename           string
	Expire             time.Duration
	ContentDisposition string
}

type StorageObject interface {
	Close() error
	Read(b []byte) (n int, err error)
	ReadAt(p []byte, off int64) (n int, err error)
	Seek(offset int64, whence int) (int64, error)
	Stat() (*StorageFileInfo, error)
}

type UploadFileMetadata struct {
	FileName string `json:"file_name"`
	MimeType string `json:"mime_type"`
	FileSize int64  `json:"file_size"`
}

type IStorage interface {
	PutObject(ctx context.Context, filePath string, fileReader io.Reader, fileSrc UploadFileMetadata) error
	PutObjectFile(ctx context.Context, bucketName, filePath string, fileReader io.Reader, fileSrc UploadFileMetadata) error
	PresignedGetObject(ctx context.Context, filePath string, param *PresignedGetObjectParam) (string, error)
	UnsafePublicObject(filePath string) (string, error)
	DeleteObject(ctx context.Context, filePath string) error
	ListObjects(ctx context.Context, opt *ListObjectOptions) ([]string, error)
	GetObject(ctx context.Context, filePath string) (StorageObject, error)
	CreateBucket(ctx context.Context, name string) error
	BucketName() string
}
