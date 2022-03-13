package goseidon

type BinaryFile = []byte

type UploadFileParam struct {
	FileData BinaryFile
	FileName string
}

type UploadFileResult struct {
	FileName string
}

type Uploader interface {
	UploadFile(UploadFileParam) (*UploadFileResult, error)
}

type RetrieveFileParam struct {
	Id string
}

type RetrieveFileResult struct {
	File BinaryFile
}

type Retriever interface {
	RetrieveFile(RetrieveFileParam) (*RetrieveFileResult, error)
}

type Storage interface {
	Uploader
	Retriever
}
