package g_storage

import (
	"context"

	gstorage "cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type GoogleConfig struct {
	BucketName   string
	GoogleClient *gstorage.Client
}

type GoogleStorageOption interface {
	Apply(c *GoogleConfig) error
}

type withGoogleClient struct {
	bucket string
	cl     *gstorage.Client
}

func (o *withGoogleClient) Apply(c *GoogleConfig) error {
	c.BucketName = o.bucket
	c.GoogleClient = o.cl
	return nil
}

func WithGoogleClient(bucketName string, googleClient *gstorage.Client) GoogleStorageOption {
	return &withGoogleClient{
		bucket: bucketName,
		cl:     googleClient,
	}
}

type withCredentialPath struct {
	bucket string
	path   string
}

func (o *withCredentialPath) Apply(c *GoogleConfig) error {
	ctx := context.Background()
	opt := option.WithCredentialsFile(o.path)
	gsClient, err := gstorage.NewClient(ctx, opt)
	if err != nil {
		return err
	}

	c.GoogleClient = gsClient
	return nil
}

func WithCredentialPath(bucketName string, path string) GoogleStorageOption {
	return &withCredentialPath{
		bucket: bucketName,
		path:   path,
	}
}
