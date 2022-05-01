package aws_s3

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	goseidon "github.com/go-seidon/core"
	"github.com/go-seidon/core/internal/clock"
)

type AwsS3Storage struct {
	Config *AwsS3Config
	Client AwsS3Client
	Clock  clock.Clock
}

type AwsS3Client interface {
	PutObject(*s3.PutObjectInput) (*s3.PutObjectOutput, error)
	GetObject(*s3.GetObjectInput) (*s3.GetObjectOutput, error)
	DeleteObject(*s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error)
}

func (s *AwsS3Storage) UploadFile(ctx context.Context, p goseidon.UploadFileParam) (*goseidon.UploadFileResult, error) {
	if ctx == nil {
		return nil, fmt.Errorf("invalid context")
	}

	_, err := s.Client.PutObject(&s3.PutObjectInput{
		Body:   bytes.NewReader(p.FileData),
		Bucket: aws.String(s.Config.BucketName),
		Key:    aws.String(p.FileId),
	})
	if err != nil {
		return nil, err
	}

	uploadedAt := s.Clock.Now()
	res := &goseidon.UploadFileResult{
		FileId:     p.FileId,
		FileName:   p.FileName,
		UploadedAt: uploadedAt,
	}
	return res, nil
}

func (s *AwsS3Storage) RetrieveFile(ctx context.Context, p goseidon.RetrieveFileParam) (*goseidon.RetrieveFileResult, error) {
	if ctx == nil {
		return nil, fmt.Errorf("invalid context")
	}

	out, err := s.Client.GetObject(&s3.GetObjectInput{
		Key:    aws.String(p.Id),
		Bucket: aws.String(s.Config.BucketName),
	})
	if err != nil {
		return nil, err
	}

	fileData, err := io.ReadAll(out.Body)
	if err != nil {
		return nil, err
	}

	retrievedAt := s.Clock.Now()
	res := &goseidon.RetrieveFileResult{
		File:        fileData,
		RetrievedAt: retrievedAt,
	}
	return res, nil
}

func (s *AwsS3Storage) DeleteFile(ctx context.Context, p goseidon.DeleteFileParam) (*goseidon.DeleteFileResult, error) {
	if ctx == nil {
		return nil, fmt.Errorf("invalid context")
	}

	_, err := s.Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.Config.BucketName),
		Key:    aws.String(p.Id),
	})
	if err != nil {
		return nil, err
	}

	deletedAt := s.Clock.Now()
	res := &goseidon.DeleteFileResult{
		Id:        p.Id,
		DeletedAt: deletedAt,
	}
	return res, nil
}

func NewAwsS3Storage(opt AwsS3StorageOption) (*AwsS3Storage, error) {
	if opt == nil {
		return nil, fmt.Errorf("invalid aws s3 option")
	}

	cfg := &AwsS3Config{}
	err := opt.Apply(cfg)
	if err != nil {
		return nil, err
	}

	clock, _ := clock.NewClock()

	storage := &AwsS3Storage{
		Client: cfg.Client,
		Clock:  clock,
		Config: cfg,
	}
	return storage, nil
}
