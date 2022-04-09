package aws_s3

import (
	"bytes"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	goseidon "github.com/go-seidon/core"
)

type AwsS3Credential struct {
	Region, AccessKeyId, SecretAccessKey string
}

type AwsS3Config struct {
	Credential *AwsS3Credential
	BucketName string
}

type AwsS3Client interface {
	PutObject(*s3.PutObjectInput) (*s3.PutObjectOutput, error)
	GetObject(*s3.GetObjectInput) (*s3.GetObjectOutput, error)
}

type AwsS3Storage struct {
	Client AwsS3Client
	Config *AwsS3Config
}

func (s *AwsS3Storage) UploadFile(p goseidon.UploadFileParam) (*goseidon.UploadFileResult, error) {
	_, err := s.Client.PutObject(&s3.PutObjectInput{
		Body:   bytes.NewReader(p.FileData),
		Bucket: aws.String(s.Config.BucketName),
		Key:    aws.String(p.FileName),
	})
	if err != nil {
		return nil, err
	}
	res := &goseidon.UploadFileResult{
		FileName: p.FileName,
	}
	return res, nil
}

func (s *AwsS3Storage) RetrieveFile(p goseidon.RetrieveFileParam) (*goseidon.RetrieveFileResult, error) {
	out, err := s.Client.GetObject(&s3.GetObjectInput{
		Key:    aws.String(p.Id),
		Bucket: aws.String(s.Config.BucketName),
	})
	if err != nil {
		return nil, err
	}

	fileData, err := io.ReadAll(out.Body)
	if err != nil {
		return nil, err
	}

	res := &goseidon.RetrieveFileResult{
		File: fileData,
	}
	return res, nil
}

func NewAwsS3Client(cr AwsS3Credential) (AwsS3Client, error) {
	awsCfg := &aws.Config{
		Region: aws.String(cr.Region),
		Credentials: credentials.NewStaticCredentials(
			cr.AccessKeyId, cr.SecretAccessKey, "",
		),
	}
	session, err := session.NewSession(awsCfg)
	if err != nil {
		return nil, err
	}
	client := s3.New(session)
	return client, nil
}

func NewAwsS3Storage(c *AwsS3Config) (*AwsS3Storage, error) {
	if c == nil {
		return nil, fmt.Errorf("invalid aws s3 config")
	}
	if c.Credential == nil {
		return nil, fmt.Errorf("invalid aws s3 credential")
	}
	if c.Credential.Region == "" {
		return nil, fmt.Errorf("invalid aws s3 region")
	}
	if c.Credential.AccessKeyId == "" {
		return nil, fmt.Errorf("invalid aws s3 access key")
	}
	if c.Credential.SecretAccessKey == "" {
		return nil, fmt.Errorf("invalid aws s3 secret access key")
	}
	if c.BucketName == "" {
		return nil, fmt.Errorf("invalid aws s3 bucket name")
	}

	cl, _ := NewAwsS3Client(*c.Credential)

	storage := &AwsS3Storage{
		Client: cl,
		Config: c,
	}
	return storage, nil
}
