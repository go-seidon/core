package aws_s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type AwsS3Config struct {
	Region, AccessKeyId, SecretAccessKey string
	BucketName                           string

	Client AwsS3Client
}

type AwsS3StorageOption interface {
	Apply(c *AwsS3Config) error
}

type withStaticCredential struct {
	region, accessKey, secretKey, bucketName string
}

func (o *withStaticCredential) Apply(c *AwsS3Config) error {
	c.BucketName = o.bucketName
	c.AccessKeyId = o.accessKey
	c.SecretAccessKey = o.secretKey
	c.Region = o.region

	cr := credentials.NewStaticCredentials(
		c.AccessKeyId, c.SecretAccessKey, "",
	)
	awsCfg := &aws.Config{
		Region:      aws.String(c.Region),
		Credentials: cr,
	}
	session, err := session.NewSession(awsCfg)
	if err != nil {
		return err
	}

	c.Client = s3.New(session)
	return nil
}

func WithStatisCredential(region, accessKey, secretKey, bucketName string) AwsS3StorageOption {
	return &withStaticCredential{
		region:     region,
		accessKey:  accessKey,
		secretKey:  secretKey,
		bucketName: bucketName,
	}
}
