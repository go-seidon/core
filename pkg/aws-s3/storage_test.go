package aws_s3_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	goseidon "github.com/go-seidon/core"
	awsmock "github.com/go-seidon/core/internal/aws"
	"github.com/go-seidon/core/internal/clock"
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

	Context("NewAwsS3Storage function", func() {
		When("option is invalid", func() {
			It("should return error", func() {
				s, err := aws_s3.NewAwsS3Storage(nil)

				Expect(err).To(Equal(fmt.Errorf("invalid aws s3 option")))
				Expect(s).To(BeNil())
			})
		})

		When("failed apply option", func() {
			It("should return error", func() {
				s, err := aws_s3.NewAwsS3Storage(&withFailedOption{})

				Expect(err).To(Equal(fmt.Errorf("failed apply option")))
				Expect(s).To(BeNil())
			})
		})

		When("success create storage", func() {
			It("should return aws_s3 storage", func() {
				s, err := aws_s3.NewAwsS3Storage(&withSuccessOption{})

				Expect(s).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})
	})

	Context("UploadFile method", func() {
		var (
			ctx         context.Context
			s           *aws_s3.AwsS3Storage
			p           goseidon.UploadFileParam
			cfg         *aws_s3.AwsS3Config
			cl          *awsmock.MockAwsS3Client
			clo         *clock.MockClock
			currentTime time.Time
		)

		BeforeEach(func() {
			ctx = context.Background()
			ctrl := gomock.NewController(t)
			cfg = &aws_s3.AwsS3Config{
				Region:          "mock-region",
				AccessKeyId:     "mock-access-key-id",
				SecretAccessKey: "mock-secret-access-key",
				BucketName:      "mock-bucket-name",
			}
			cl = awsmock.NewMockAwsS3Client(ctrl)
			clo = clock.NewMockClock(ctrl)
			currentTime = time.Now()

			s = &aws_s3.AwsS3Storage{
				Config: cfg,
				Client: cl,
				Clock:  clo,
			}
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
					Key:    aws.String(p.FileId),
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
					Key:    aws.String(p.FileId),
				}
				out := &s3.PutObjectOutput{}
				cl.EXPECT().
					PutObject(gomock.Eq(param)).
					Return(out, nil).
					Times(1)
				clo.EXPECT().Now().Return(currentTime)

				res, err := s.UploadFile(ctx, p)

				eRes := &goseidon.UploadFileResult{
					FileId:     p.FileId,
					FileName:   p.FileName,
					UploadedAt: currentTime,
				}
				Expect(res).To(Equal(eRes))
				Expect(err).To(BeNil())
			})
		})
	})

	Context("RetrieveFile method", func() {
		var (
			ctx         context.Context
			s           *aws_s3.AwsS3Storage
			p           goseidon.RetrieveFileParam
			cl          *awsmock.MockAwsS3Client
			cfg         *aws_s3.AwsS3Config
			clo         *clock.MockClock
			currentTime time.Time
		)

		BeforeEach(func() {
			ctx = context.Background()
			ctrl := gomock.NewController(t)
			cl = awsmock.NewMockAwsS3Client(ctrl)
			cfg = &aws_s3.AwsS3Config{
				Region:          "mock-region",
				AccessKeyId:     "mock-access-key-id",
				SecretAccessKey: "mock-secret-access-key",
				BucketName:      "mock-bucket-name",
			}
			clo = clock.NewMockClock(ctrl)
			currentTime = time.Now()

			s = &aws_s3.AwsS3Storage{
				Config: cfg,
				Client: cl,
				Clock:  clo,
			}
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
				clo.EXPECT().Now().Return(currentTime)

				res, err := s.RetrieveFile(ctx, p)

				eRes := &goseidon.RetrieveFileResult{
					File:        []byte{},
					RetrievedAt: currentTime,
				}
				Expect(res).To(Equal(eRes))
				Expect(err).To(BeNil())
			})
		})
	})

	Context("DeleteFile method", func() {
		var (
			ctx         context.Context
			s           *aws_s3.AwsS3Storage
			p           goseidon.DeleteFileParam
			cl          *awsmock.MockAwsS3Client
			cfg         *aws_s3.AwsS3Config
			clo         *clock.MockClock
			currentTime time.Time
		)

		BeforeEach(func() {
			ctx = context.Background()
			p = goseidon.DeleteFileParam{
				Id: "mock-file-id",
			}
			ctrl := gomock.NewController(t)
			cl = awsmock.NewMockAwsS3Client(ctrl)
			clo = clock.NewMockClock(ctrl)
			currentTime = time.Now()
			cfg = &aws_s3.AwsS3Config{
				Region:          "mock-region",
				AccessKeyId:     "mock-access-key-id",
				SecretAccessKey: "mock-secret-access-key",
				BucketName:      "mock-bucket-name",
			}
			s = &aws_s3.AwsS3Storage{
				Config: cfg,
				Client: cl,
				Clock:  clo,
			}
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
				param := &s3.DeleteObjectInput{
					Key:    aws.String(p.Id),
					Bucket: aws.String(cfg.BucketName),
				}

				cl.EXPECT().
					DeleteObject(gomock.Eq(param)).
					Return(nil, nil).
					Times(1)

				clo.EXPECT().Now().Return(currentTime)

				res, err := s.DeleteFile(ctx, p)

				eRes := &goseidon.DeleteFileResult{
					Id:        p.Id,
					DeletedAt: currentTime,
				}
				Expect(res).To(Equal(eRes))
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

type withFailedOption struct {
}

func (o *withFailedOption) Apply(c *aws_s3.AwsS3Config) error {
	return fmt.Errorf("failed apply option")
}

type withSuccessOption struct {
}

func (o *withSuccessOption) Apply(c *aws_s3.AwsS3Config) error {
	return nil
}
