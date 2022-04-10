package io

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
)

type FileManager interface {
	IsExists(path string) bool
	CreateDir(path string, perm fs.FileMode) error
	WriteFile(name string, data []byte, perm fs.FileMode) error
	Open(path string) (*os.File, error)
	ReadFile(file *os.File) ([]byte, error)
	RemoveFile(path string) error
}

type fileManager struct {
}

func (fm *fileManager) IsExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func (fm *fileManager) CreateDir(path string, perm fs.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (fm *fileManager) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (fm *fileManager) Open(path string) (*os.File, error) {
	return os.Open(path)
}

func (fm *fileManager) ReadFile(file *os.File) ([]byte, error) {
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

func (fm *fileManager) RemoveFile(path string) error {
	return os.Remove(path)
}

func NewFileManager() (FileManager, error) {
	s := &fileManager{}
	return s, nil
}
