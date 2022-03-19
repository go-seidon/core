package aws_s3_test

import (
	"bytes"
	"fmt"
	"testing"

	// goseidon "github.com/go-seidon/core"
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
					File: nil,
				}
				Expect(res).To(Equal(eRes))
				Expect(err).To(BeNil())
			})
		})
	})
})

type readCloser struct {
}

func (rc *readCloser) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (rc *readCloser) Close() error {
	return nil
}
