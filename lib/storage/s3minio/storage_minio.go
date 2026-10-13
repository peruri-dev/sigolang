package s3minio

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"sigolang/config"
	"sigolang/lib/storage"
	"slices"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	// [OTEL]
	// "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type MinioConnectCfg struct {
	Endpoint   string
	AccessKey  string
	SecretKey  string
	UseSSL     bool
	BucketName string
}

var zeroDuration time.Duration

type Storage struct {
	Impl any

	bucketName string
}

func parseURL(uri string) (*MinioConnectCfg, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	cfg := &MinioConnectCfg{
		Endpoint: u.Host,
		UseSSL:   true,
	}
	if user := u.User; user != nil {
		cfg.AccessKey = user.Username()
		if pass, ok := user.Password(); ok {
			cfg.SecretKey = pass
		}
	}
	for k, v := range u.Query() {
		switch k {
		case "sslmode":
			if len(v) == 0 {
				return nil, fmt.Errorf("minio config missing sslmode value")
			}
			val := strings.ToLower(v[0])
			valids := []string{"disable", "enable"}
			if !slices.Contains(valids, val) {
				return nil, fmt.Errorf("minio config invalid sslmode %s", val)
			}
			if val == valids[0] {
				cfg.UseSSL = false
			}
		default:
			return nil, fmt.Errorf("minio config unknown query params %s", k)
		}
	}
	bucketName := strings.TrimPrefix(u.Path, "/")
	if bucketName == "" {
		return nil, fmt.Errorf("minio config missing bucket name")
	}
	cfg.BucketName = bucketName

	return cfg, nil
}

type customWriter struct {
}

func (w *customWriter) Write(p []byte) (n int, err error) {
	slog.Info(string(p))
	return len(p), nil
}

func RegisterMinioStorage() {
	storage.RegisterStorage(&storage.StorageFactory{
		Prefixes: []string{"s3://"},
		Create: func(c *config.StorageConfig) (s storage.IStorage, err error) {
			appCfg := config.Get()
			cfg, err := parseURL(c.StorageUri)
			if err != nil {
				return nil, err
			}

			cloudStorage := &Storage{}

			var transport http.RoundTripper
			transport, err = minio.DefaultTransport(cfg.UseSSL)
			if err != nil {
				return nil, err
			}

			// [OTEL]
			// transport = otelhttp.NewTransport(transport)

			minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
				Creds:     credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
				Transport: transport,
				Secure:    cfg.UseSSL,
			})
			if err != nil {
				return nil, err
			}
			slog.Info("MinIO client created, trying to hit...", slog.String("bucket", cfg.BucketName))
			cw := customWriter{}

			if appCfg.MinioInitLog {
				minioClient.TraceOn(&cw)
				defer minioClient.TraceOff()
			}

			// Use a context with a timeout for the connection check
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Check connectivity by trying to list a (possibly non-existent) bucket.
			// The specific bucket name does not matter; we are checking if the server
			// responds.
			found, err := minioClient.BucketExists(ctx, cfg.BucketName)
			if err != nil {
				// This is where you handle connection errors (e.g., network issues,
				// incorrect credentials, server not running)
				return nil, fmt.Errorf("MinIO connection failed: %v", err)
			}

			// _, err = minioClient.ListBuckets(context.TODO())
			// if err != nil {
			// 	return nil, err
			// }

			if found {
				cloudStorage.Impl = minioClient
				cloudStorage.bucketName = cfg.BucketName

				slog.Info("minio connected", slog.String("Endpoint", cfg.Endpoint), slog.String("bucket", cloudStorage.bucketName))

				return cloudStorage, err
			}

			return nil, fmt.Errorf("MinIO server is connected, but the bucket does not exist")
		},
	})
}

func GetMC(sto *Storage) *minio.Client {
	if mc, ok := sto.Impl.(*minio.Client); ok {
		return mc
	}
	return nil
}

func (s *Storage) BucketName() string {
	return s.bucketName
}

