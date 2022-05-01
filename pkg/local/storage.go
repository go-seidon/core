package local

import (
	"context"
	"fmt"
	"io/fs"

	goseidon "github.com/go-seidon/core"
	"github.com/go-seidon/core/internal/clock"
	"github.com/go-seidon/core/internal/io"
)

type LocalStorage struct {
	Config *LocalConfig
	Client io.FileManager
	Clock  clock.Clock
}

func (s *LocalStorage) UploadFile(ctx context.Context, p goseidon.UploadFileParam) (*goseidon.UploadFileResult, error) {
	if ctx == nil {
		return nil, fmt.Errorf("invalid context")
	}

	rwPermission := fs.FileMode(0644)

	if !s.Client.IsExists(s.Config.StorageDir) {
		err := s.Client.CreateDir(s.Config.StorageDir, rwPermission)
		if err != nil {
			return nil, fmt.Errorf("failed create storage dir: %s", s.Config.StorageDir)
		}
	}

	path := fmt.Sprintf("%s/%s", s.Config.StorageDir, p.FileId)
	if s.Client.IsExists(path) {
		return nil, fmt.Errorf("file already exists")
	}

	err := s.Client.WriteFile(path, p.FileData, rwPermission)
	if err != nil {
		return nil, fmt.Errorf("failed storing file")
	}

	uploadedAt := s.Clock.Now()
	res := &goseidon.UploadFileResult{
		FileId:     p.FileId,
		FileName:   p.FileName,
		UploadedAt: uploadedAt,
	}
	return res, nil
}

func (s *LocalStorage) RetrieveFile(ctx context.Context, p goseidon.RetrieveFileParam) (*goseidon.RetrieveFileResult, error) {
	if ctx == nil {
		return nil, fmt.Errorf("invalid context")
	}

	path := fmt.Sprintf("%s/%s", s.Config.StorageDir, p.Id)
	if !s.Client.IsExists(path) {
		return nil, fmt.Errorf("file is not found")
	}

	file, err := s.Client.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed open file")
	}
	defer file.Close()

	binFile, err := s.Client.ReadFile(file)
	if err != nil {
		return nil, err
	}

	retrievedAt := s.Clock.Now()
	res := &goseidon.RetrieveFileResult{
		File:        binFile,
		RetrievedAt: retrievedAt,
	}
	return res, nil
}

func (s *LocalStorage) DeleteFile(ctx context.Context, p goseidon.DeleteFileParam) (*goseidon.DeleteFileResult, error) {
	if ctx == nil {
		return nil, fmt.Errorf("invalid context")
	}

	path := fmt.Sprintf("%s/%s", s.Config.StorageDir, p.Id)
	if !s.Client.IsExists(path) {
		return nil, fmt.Errorf("file is not found")
	}

	err := s.Client.RemoveFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed delete file")
	}

	deletedAt := s.Clock.Now()
	res := &goseidon.DeleteFileResult{
		Id:        p.Id,
		DeletedAt: deletedAt,
	}
	return res, nil
}

func NewLocalStorage(opt LocalStorageOption) (*LocalStorage, error) {
	if opt == nil {
		return nil, fmt.Errorf("invalid storage option")
	}

	cfg := &LocalConfig{}
	err := opt.Apply(cfg)
	if err != nil {
		return nil, err
	}

	client, _ := io.NewFileManager()
	clock, _ := clock.NewClock()
	s := &LocalStorage{
		Config: cfg,
		Client: client,
		Clock:  clock,
	}
	return s, nil
}
