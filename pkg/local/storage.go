package local

import (
	"fmt"
	"io/fs"
	"strings"
	"time"

	goseidon "github.com/go-seidon/core"
	"github.com/go-seidon/core/internal/io"
)

type LocalConfig struct {
	StorageDir string
}

type LocalStorage struct {
	config      *LocalConfig
	FileManager io.FileManager
}

func (s *LocalStorage) UploadFile(p goseidon.UploadFileParam) (*goseidon.UploadFileResult, error) {
	rwPermission := fs.FileMode(0644)

	if !s.FileManager.IsExists(s.config.StorageDir) {
		err := s.FileManager.CreateDir(s.config.StorageDir, rwPermission)
		if err != nil {
			return nil, fmt.Errorf("failed create storage dir: %s", s.config.StorageDir)
		}
	}

	path := fmt.Sprintf("%s/%s", s.config.StorageDir, p.FileName)
	if s.FileManager.IsExists(path) {
		return nil, fmt.Errorf("file already exists")
	}

	err := s.FileManager.WriteFile(path, p.FileData, rwPermission)
	if err != nil {
		return nil, fmt.Errorf("failed storing file")
	}

	res := &goseidon.UploadFileResult{
		FileName: p.FileName,
	}
	return res, nil
}

func (s *LocalStorage) RetrieveFile(p goseidon.RetrieveFileParam) (*goseidon.RetrieveFileResult, error) {
	path := fmt.Sprintf("%s/%s", s.config.StorageDir, p.Id)
	if !s.FileManager.IsExists(path) {
		return nil, fmt.Errorf("file is not found")
	}

	file, err := s.FileManager.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed open file")
	}
	defer file.Close()

	binFile, err := s.FileManager.ReadFile(file)
	if err != nil {
		return nil, err
	}

	res := &goseidon.RetrieveFileResult{
		File: binFile,
	}
	return res, nil
}

func (s *LocalStorage) DeleteFile(p goseidon.DeleteFileParam) (*goseidon.DeleteFileResult, error) {
	path := fmt.Sprintf("%s/%s", s.config.StorageDir, p.Id)
	if !s.FileManager.IsExists(path) {
		return nil, fmt.Errorf("file is not found")
	}

	err := s.FileManager.RemoveFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed delete file")
	}

	res := &goseidon.DeleteFileResult{
		Id:        p.Id,
		DeletedAt: time.Now(),
	}
	return res, nil
}

func NewLocalStorage(c *LocalConfig) (*LocalStorage, error) {
	if c == nil {
		return nil, fmt.Errorf("invalid storage config")
	}
	fm, _ := io.NewFileManager()
	s := &LocalStorage{
		config:      c,
		FileManager: fm,
	}
	return s, nil
}

func NewLocalConfig(sDir string) (*LocalConfig, error) {
	if sDir == "" {
		return nil, fmt.Errorf("invalid storage directory")
	}
	sDir = strings.ToLower(sDir)
	sDir = strings.TrimSuffix(sDir, "/")
	c := &LocalConfig{
		StorageDir: sDir,
	}
	return c, nil
}
