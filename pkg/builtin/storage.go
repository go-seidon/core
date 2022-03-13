package builtin

import (
	"fmt"

	goseidon "github.com/go-seidon/core"
)

type builtinStorage struct {
}

func (s *builtinStorage) UploadFile(p goseidon.UploadFileParam) (*goseidon.UploadFileResult, error) {
	return nil, fmt.Errorf("failed upload file")
}

func (s *builtinStorage) RetrieveFile(p goseidon.RetrieveFileParam) (*goseidon.RetrieveFileResult, error) {
	return nil, fmt.Errorf("failed retrieve file")
}

func NewStorage() (*builtinStorage, error) {
	return &builtinStorage{}, nil
}