func (s *Storage) CreateBucket(ctx context.Context, name string) error {
	mc := GetMC(s)
	if mc == nil {
		return fmt.Errorf("storage is not minio")
	}

	ok, err := mc.BucketExists(ctx, name)
	if err != nil {
		return err
	}
	if !ok {
		err = mc.MakeBucket(ctx, name, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) PutObject(ctx context.Context, filePath string, fileReader io.Reader, fileSrc storage.UploadFileMetadata) error {
	return s.PutObjectFile(ctx, s.BucketName(), filePath, fileReader, fileSrc)
}

func (s *Storage) PutObjectFile(ctx context.Context, bucketName, filePath string, fileReader io.Reader, fileSrc storage.UploadFileMetadata) error {
	mc := GetMC(s)
	if mc == nil {
		return fmt.Errorf("storage is not minio")
	}

	// If no bucketName provided, use default one
	if bucketName == "" {
		bucketName = s.BucketName()
	}

	filePath = strings.TrimPrefix(filePath, "/")

	// Ensure bucket exists
	err := s.CreateBucket(ctx, bucketName)
	if err != nil {
		return errors.Wrap(err, "unable to check or create bucket")
	}

	slog.InfoContext(ctx, "try to minio object uploaded",
		slog.String("bucket", bucketName),
		slog.String("filePath", filePath),
		slog.String("contentType", fileSrc.MimeType),
		slog.Int64("fileSize", fileSrc.FileSize),
	)

	_, err = mc.PutObject(ctx, bucketName, filePath, fileReader, fileSrc.FileSize, minio.PutObjectOptions{
		ContentType: fileSrc.MimeType,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to put object %s in bucket %s", filePath, bucketName)
	}

	return nil
}

func (s *Storage) PresignedGetObject(ctx context.Context, filePath string, param *storage.PresignedGetObjectParam) (string, error) {
	if strings.HasPrefix(filePath, "http://") || strings.HasPrefix(filePath, "https://") {
		return filePath, nil
	}

	filePath = strings.TrimPrefix(filePath, "/")

	mc := GetMC(s)
	if mc == nil {
		return "", fmt.Errorf("storage is not minio")
	}

	bucketName := s.BucketName()

	reqParams := make(url.Values)
	expires := 24 * time.Hour
	if param != nil {
		defaultCD := "attachment"
		if param.ContentDisposition == "" {
			param.ContentDisposition = defaultCD
		} else {
			allowedCD := []string{"attachment", "inline"}
			if !slices.Contains(allowedCD, param.ContentDisposition) {
				slog.WarnContext(ctx, "invalid content-disposition for presigning url", slog.String("value", param.ContentDisposition))
				param.ContentDisposition = defaultCD
			}
		}
		if param.Expire != zeroDuration {
			expires = param.Expire
		}
		if param.Filename != "" {
			safeName := url.PathEscape(param.Filename)
			reqParams.Add("response-content-disposition", fmt.Sprintf("%s; filename=\"%s\"", param.ContentDisposition, safeName))
		}
	}
	presignedURL, err := mc.PresignedGetObject(ctx, bucketName, filePath, expires, reqParams)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

func (s *Storage) UnsafePublicObject(filePath string) (string, error) {
	mc := GetMC(s)
	if mc == nil {
		return "", fmt.Errorf("storage is not minio")
	}

	bucketName := s.BucketName()
	endpoint := mc.EndpointURL()

	publicURL := fmt.Sprintf("%s/%s/%s", endpoint, bucketName, filePath)
	return publicURL, nil
}

func (s *Storage) DeleteObject(ctx context.Context, filePath string) error {
	mc := GetMC(s)
	if mc == nil {
		return fmt.Errorf("storage is not minio")
	}

	bucketName := s.BucketName()
	filePath = strings.TrimPrefix(filePath, "/")

	err := mc.RemoveObject(ctx, bucketName, filePath, minio.RemoveObjectOptions{})
	if err != nil {
		return errors.Wrapf(err, "failed to delete object %s from bucket %s", filePath, bucketName)
	}

	slog.InfoContext(ctx, "minio object deleted", slog.String("bucket", bucketName), slog.String("filePath", filePath))
	return nil
}

func (s *Storage) ListObjects(ctx context.Context, opt *storage.ListObjectOptions) ([]string, error) {
	mc := GetMC(s)
	if mc == nil {
		return nil, fmt.Errorf("storage is not minio")
	}

	bucketName := s.BucketName()

	objects := make([]string, 0)
	optMinio := minio.ListObjectsOptions{}
	if opt != nil {
		optMinio.Recursive = opt.Recursive
	}

	for object := range mc.ListObjects(ctx, bucketName, optMinio) {
		if object.Err != nil {
			return nil, object.Err
		}
		slog.InfoContext(ctx, "minio object listed", slog.String("bucket", bucketName), slog.String("filePath", object.Key))
		objects = append(objects, object.Key)
	}

	return objects, nil
}

type minioObject struct {
	Object *minio.Object
}

func (obj *minioObject) Close() error {
	return obj.Object.Close()
}
func (obj *minioObject) Read(b []byte) (n int, err error) {
	return obj.Object.Read(b)
}
func (obj *minioObject) ReadAt(p []byte, off int64) (n int, err error) {
	return obj.Object.ReadAt(p, off)
}
func (obj *minioObject) Seek(offset int64, whence int) (int64, error) {
	return obj.Object.Seek(offset, whence)
}
func (obj *minioObject) Stat() (*storage.StorageFileInfo, error) {
	stat, err := obj.Object.Stat()
	if err != nil {
		return nil, err
	}

	return &storage.StorageFileInfo{
		LastModified: stat.LastModified,
		Size:         stat.Size,
		ContentType:  stat.ContentType,
	}, nil
}

func (s *Storage) GetObject(ctx context.Context, filePath string) (storage.StorageObject, error) {
	mc := GetMC(s)
	if mc == nil {
		return nil, fmt.Errorf("storage is not minio")
	}

	bucketName := s.BucketName()
	object, err := mc.GetObject(ctx, bucketName, filePath, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	obj := &minioObject{
		Object: object,
	}
	slog.InfoContext(ctx, "minio object downloaded", slog.String("bucket", bucketName), slog.String("filePath", filePath))

	return obj, nil
}
