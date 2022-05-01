package g_storage_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	goseidon "github.com/go-seidon/core"
	"github.com/go-seidon/core/internal/clock"
	g_cloud "github.com/go-seidon/core/internal/g-cloud"
	g_storage "github.com/go-seidon/core/pkg/g-storage"
	gomock "github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGoogle(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Google Package")
}

var _ = Describe("Storage", func() {
	var (
		t GinkgoTInterface
	)

	BeforeEach(func() {
		t = GinkgoT()
	})

	Context("NewGoogleStorage function", func() {
		When("option is invalid", func() {
			It("should return error", func() {
				s, err := g_storage.NewGoogleStorage(nil)

				Expect(s).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("invalid google option")))
			})
		})

		When("failed apply option", func() {
			It("should return error", func() {

				s, err := g_storage.NewGoogleStorage(&withFailedApply{})

				Expect(s).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("failed apply option")))
			})
		})

		When("success apply option", func() {
			It("should return storage", func() {

				s, err := g_storage.NewGoogleStorage(&withSuccessApply{})

				Expect(s).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})
	})

	Context("UploadFile method", func() {
		var (
			ctx         context.Context
			s           *g_storage.GoogleStorage
			cl          *g_cloud.MockGoogleStorageClient
			wc          *g_cloud.MockWriteCloser
			cfg         *g_storage.GoogleConfig
			p           goseidon.UploadFileParam
			clo         *clock.MockClock
			currentTime time.Time
		)

		BeforeEach(func() {
			ctx = context.Background()
			cfg = &g_storage.GoogleConfig{
				BucketName:   "bucket-name",
				GoogleClient: &storage.Client{},
			}
			currentTime = time.Now()
			ctrl := gomock.NewController(t)
			cl = g_cloud.NewMockGoogleStorageClient(ctrl)
			wc = g_cloud.NewMockWriteCloser(ctrl)
			clo = clock.NewMockClock(ctrl)
			s = &g_storage.GoogleStorage{
				Client: cl,
				Clock:  clo,
				Config: cfg,
			}
			p = goseidon.UploadFileParam{
				FileData: make([]byte, 1),
				FileName: "file-name.jpg",
				FileSize: 1,
			}
		})

		When("context is invalid", func() {
			It("should return error", func() {
				res, err := s.UploadFile(nil, p)

				Expect(res).To(BeNil())
				Expect(err.Error()).To(Equal("invalid context"))
			})
		})

		When("failed copy file", func() {
			It("should return error", func() {
				cl.EXPECT().
					NewWriter(gomock.Eq(ctx), gomock.Eq(cfg.BucketName), gomock.Eq(p.FileName)).
					Return(wc).
					Times(1)
				buf := bytes.NewBuffer(p.FileData)
				cl.EXPECT().
					Copy(gomock.Eq(wc), gomock.Eq(buf)).
					Return(int64(0), fmt.Errorf("failed copy file")).
					Times(1)

				res, err := s.UploadFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err.Error()).To(Equal("failed copy file"))
			})
		})

		When("failed close file", func() {
			It("should return error", func() {
				wc.EXPECT().
					Close().
					Return(fmt.Errorf("failed close file")).
					Times(1)
				cl.EXPECT().
					NewWriter(gomock.Eq(ctx), gomock.Eq(cfg.BucketName), gomock.Eq(p.FileName)).
					Return(wc).
					Times(1)
				buf := bytes.NewBuffer(p.FileData)
				cl.EXPECT().
					Copy(gomock.Eq(wc), gomock.Eq(buf)).
					Return(int64(0), nil).
					Times(1)

				res, err := s.UploadFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err.Error()).To(Equal("failed close file"))
			})
		})

		When("success upload file", func() {
			It("should return result", func() {
				wc.EXPECT().
					Close().
					Return(nil).
					Times(1)
				cl.EXPECT().
					NewWriter(gomock.Eq(ctx), gomock.Eq(cfg.BucketName), gomock.Eq(p.FileName)).
					Return(wc).
					Times(1)
				buf := bytes.NewBuffer(p.FileData)
				cl.EXPECT().
					Copy(gomock.Eq(wc), gomock.Eq(buf)).
					Return(int64(0), nil).
					Times(1)
				clo.EXPECT().Now().Return(currentTime).Times(1)

				res, err := s.UploadFile(ctx, p)

				eRes := &goseidon.UploadFileResult{
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
			s           goseidon.Storage
			cfg         *g_storage.GoogleConfig
			cl          *g_cloud.MockGoogleStorageClient
			rc          *g_cloud.MockReadCloser
			p           goseidon.RetrieveFileParam
			clo         *clock.MockClock
			currentTime time.Time
		)

		BeforeEach(func() {
			ctx = context.Background()
			cfg = &g_storage.GoogleConfig{
				BucketName:   "bucket-name",
				GoogleClient: &storage.Client{},
			}
			ctrl := gomock.NewController(t)
			cl = g_cloud.NewMockGoogleStorageClient(ctrl)
			rc = g_cloud.NewMockReadCloser(ctrl)
			clo = clock.NewMockClock(ctrl)
			currentTime = time.Now()
			s = &g_storage.GoogleStorage{
				Client: cl,
				Config: cfg,
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

		When("failed create reader", func() {
			It("should return error", func() {
				cl.EXPECT().
					NewReader(gomock.Eq(ctx), gomock.Eq(cfg.BucketName), gomock.Eq(p.Id)).
					Return(nil, fmt.Errorf("failed create reader")).
					Times(1)

				res, err := s.RetrieveFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err.Error()).To(Equal("failed create reader"))
			})
		})

		When("failed read data", func() {
			It("should return error", func() {
				rc.EXPECT().
					Close().
					Times(1)
				rc.EXPECT().
					Read(gomock.Any()).
					Return(0, fmt.Errorf("failed read data")).
					Times(1)
				cl.EXPECT().
					NewReader(gomock.Eq(ctx), gomock.Eq(cfg.BucketName), gomock.Eq(p.Id)).
					Return(rc, nil).
					Times(1)

				res, err := s.RetrieveFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err.Error()).To(Equal("failed read data"))
			})
		})

		When("success retrieve data", func() {
			It("should return result", func() {
				rc.EXPECT().
					Close().
					Times(1)
				rc.EXPECT().
					Read(gomock.Any()).
					Return(1, io.EOF).
					Times(1)
				cl.EXPECT().
					NewReader(gomock.Eq(ctx), gomock.Eq(cfg.BucketName), gomock.Eq(p.Id)).
					Return(rc, nil).
					Times(1)
				clo.EXPECT().Now().Return(currentTime)

				res, err := s.RetrieveFile(ctx, p)

				eRes := &goseidon.RetrieveFileResult{
					File:        make([]byte, 1),
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
			s           goseidon.Storage
			cfg         *g_storage.GoogleConfig
			cl          *g_cloud.MockGoogleStorageClient
			clo         *clock.MockClock
			p           goseidon.DeleteFileParam
			currentTime time.Time
		)

		BeforeEach(func() {
			ctx = context.Background()
			cfg = &g_storage.GoogleConfig{
				BucketName:   "bucket-name",
				GoogleClient: &storage.Client{},
			}
			ctrl := gomock.NewController(t)
			cl = g_cloud.NewMockGoogleStorageClient(ctrl)
			currentTime = time.Now()
			clo = clock.NewMockClock(ctrl)
			s = &g_storage.GoogleStorage{
				Config: cfg,
				Client: cl,
				Clock:  clo,
			}
			p = goseidon.DeleteFileParam{
				Id: "mock-file-id",
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
				cl.EXPECT().
					Delete(gomock.Eq(ctx), gomock.Eq(cfg.BucketName), gomock.Eq(p.Id)).
					Return(fmt.Errorf("failed delete file")).
					Times(1)

				res, err := s.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err.Error()).To(Equal("failed delete file"))
			})
		})

		When("success delete file", func() {
			It("should return result", func() {
				cl.EXPECT().
					Delete(gomock.Eq(ctx), gomock.Eq(cfg.BucketName), gomock.Eq(p.Id)).
					Return(nil).
					Times(1)
				clo.EXPECT().Now().Return(currentTime).Times(1)

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

type withFailedApply struct {
}

func (o *withFailedApply) Apply(c *g_storage.GoogleConfig) error {
	return fmt.Errorf("failed apply option")
}

type withSuccessApply struct {
}

func (o *withSuccessApply) Apply(c *g_storage.GoogleConfig) error {
	return nil
}
