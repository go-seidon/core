package aws_s3_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

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

	Context("NewAwsS3Storage function", func() {
		When("success create storage", func() {
			It("should return aws_s3 storage", func() {
				ctrl := gomock.NewController(t)
				cl := aws_s3.NewMockAwsS3Client(ctrl)
				op := &aws_s3.AwsS3Option{}
				s, err := aws_s3.NewAwsS3Storage(cl, op)

				Expect(s).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})

		When("client is invalid", func() {
			It("should return error", func() {
				op := &aws_s3.AwsS3Option{}
				s, err := aws_s3.NewAwsS3Storage(nil, op)

				Expect(err).To(Equal(fmt.Errorf("invalid aws s3 client")))
				Expect(s).To(BeNil())
			})
		})

		When("option is invalid", func() {
			It("should return error", func() {
				ctrl := gomock.NewController(t)
				cl := aws_s3.NewMockAwsS3Client(ctrl)
				s, err := aws_s3.NewAwsS3Storage(cl, nil)

				Expect(err).To(Equal(fmt.Errorf("invalid aws s3 option")))
				Expect(s).To(BeNil())
			})
		})
	})

	Context("NewAwsS3Client function", func() {
		When("success create client", func() {
			It("should return aws s3 client", func() {
				cr := &aws_s3.AwsS3Credential{}

				cl, err := aws_s3.NewAwsS3Client(cr)

				Expect(err).To(BeNil())
				Expect(cl).ToNot(BeNil())
			})
		})

		When("credential is invalid", func() {
			It("should return error", func() {
				cl, err := aws_s3.NewAwsS3Client(nil)

				Expect(err).To(Equal(fmt.Errorf("invalid aws s3 credential")))
				Expect(cl).To(BeNil())
			})
		})

		When("failed create session", func() {
			It("should return error", func() {
				os.Setenv("AWS_SDK_LOAD_CONFIG", "true")
				os.Setenv("AWS_S3_USE_ARN_REGION", "invalid_value")
				cr := &aws_s3.AwsS3Credential{}
				cl, err := aws_s3.NewAwsS3Client(cr)

				Expect(err).To(Equal(fmt.Errorf("failed to load environment config, invalid value for environment variable, AWS_S3_USE_ARN_REGION=invalid_value, need true or false")))
				Expect(cl).To(BeNil())
			})
		})
	})

	Context("NewAwsS3Credential function", func() {
		When("region is invalid", func() {
			It("should return error", func() {
				cr, err := aws_s3.NewAwsS3Credential("", "some-access-key", "some-secret-key")

				Expect(err).To(Equal(fmt.Errorf("invalid aws s3 region")))
				Expect(cr).To(BeNil())
			})
		})

		When("accessKeyId is invalid", func() {
			It("should return error", func() {
				cr, err := aws_s3.NewAwsS3Credential("ap-southeast-2", "", "some-secret-key")

				Expect(err).To(Equal(fmt.Errorf("invalid aws s3 access key id")))
				Expect(cr).To(BeNil())
			})
		})

		When("secretAccessKey is invalid", func() {
			It("should return error", func() {
				cr, err := aws_s3.NewAwsS3Credential("ap-southeast-2", "some-access-key", "")

				Expect(err).To(Equal(fmt.Errorf("invalid aws s3 secret access key")))
				Expect(cr).To(BeNil())
			})
		})

		When("all param is valid", func() {
			It("should return aws s3 credential", func() {
				cr, err := aws_s3.NewAwsS3Credential("ap-southeast-2", "some-access-key", "some-secret-key")

				Expect(err).To(BeNil())
				Expect(cr).ToNot(BeNil())
			})
		})
	})

	Context("NewAwsS3Option function", func() {
		When("bucketName is invalid", func() {
			It("should return error", func() {
				op, err := aws_s3.NewAwsS3Option("")

				Expect(err).To(Equal(fmt.Errorf("invalid aws s3 bucket name")))
				Expect(op).To(BeNil())
			})
		})

		When("all param is valid", func() {
			It("should return aws s3 option", func() {
				op, err := aws_s3.NewAwsS3Option("some-bucket-name")

				Expect(err).To(BeNil())
				Expect(op).ToNot(BeNil())
			})
		})
	})

	Context("UploadFile method", func() {
		var (
			s  goseidon.Storage
			p  goseidon.UploadFileParam
			cl *aws_s3.MockAwsS3Client
			op *aws_s3.AwsS3Option
		)

		BeforeEach(func() {
			ctrl := gomock.NewController(t)
			cl = aws_s3.NewMockAwsS3Client(ctrl)
			op = &aws_s3.AwsS3Option{
				BucketName: "mock-bucket-name",
			}
			storage, _ := aws_s3.NewAwsS3Storage(cl, op)
			s = storage
			p = goseidon.UploadFileParam{}
		})

		When("failed upload file", func() {
			It("should return error", func() {
				param := &s3.PutObjectInput{
					Body:   bytes.NewReader(p.FileData),
					Bucket: aws.String(op.BucketName),
					Key:    aws.String(p.FileName),
				}
				cl.EXPECT().
					PutObject(gomock.Eq(param)).
					Return(nil, fmt.Errorf("failed upload file")).
					Times(1)
				res, err := s.UploadFile(p)

				Expect(res).To(BeNil())
				Expect(err.Error()).To(Equal("failed upload file"))
			})
		})

		When("success upload file", func() {
			It("should return result", func() {
				param := &s3.PutObjectInput{
					Body:   bytes.NewReader(p.FileData),
					Bucket: aws.String(op.BucketName),
					Key:    aws.String(p.FileName),
				}
				out := &s3.PutObjectOutput{}
				cl.EXPECT().
					PutObject(gomock.Eq(param)).
					Return(out, nil).
					Times(1)
				res, err := s.UploadFile(p)

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
			s  goseidon.Storage
			p  goseidon.RetrieveFileParam
			cl *aws_s3.MockAwsS3Client
			op *aws_s3.AwsS3Option
		)

		BeforeEach(func() {
			ctrl := gomock.NewController(t)
			cl = aws_s3.NewMockAwsS3Client(ctrl)
			op = &aws_s3.AwsS3Option{
				BucketName: "mock-bucket-name",
			}
			storage, _ := aws_s3.NewAwsS3Storage(cl, op)
			s = storage
			p = goseidon.RetrieveFileParam{
				Id: "mock-file-id",
			}
		})

		When("failed retrieve file", func() {
			It("should return error", func() {
				param := &s3.GetObjectInput{
					Key:    aws.String(p.Id),
					Bucket: aws.String(op.BucketName),
				}
				cl.EXPECT().
					GetObject(gomock.Eq(param)).
					Return(nil, fmt.Errorf("failed retrieve file")).
					Times(1)
				res, err := s.RetrieveFile(p)

				Expect(res).To(BeNil())
				Expect(err.Error()).To(Equal("failed retrieve file"))
			})
		})

		When("failed read file", func() {
			It("should return error", func() {
				param := &s3.GetObjectInput{
					Key:    aws.String(p.Id),
					Bucket: aws.String(op.BucketName),
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

				res, err := s.RetrieveFile(p)

				Expect(res).To(BeNil())
				Expect(err.Error()).To(Equal("failed read file"))
			})
		})

		When("success retrieve file", func() {
			It("should return result", func() {
				param := &s3.GetObjectInput{
					Key:    aws.String(p.Id),
					Bucket: aws.String(op.BucketName),
				}

				out := &s3.GetObjectOutput{
					Body: &readCloser{},
				}
				cl.EXPECT().
					GetObject(gomock.Eq(param)).
					Return(out, nil).
					Times(1)
				res, err := s.RetrieveFile(p)

				eRes := &goseidon.RetrieveFileResult{
					File: []byte{},
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
