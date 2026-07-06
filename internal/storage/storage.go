package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type ObjectInfo struct {
	Path       string    `json:"path"`
	Name       string    `json:"name"`
	Size       int64     `json:"size"`
	ModifiedAt time.Time `json:"modified_at"`
}

type Storage interface {
	Put(ctx context.Context, path string, reader io.Reader) error
	Get(ctx context.Context, path string) (io.ReadCloser, error)
	List(ctx context.Context, prefix string) ([]ObjectInfo, error)
	Delete(ctx context.Context, path string) error
}

type LocalStorage struct {
	baseDir string
}

func NewLocalStorage(baseDir string) (*LocalStorage, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, err
	}
	return &LocalStorage{baseDir: baseDir}, nil
}

func (l *LocalStorage) Put(ctx context.Context, path string, reader io.Reader) error {
	fullPath := filepath.Join(l.baseDir, path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}
	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, reader)
	return err
}

func (l *LocalStorage) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(l.baseDir, path)
	return os.Open(fullPath)
}

func (l *LocalStorage) List(ctx context.Context, prefix string) ([]ObjectInfo, error) {
	var objects []ObjectInfo
	searchDir := filepath.Join(l.baseDir, prefix)
	err := filepath.WalkDir(searchDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		if !d.IsDir() {
			relPath, err := filepath.Rel(l.baseDir, path)
			if err != nil {
				return err
			}
			info, err := d.Info()
			if err != nil {
				return err
			}
			objects = append(objects, ObjectInfo{
				Path:       relPath,
				Name:       d.Name(),
				Size:       info.Size(),
				ModifiedAt: info.ModTime(),
			})
		}
		return nil
	})
	return objects, err
}

func (l *LocalStorage) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(l.baseDir, path)
	return os.RemoveAll(fullPath)
}

type S3Storage struct {
	client *s3.Client
	bucket string
}

func NewS3Storage(bucket, region, endpoint, accessKey, secretKey string, useSSL bool) (*S3Storage, error) {
	var opts []func(*config.LoadOptions) error

	if region != "" {
		opts = append(opts, config.WithRegion(region))
	}

	if accessKey != "" && secretKey != "" {
		opts = append(opts, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		))
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load S3 config: %w", err)
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
			o.UsePathStyle = true
		}
	})

	return &S3Storage{
		client: s3Client,
		bucket: bucket,
	}, nil
}

func (s *S3Storage) Put(ctx context.Context, path string, reader io.Reader) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
		Body:   strings.NewReader(string(data)),
	})
	return err
}

func (s *S3Storage) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}
	return out.Body, nil
}

func (s *S3Storage) List(ctx context.Context, prefix string) ([]ObjectInfo, error) {
	var objects []ObjectInfo
	paginator := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(prefix),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, obj := range page.Contents {
			key := aws.ToString(obj.Key)
			name := key
			if idx := strings.LastIndex(key, "/"); idx != -1 {
				name = key[idx+1:]
			}
			objects = append(objects, ObjectInfo{
				Path:       key,
				Name:       name,
				Size:       aws.ToInt64(obj.Size),
				ModifiedAt: aws.ToTime(obj.LastModified),
			})
		}
	}

	return objects, nil
}

func (s *S3Storage) Delete(ctx context.Context, path string) error {
	paginator := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(path),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		for _, obj := range page.Contents {
			_, err = s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
				Bucket: aws.String(s.bucket),
				Key:    obj.Key,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
