package goseidon

import (
	"context"
	"time"
)

type BinaryFile = []byte

type UploadFileParam struct {
	FileData BinaryFile
	FileId   string
	FileName string
	FileSize int64
}

type UploadFileResult struct {
	FileId     string
	FileName   string
	UploadedAt time.Time
}

type Uploader interface {
	UploadFile(ctx context.Context, p UploadFileParam) (*UploadFileResult, error)
}

type RetrieveFileParam struct {
	Id string
}

type RetrieveFileResult struct {
	File        BinaryFile
	RetrievedAt time.Time
}

type Retriever interface {
	RetrieveFile(ctx context.Context, p RetrieveFileParam) (*RetrieveFileResult, error)
}

type DeleteFileParam struct {
	Id string
}

type DeleteFileResult struct {
	Id        string
	DeletedAt time.Time
}

type Deleter interface {
	DeleteFile(ctx context.Context, p DeleteFileParam) (*DeleteFileResult, error)
}

type Storage interface {
	Uploader
	Retriever
	Deleter
}
