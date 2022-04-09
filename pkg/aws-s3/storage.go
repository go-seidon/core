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

type AwsS3Option struct {
	Credential *AwsS3Credential
	BucketName string
}

type AwsS3Client interface {
	PutObject(*s3.PutObjectInput) (*s3.PutObjectOutput, error)
	GetObject(*s3.GetObjectInput) (*s3.GetObjectOutput, error)
}

type awsS3Storage struct {
	client AwsS3Client
	option *AwsS3Option
}

func (s *awsS3Storage) UploadFile(p goseidon.UploadFileParam) (*goseidon.UploadFileResult, error) {
	_, err := s.client.PutObject(&s3.PutObjectInput{
		Body:   bytes.NewReader(p.FileData),
		Bucket: aws.String(s.option.BucketName),
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

func (s *awsS3Storage) RetrieveFile(p goseidon.RetrieveFileParam) (*goseidon.RetrieveFileResult, error) {
	out, err := s.client.GetObject(&s3.GetObjectInput{
		Key:    aws.String(p.Id),
		Bucket: aws.String(s.option.BucketName),
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

func NewAwsS3Option(bucketName string) (*AwsS3Option, error) {
	if bucketName == "" {
		return nil, fmt.Errorf("invalid aws s3 bucket name")
	}
	op := &AwsS3Option{
		BucketName: bucketName,
	}
	return op, nil
}

func NewAwsS3Credential(region, accessKeyId, secretAccessKey string) (*AwsS3Credential, error) {
	if region == "" {
		return nil, fmt.Errorf("invalid aws s3 region")
	}
	if accessKeyId == "" {
		return nil, fmt.Errorf("invalid aws s3 access key id")
	}
	if secretAccessKey == "" {
		return nil, fmt.Errorf("invalid aws s3 secret access key")
	}
	cr := &AwsS3Credential{
		Region:          region,
		AccessKeyId:     accessKeyId,
		SecretAccessKey: secretAccessKey,
	}
	return cr, nil
}

func NewAwsS3Client(cr *AwsS3Credential) (AwsS3Client, error) {
	if cr == nil {
		return nil, fmt.Errorf("invalid aws s3 credential")
	}
	config := &aws.Config{
		Region: aws.String(cr.Region),
		Credentials: credentials.NewStaticCredentials(
			cr.AccessKeyId, cr.SecretAccessKey, "",
		),
	}
	session, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}
	client := s3.New(session)
	return client, nil
}

func NewAwsS3Storage(cl AwsS3Client, op *AwsS3Option) (goseidon.Storage, error) {
	if cl == nil {
		return nil, fmt.Errorf("invalid aws s3 client")
	}
	if op == nil {
		return nil, fmt.Errorf("invalid aws s3 option")
	}

	storage := &awsS3Storage{
		client: cl,
		option: op,
	}
	return storage, nil
}
