package goseidon

import "time"

type BinaryFile = []byte

type UploadFileParam struct {
	FileData BinaryFile
	FileName string
	FileSize int64
}

type UploadFileResult struct {
	FileName string
}

type Uploader interface {
	UploadFile(p UploadFileParam) (*UploadFileResult, error)
}

type RetrieveFileParam struct {
	Id string
}

type RetrieveFileResult struct {
	File BinaryFile
}

type Retriever interface {
	RetrieveFile(p RetrieveFileParam) (*RetrieveFileResult, error)
}

type DeleteFileParam struct {
	Id string
}

type DeleteFileResult struct {
	Id        string
	DeletedAt time.Time
}

type Deleter interface {
	DeleteFile(p DeleteFileParam) (*DeleteFileResult, error)
}

type Storage interface {
	Uploader
	Retriever
	Deleter
}
