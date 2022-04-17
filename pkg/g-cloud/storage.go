package g_cloud

import (
	"bytes"
	"context"
	"fmt"
	"io"

	gstorage "cloud.google.com/go/storage"
	goseidon "github.com/go-seidon/core"
	"github.com/go-seidon/core/internal/clock"
	"google.golang.org/api/option"
)

type GoogleConfig struct {
	ProjectId      string
	CredentialPath string
	BucketName     string
}

type GoogleStorage struct {
	Config *GoogleConfig
	Client GoogleStorageClient
	Clock  clock.Clock
}

func (s *GoogleStorage) UploadFile(ctx context.Context, p goseidon.UploadFileParam) (*goseidon.UploadFileResult, error) {
	if ctx == nil {
		return nil, fmt.Errorf("invalid context")
	}

	wc := s.Client.NewWriter(ctx, s.Config.BucketName, p.FileName)
	buf := bytes.NewBuffer(p.FileData)
	_, err := s.Client.Copy(wc, buf)
	if err != nil {
		return nil, err
	}

	err = wc.Close()
	if err != nil {
		return nil, err
	}

	uploadedAt := s.Clock.Now()
	res := &goseidon.UploadFileResult{
		FileName:   p.FileName,
		UploadedAt: uploadedAt,
	}
	return res, nil
}

func (s *GoogleStorage) RetrieveFile(ctx context.Context, p goseidon.RetrieveFileParam) (*goseidon.RetrieveFileResult, error) {
	if ctx == nil {
		return nil, fmt.Errorf("invalid context")
	}

	rc, err := s.Client.NewReader(ctx, s.Config.BucketName, p.Id)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	fileData, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	res := &goseidon.RetrieveFileResult{
		File: fileData,
	}
	return res, nil
}

func (s *GoogleStorage) DeleteFile(ctx context.Context, p goseidon.DeleteFileParam) (*goseidon.DeleteFileResult, error) {
	if ctx == nil {
		return nil, fmt.Errorf("invalid context")
	}

	err := s.Client.Delete(ctx, s.Config.BucketName, p.Id)
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

func NewGoogleConfig(projectId, credentialPath, bucketName string) (*GoogleConfig, error) {
	if projectId == "" {
		return nil, fmt.Errorf("invalid google project id")
	}
	if credentialPath == "" {
		return nil, fmt.Errorf("invalid google credential path")
	}
	if bucketName == "" {
		return nil, fmt.Errorf("invalid google bucket name")
	}
	c := &GoogleConfig{
		ProjectId:      projectId,
		CredentialPath: credentialPath,
		BucketName:     bucketName,
	}
	return c, nil
}

func NewGoogleStorage(config *GoogleConfig) (*GoogleStorage, error) {
	if config == nil {
		return nil, fmt.Errorf("invalid google config")
	}

	ctx := context.Background()
	gsClient, err := gstorage.NewClient(
		ctx, option.WithCredentialsFile(config.CredentialPath),
	)
	if err != nil {
		return nil, err
	}
	client := &googleStorageClient{
		client: gsClient,
	}

	clock, _ := clock.NewClock()

	s := &GoogleStorage{
		Config: config,
		Client: client,
		Clock:  clock,
	}
	return s, nil
}
