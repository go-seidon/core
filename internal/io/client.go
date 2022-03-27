package io

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
)

type FileManagerService interface {
	IsFileExists(path string) bool
	WriteFile(name string, data []byte, perm fs.FileMode) error
	Open(path string) (*os.File, error)
	ReadFile(file *os.File) ([]byte, error)
}

type fileManager struct {
}

func (s *fileManager) IsFileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func (s *fileManager) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (s *fileManager) Open(path string) (*os.File, error) {
	return os.Open(path)
}

func (s *fileManager) ReadFile(file *os.File) ([]byte, error) {
	if file == nil {
		return nil, fmt.Errorf("invalid file")
	}

	reader := bufio.NewReader(file)
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func NewFileManager() (FileManagerService, error) {
	s := &fileManager{}
	return s, nil
}
