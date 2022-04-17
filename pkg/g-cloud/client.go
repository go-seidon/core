package g_cloud

import (
	"context"
	"io"

	gstorage "cloud.google.com/go/storage"
)

type Writer interface {
	Write(p []byte) (n int, err error)
}

type Reader interface {
	Read(p []byte) (n int, err error)
}

type Closer interface {
	Close() error
}

type WriteCloser interface {
	Writer
	Closer
}

type ReadCloser interface {
	Reader
	Closer
}

type GoogleStorageClient interface {
	NewWriter(ctx context.Context, bucketName, fileId string) WriteCloser
	NewReader(ctx context.Context, bucketName, fileId string) (ReadCloser, error)
	Delete(ctx context.Context, bucketName, fileId string) error
	Copy(dst Writer, src Reader) (written int64, err error)
}

type googleStorageClient struct {
	client *gstorage.Client
}

func (c *googleStorageClient) NewWriter(ctx context.Context, bucketName, fileId string) WriteCloser {
	return c.client.
		Bucket(bucketName).
		Object(fileId).
		NewWriter(ctx)
}

func (c *googleStorageClient) NewReader(ctx context.Context, bucketName, fileId string) (ReadCloser, error) {
	return c.client.
		Bucket(bucketName).
		Object(fileId).
		NewReader(ctx)
}

func (c *googleStorageClient) Copy(dst Writer, src Reader) (written int64, err error) {
	return io.Copy(dst, src)
}

func (c *googleStorageClient) Delete(ctx context.Context, bucketName, fileId string) error {
	return c.client.Bucket(bucketName).Object(fileId).Delete(ctx)
}
