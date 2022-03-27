package local

import (
	"fmt"
	"io/fs"
	"strings"

	goseidon "github.com/go-seidon/core"
	"github.com/go-seidon/core/internal/io"
)

type LocalConfig struct {
	StorageDir string
}

type localStorage struct {
	config      *LocalConfig
	fileManager io.FileManagerService
}

func (s *localStorage) UploadFile(p goseidon.UploadFileParam) (*goseidon.UploadFileResult, error) {
	path := fmt.Sprintf("%s/%s", s.config.StorageDir, p.FileName)
	if s.fileManager.IsFileExists(path) {
		return nil, fmt.Errorf("file already exists")
	}

	permission := fs.FileMode(0644)
	err := s.fileManager.WriteFile(path, p.FileData, permission)
	if err != nil {
		return nil, fmt.Errorf("failed storing file")
	}

	res := &goseidon.UploadFileResult{
		FileName: p.FileName,
	}
	return res, nil
}

func (s *localStorage) RetrieveFile(p goseidon.RetrieveFileParam) (*goseidon.RetrieveFileResult, error) {
	path := fmt.Sprintf("%s/%s", s.config.StorageDir, p.Id)
	if !s.fileManager.IsFileExists(path) {
		return nil, fmt.Errorf("file is not found")
	}

	file, err := s.fileManager.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed open file")
	}
	defer file.Close()

	binFile, err := s.fileManager.ReadFile(file)
	if err != nil {
		return nil, err
	}

	res := &goseidon.RetrieveFileResult{
		File: binFile,
	}
	return res, nil
}

func NewLocalStorage(c *LocalConfig, fm io.FileManagerService) (goseidon.Storage, error) {
	if c == nil {
		return nil, fmt.Errorf("invalid storage config")
	}
	if fm == nil {
		return nil, fmt.Errorf("invalid file manager")
	}
	s := &localStorage{
		config:      c,
		fileManager: fm,
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
