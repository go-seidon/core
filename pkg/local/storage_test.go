package local_test

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"testing"
	"time"

	goseidon "github.com/go-seidon/core"
	"github.com/go-seidon/core/internal/clock"
	"github.com/go-seidon/core/internal/io"
	"github.com/go-seidon/core/pkg/local"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLocal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Local Package")
}

var _ = Describe("Storage", func() {
	Context("NewLocalStorage function", func() {
		var (
			c *local.LocalConfig
		)

		BeforeEach(func() {
			c = &local.LocalConfig{
				StorageDir: "storage",
			}
		})

		When("config is invalid", func() {
			It("should return error", func() {
				s, err := local.NewLocalStorage(nil)

				Expect(s).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("invalid storage config")))
			})
		})

		When("success create storage", func() {
			It("should return local storage", func() {
				s, err := local.NewLocalStorage(c)

				Expect(s).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})
	})

	Context("NewLocalConfig function", func() {
		When("storage directory is invalid", func() {
			It("should return error", func() {
				c, err := local.NewLocalConfig("")

				Expect(c).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("invalid storage directory")))
			})
		})

		When("all param is valid", func() {
			It("should return config", func() {
				c, err := local.NewLocalConfig("storage/custom-dir/")

				eConfig := &local.LocalConfig{
					StorageDir: "storage/custom-dir",
				}
				Expect(c).To(Equal(eConfig))
				Expect(err).To(BeNil())
			})
		})
	})

	Context("UploadFile method", func() {
		var (
			ctx         context.Context
			s           goseidon.Storage
			p           goseidon.UploadFileParam
			c           *local.LocalConfig
			fm          *io.MockFileManager
			clo         *clock.MockClock
			currentTime time.Time
		)

		BeforeEach(func() {
			ctx = context.Background()
			p = goseidon.UploadFileParam{
				FileName: "image.jpg",
				FileData: make([]byte, 1),
			}
			c = &local.LocalConfig{
				StorageDir: "storage",
			}
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			fm = io.NewMockFileManager(ctrl)
			currentTime = time.Now()
			clo = clock.NewMockClock(ctrl)
			storage, _ := local.NewLocalStorage(c)
			storage.FileManager = fm
			storage.Clock = clo
			s = storage
		})

		When("context is invalid", func() {
			It("should return error", func() {
				res, err := s.UploadFile(nil, p)

				Expect(res).To(BeNil())
				Expect(err.Error()).To(Equal("invalid context"))
			})
		})

		When("failed create storage dir", func() {
			It("should return error", func() {
				fm.EXPECT().
					IsExists(gomock.Eq(c.StorageDir)).
					Return(false).
					Times(1)

				fm.EXPECT().
					CreateDir(gomock.Eq(c.StorageDir), gomock.Eq(fs.FileMode(0644))).
					Return(fmt.Errorf("invalid storage dir")).
					Times(1)

				res, err := s.UploadFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("failed create storage dir: %s", c.StorageDir)))
			})
		})

		When("file already exists", func() {
			It("should return error", func() {
				fm.EXPECT().
					IsExists(gomock.Eq(c.StorageDir)).
					Return(true).
					Times(1)

				fm.EXPECT().
					IsExists(gomock.Eq(c.StorageDir + "/" + p.FileName)).
					Return(true).
					Times(1)

				res, err := s.UploadFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("file already exists")))
			})
		})

		When("failed upload file", func() {
			It("should return error", func() {
				fm.EXPECT().
					IsExists(gomock.Eq(c.StorageDir)).
					Return(true).
					Times(1)

				fm.EXPECT().
					IsExists(gomock.Eq(c.StorageDir + "/" + p.FileName)).
					Return(false).
					Times(1)

				fm.EXPECT().
					WriteFile(
						gomock.Eq(c.StorageDir+"/"+p.FileName),
						gomock.Eq(make([]byte, 1)),
						gomock.Eq(fs.FileMode(0644)),
					).
					Return(fmt.Errorf("access denied")).
					Times(1)
				res, err := s.UploadFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("failed storing file")))
			})
		})

		When("success upload file", func() {
			It("should return result", func() {
				fm.EXPECT().
					IsExists(gomock.Eq(c.StorageDir)).
					Return(true).
					Times(1)

				fm.EXPECT().
					IsExists(gomock.Eq(c.StorageDir + "/" + p.FileName)).
					Return(false).
					Times(1)

				fm.EXPECT().
					WriteFile(
						gomock.Eq(c.StorageDir+"/"+p.FileName),
						gomock.Eq(make([]byte, 1)),
						gomock.Eq(fs.FileMode(0644)),
					).
					Return(nil).
					Times(1)
				clo.EXPECT().Now().Return(currentTime)

				res, err := s.UploadFile(ctx, p)

				eRes := &goseidon.UploadFileResult{
					FileName:   p.FileName,
					UploadedAt: time.Now(),
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
			c   *local.LocalConfig
			fm  *io.MockFileManager
		)

		BeforeEach(func() {
			ctx = context.Background()
			p = goseidon.RetrieveFileParam{
				Id: "unique-access-id",
			}
			c = &local.LocalConfig{
				StorageDir: "storage",
			}
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			fm = io.NewMockFileManager(ctrl)
			storage, _ := local.NewLocalStorage(c)
			storage.FileManager = fm
			s = storage
		})

		When("context is invalid", func() {
			It("should return error", func() {
				res, err := s.RetrieveFile(nil, p)

				Expect(res).To(BeNil())
				Expect(err.Error()).To(Equal("invalid context"))
			})
		})

		When("file is not available", func() {
			It("should return error", func() {
				fm.EXPECT().
					IsExists(gomock.Eq(c.StorageDir + "/" + p.Id)).
					Return(false).
					Times(1)
				res, err := s.RetrieveFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("file is not found")))
			})
		})

		When("failed open file", func() {
			It("should return error", func() {
				fm.EXPECT().
					IsExists(gomock.Eq(c.StorageDir + "/" + p.Id)).
					Return(true).
					Times(1)

				fm.EXPECT().
					Open(gomock.Eq(c.StorageDir+"/"+p.Id)).
					Return(nil, fmt.Errorf("access denied")).
					Times(1)

				res, err := s.RetrieveFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("failed open file")))
			})
		})

		When("failed read file", func() {
			It("should return error", func() {
				fm.EXPECT().
					IsExists(gomock.Eq(c.StorageDir + "/" + p.Id)).
					Return(true).
					Times(1)

				file := &os.File{}
				fm.EXPECT().
					Open(gomock.Eq(c.StorageDir+"/"+p.Id)).
					Return(file, nil).
					Times(1)

				fm.EXPECT().
					ReadFile(gomock.Eq(file)).
					Return(nil, fmt.Errorf("failed read file")).
					Times(1)

				res, err := s.RetrieveFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("failed read file")))
			})
		})

		When("success retrieve file", func() {
			It("should return result", func() {
				fm.EXPECT().
					IsExists(gomock.Eq(c.StorageDir + "/" + p.Id)).
					Return(true).
					Times(1)

				file := &os.File{}
				fm.EXPECT().
					Open(gomock.Eq(c.StorageDir+"/"+p.Id)).
					Return(file, nil).
					Times(1)

				binFile := make([]byte, 1)
				fm.EXPECT().
					ReadFile(gomock.Eq(file)).
					Return(binFile, nil).
					Times(1)

				res, err := s.RetrieveFile(ctx, p)

				eRes := &goseidon.RetrieveFileResult{
					File: binFile,
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
			c   *local.LocalConfig
			fm  *io.MockFileManager
		)

		BeforeEach(func() {
			ctx = context.Background()
			p = goseidon.DeleteFileParam{
				Id: "unique-access-id",
			}
			c = &local.LocalConfig{
				StorageDir: "storage",
			}
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			fm = io.NewMockFileManager(ctrl)
			storage, _ := local.NewLocalStorage(c)
			storage.FileManager = fm
			s = storage
		})

		When("context is invalid", func() {
			It("should return error", func() {
				res, err := s.DeleteFile(nil, p)

				Expect(res).To(BeNil())
				Expect(err.Error()).To(Equal("invalid context"))
			})
		})

		When("file is not found", func() {
			It("should return error", func() {
				fm.EXPECT().
					IsExists(gomock.Eq(c.StorageDir + "/" + p.Id)).
					Return(false).
					Times(1)

				res, err := s.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("file is not found")))
			})
		})

		When("failed remove file", func() {
			It("should return error", func() {
				fm.EXPECT().
					IsExists(gomock.Eq(c.StorageDir + "/" + p.Id)).
					Return(true).
					Times(1)

				fm.EXPECT().
					RemoveFile(gomock.Eq(c.StorageDir + "/" + p.Id)).
					Return(fmt.Errorf("invalid permission")).
					Times(1)

				res, err := s.DeleteFile(ctx, p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("failed delete file")))
			})
		})

		When("success remove file", func() {
			It("should return result", func() {
				currentTime := time.Now()

				fm.EXPECT().
					IsExists(gomock.Eq(c.StorageDir + "/" + p.Id)).
					Return(true).
					Times(1)

				fm.EXPECT().
					RemoveFile(gomock.Eq(c.StorageDir + "/" + p.Id)).
					Return(nil).
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
