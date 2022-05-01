package g_storage

import (
	"bytes"
	"context"
	"fmt"
	"io"

	goseidon "github.com/go-seidon/core"
	"github.com/go-seidon/core/internal/clock"
	g_cloud "github.com/go-seidon/core/internal/g-cloud"
)

type GoogleStorage struct {
	Config *GoogleConfig
	Client g_cloud.GoogleStorageClient
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

func NewGoogleStorage(opt GoogleStorageOption) (*GoogleStorage, error) {
	if opt == nil {
		return nil, fmt.Errorf("invalid google option")
	}

	config := &GoogleConfig{}
	err := opt.Apply(config)
	if err != nil {
		return nil, err
	}

	client, _ := g_cloud.NewGoogleStorageClient(config.GoogleClient)
	clock, _ := clock.NewClock()

	s := &GoogleStorage{
		Config: config,
		Client: client,
		Clock:  clock,
	}
	return s, nil
}
