package aws_s3_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	goseidon "github.com/go-seidon/core"
	aws_s3 "github.com/go-seidon/core/pkg/aws-s3"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAwsS3(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AwsS3 Package")
}

var _ = Describe("Storage", func() {
	var (
		t GinkgoTInterface
	)

	BeforeEach(func() {
		t = GinkgoT()
	})

	Context("NewAwsS3Config function", func() {
		When("all param is valid", func() {
			It("should return aws_s3 config", func() {
				res, err := aws_s3.NewAwsS3Config(
					"mock-region",
					"mock-access-key-id",
					"mock-secret-access-key",
					"mock-bucket-name",
				)

				Expect(res).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})

		When("region is invalid", func() {
			It("should return error", func() {
				res, err := aws_s3.NewAwsS3Config(
					"",
					"mock-access-key-id",
					"mock-secret-access-key",
					"mock-bucket-name",
				)

				Expect(err).To(Equal(fmt.Errorf("invalid aws s3 region")))
				Expect(res).To(BeNil())
			})
		})

		When("access key is invalid", func() {
			It("should return error", func() {
				res, err := aws_s3.NewAwsS3Config(
					"mock-region",
					"",
					"mock-secret-access-key",
					"mock-bucket-name",
				)

				Expect(err).To(Equal(fmt.Errorf("invalid aws s3 access key")))
				Expect(res).To(BeNil())
			})
		})

		When("secret access key is invalid", func() {
			It("should return error", func() {
				res, err := aws_s3.NewAwsS3Config(
					"mock-region",
					"mock-access-key-id",
					"",
					"mock-bucket-name",
				)

				Expect(err).To(Equal(fmt.Errorf("invalid aws s3 secret access key")))
				Expect(res).To(BeNil())
			})
		})

		When("bucket name is invalid", func() {
			It("should return error", func() {
				res, err := aws_s3.NewAwsS3Config(
					"mock-region",
					"mock-access-key-id",
					"mock-secret-access-key",
					"",
				)

				Expect(err).To(Equal(fmt.Errorf("invalid aws s3 bucket name")))
				Expect(res).To(BeNil())
			})
		})

	})

	Context("NewAwsS3Client function", func() {
		When("success create client", func() {
			It("should return aws s3 client", func() {
				cr := aws_s3.AwsS3Credential{}

				cl, err := aws_s3.NewAwsS3Client(cr)

				Expect(err).To(BeNil())
				Expect(cl).ToNot(BeNil())
			})
		})

		When("failed create session", func() {
			It("should return error", func() {
				os.Setenv("AWS_SDK_LOAD_CONFIG", "true")
				os.Setenv("AWS_S3_USE_ARN_REGION", "invalid_value")
				cr := aws_s3.AwsS3Credential{}
				cl, err := aws_s3.NewAwsS3Client(cr)

				Expect(err).To(Equal(fmt.Errorf("failed to load environment config, invalid value for environment variable, AWS_S3_USE_ARN_REGION=invalid_value, need true or false")))
				Expect(cl).To(BeNil())
			})
		})
	})

	Context("NewAwsS3Storage function", func() {
		When("all param is valid", func() {
			It("should return aws_s3 storage", func() {
				ctrl := gomock.NewController(t)
				cl := aws_s3.NewMockAwsS3Client(ctrl)
				cfg, _ := aws_s3.NewAwsS3Config(
					"mock-region",
					"mock-access-key-id",
					"mock-secret-access-key",
					"mock-bucket-name",
				)
				s, err := aws_s3.NewAwsS3Storage(cfg)
				s.Client = cl

				Expect(s).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})

		When("config is invalid", func() {
			It("should return error", func() {
				s, err := aws_s3.NewAwsS3Storage(nil)

				Expect(err).To(Equal(fmt.Errorf("invalid aws s3 config")))
				Expect(s).To(BeNil())
			})
		})
	})

	Context("UploadFile method", func() {
		var (
			ctx context.Context
			s   goseidon.Storage
			p   goseidon.UploadFileParam
			cl  *aws_s3.MockAwsS3Client
			cfg *aws_s3.AwsS3Config
		)

		BeforeEach(func() {
			ctx = context.Background()
			ctrl := gomock.NewController(t)
			cl = aws_s3.NewMockAwsS3Client(ctrl)
			cfg, _ = aws_s3.NewAwsS3Config(
				"mock-region",
				"mock-access-key-id",
				"mock-secret-access-key",
				"mock-bucket-name",
			)
			storage, _ := aws_s3.NewAwsS3Storage(cfg)
			storage.Client = cl
			s = storage
			p = goseidon.UploadFileParam{}
		})

		When("context is invalid", func() {
			It("should return error", func() {
				res, err := s.UploadFile(nil, p)

				Expect(res).To(BeNil())
				Expect(err.Error()).To(Equal("invalid context"))
			})
		})

		When("failed upload file", func() {
			It("should return error", func() {
				param := &s3.PutObjectInput{
					Body:   bytes.NewReader(p.FileData),
					Bucket: aws.String(cfg.BucketName),
					Key:    aws.String(p.FileName),
				}
				cl.EXPECT().
					PutObject(gomock.Eq(param)).
					Return(nil, fmt.Errorf("failed upload file")).
					Times(1)
				res, err := s.UploadFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err.Error()).To(Equal("failed upload file"))
			})
		})

		When("success upload file", func() {
			It("should return result", func() {
				param := &s3.PutObjectInput{
					Body:   bytes.NewReader(p.FileData),
					Bucket: aws.String(cfg.BucketName),
					Key:    aws.String(p.FileName),
				}
				out := &s3.PutObjectOutput{}
				cl.EXPECT().
					PutObject(gomock.Eq(param)).
					Return(out, nil).
					Times(1)
				res, err := s.UploadFile(ctx, p)

				eRes := &goseidon.UploadFileResult{
					FileName: p.FileName,
				}
				Expect(res).To(Equal(eRes))
				Expect(err).To(BeNil())
			})
		})
	})

	Context("RetrieveFile method", func() {
		var (
			ctx context.Context
			s   goseidon.Storage
			p   goseidon.RetrieveFileParam
			cl  *aws_s3.MockAwsS3Client
			cfg *aws_s3.AwsS3Config
		)

		BeforeEach(func() {
			ctx = context.Background()
			ctrl := gomock.NewController(t)
			cl = aws_s3.NewMockAwsS3Client(ctrl)
			cfg, _ = aws_s3.NewAwsS3Config(
				"mock-region",
				"mock-access-key-id",
				"mock-secret-access-key",
				"mock-bucket-name",
			)
			storage, _ := aws_s3.NewAwsS3Storage(cfg)
			storage.Client = cl
			s = storage
			p = goseidon.RetrieveFileParam{
				Id: "mock-file-id",
			}
		})

		When("context is invalid", func() {
			It("should return error", func() {
				res, err := s.RetrieveFile(nil, p)

				Expect(res).To(BeNil())
				Expect(err.Error()).To(Equal("invalid context"))
			})
		})

		When("failed retrieve file", func() {
			It("should return error", func() {
				param := &s3.GetObjectInput{
					Key:    aws.String(p.Id),
					Bucket: aws.String(cfg.BucketName),
				}
				cl.EXPECT().
					GetObject(gomock.Eq(param)).
					Return(nil, fmt.Errorf("failed retrieve file")).
					Times(1)
				res, err := s.RetrieveFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err.Error()).To(Equal("failed retrieve file"))
			})
		})

		When("failed read file", func() {
			It("should return error", func() {
				param := &s3.GetObjectInput{
					Key:    aws.String(p.Id),
					Bucket: aws.String(cfg.BucketName),
				}

				out := &s3.GetObjectOutput{
					Body: &readCloser{
						readShouldError: true,
					},
				}
				cl.EXPECT().
					GetObject(gomock.Eq(param)).
					Return(out, nil).
					Times(1)

				res, err := s.RetrieveFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err.Error()).To(Equal("failed read file"))
			})
		})

		When("success retrieve file", func() {
			It("should return result", func() {
				param := &s3.GetObjectInput{
					Key:    aws.String(p.Id),
					Bucket: aws.String(cfg.BucketName),
				}

				out := &s3.GetObjectOutput{
					Body: &readCloser{},
				}
				cl.EXPECT().
					GetObject(gomock.Eq(param)).
					Return(out, nil).
					Times(1)
				res, err := s.RetrieveFile(ctx, p)

				eRes := &goseidon.RetrieveFileResult{
					File: []byte{},
				}
				Expect(res).To(Equal(eRes))
				Expect(err).To(BeNil())
			})
		})
	})

	Context("DeleteFile method", func() {
		var (
			ctx context.Context
			s   goseidon.Storage
			p   goseidon.DeleteFileParam
			cl  *aws_s3.MockAwsS3Client
			cfg *aws_s3.AwsS3Config
		)

		BeforeEach(func() {
			ctx = context.Background()
			p = goseidon.DeleteFileParam{
				Id: "mock-file-id",
			}
			ctrl := gomock.NewController(t)
			cl = aws_s3.NewMockAwsS3Client(ctrl)
			cfg, _ = aws_s3.NewAwsS3Config(
				"mock-region",
				"mock-access-key-id",
				"mock-secret-access-key",
				"mock-bucket-name",
			)
			storage, _ := aws_s3.NewAwsS3Storage(cfg)
			storage.Client = cl
			s = storage
		})

		When("context is invalid", func() {
			It("should return error", func() {
				res, err := s.DeleteFile(nil, p)

				Expect(res).To(BeNil())
				Expect(err.Error()).To(Equal("invalid context"))
			})
		})

		When("failed delete file", func() {
			It("should return error", func() {
				param := &s3.DeleteObjectInput{
					Key:    aws.String(p.Id),
					Bucket: aws.String(cfg.BucketName),
				}

				cl.EXPECT().
					DeleteObject(gomock.Eq(param)).
					Return(nil, fmt.Errorf("failed delete file")).
					Times(1)

				res, err := s.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("failed delete file")))
			})
		})

		When("success delete file", func() {
			It("should return error", func() {
				currentTime := time.Now()

				param := &s3.DeleteObjectInput{
					Key:    aws.String(p.Id),
					Bucket: aws.String(cfg.BucketName),
				}

				cl.EXPECT().
					DeleteObject(gomock.Eq(param)).
					Return(nil, nil).
					Times(1)

				res, err := s.DeleteFile(ctx, p)

				isAfterOrEqual := res.DeletedAt.After(currentTime) || res.DeletedAt.Equal(currentTime)
				Expect(res.Id).To(Equal(p.Id))
				Expect(isAfterOrEqual).To(BeTrue())
				Expect(err).To(BeNil())
			})
		})

	})
})

type readCloser struct {
	readShouldError bool
}

func (rc *readCloser) Read(p []byte) (n int, err error) {
	if rc.readShouldError {
		return 0, fmt.Errorf("failed read file")
	}
	return 0, io.EOF
}

func (rc *readCloser) Close() error {
	return nil
}
