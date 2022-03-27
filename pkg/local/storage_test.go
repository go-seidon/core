package local_test

import (
	"fmt"
	"io/fs"
	"os"
	"testing"

	goseidon "github.com/go-seidon/core"
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
			c  *local.LocalConfig
			fm io.FileManagerService
		)

		BeforeEach(func() {
			c = &local.LocalConfig{
				StorageDir: "storage",
			}
			fm, _ = io.NewFileManager()
		})

		When("config is invalid", func() {
			It("should return error", func() {
				s, err := local.NewLocalStorage(nil, fm)

				Expect(s).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("invalid storage config")))
			})
		})

		When("file manager is invalid", func() {
			It("should return error", func() {
				s, err := local.NewLocalStorage(c, nil)

				Expect(s).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("invalid file manager")))
			})
		})

		When("success create storage", func() {
			It("should return local storage", func() {
				s, err := local.NewLocalStorage(c, fm)

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
			s  goseidon.Storage
			p  goseidon.UploadFileParam
			c  *local.LocalConfig
			fm *io.MockFileManagerService
		)

		BeforeEach(func() {
			p = goseidon.UploadFileParam{
				FileName: "image.jpg",
				FileData: make([]byte, 1),
			}
			c = &local.LocalConfig{
				StorageDir: "storage",
			}
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			fm = io.NewMockFileManagerService(ctrl)
			storage, _ := local.NewLocalStorage(c, fm)
			s = storage
		})

		When("file already exists", func() {
			It("should return error", func() {
				fm.EXPECT().
					IsFileExists(gomock.Eq(c.StorageDir + "/" + p.FileName)).
					Return(true).
					Times(1)
				res, err := s.UploadFile(p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("file already exists")))
			})
		})

		When("failed upload file", func() {
			It("should return error", func() {
				fm.EXPECT().
					IsFileExists(gomock.Eq(c.StorageDir + "/" + p.FileName)).
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
				res, err := s.UploadFile(p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("failed storing file")))
			})
		})

		When("success upload file", func() {
			It("should return result", func() {
				fm.EXPECT().
					IsFileExists(gomock.Eq(c.StorageDir + "/" + p.FileName)).
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
				res, err := s.UploadFile(p)

				eRes := &goseidon.UploadFileResult{FileName: p.FileName}
				Expect(res).To(Equal(eRes))
				Expect(err).To(BeNil())
			})
		})
	})

	Context("RetrieveFile method", func() {
		var (
			s  goseidon.Storage
			p  goseidon.RetrieveFileParam
			c  *local.LocalConfig
			fm *io.MockFileManagerService
		)

		BeforeEach(func() {
			p = goseidon.RetrieveFileParam{
				Id: "unique-access-id",
			}
			c = &local.LocalConfig{
				StorageDir: "storage",
			}
			t := GinkgoT()
			ctrl := gomock.NewController(t)
			fm = io.NewMockFileManagerService(ctrl)
			storage, _ := local.NewLocalStorage(c, fm)
			s = storage
		})

		When("file is not available", func() {
			It("should return error", func() {
				fm.EXPECT().
					IsFileExists(gomock.Eq(c.StorageDir + "/" + p.Id)).
					Return(false).
					Times(1)
				res, err := s.RetrieveFile(p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("file is not found")))
			})
		})

		When("failed open file", func() {
			It("should return error", func() {
				fm.EXPECT().
					IsFileExists(gomock.Eq(c.StorageDir + "/" + p.Id)).
					Return(true).
					Times(1)

				fm.EXPECT().
					Open(gomock.Eq(c.StorageDir+"/"+p.Id)).
					Return(nil, fmt.Errorf("access denied")).
					Times(1)

				res, err := s.RetrieveFile(p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("failed open file")))
			})
		})

		When("failed read file", func() {
			It("should return error", func() {
				fm.EXPECT().
					IsFileExists(gomock.Eq(c.StorageDir + "/" + p.Id)).
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

				res, err := s.RetrieveFile(p)

				Expect(res).To(BeNil())
				Expect(err).To(Equal(fmt.Errorf("failed read file")))
			})
		})

		When("success retrieve file", func() {
			It("should return result", func() {
				fm.EXPECT().
					IsFileExists(gomock.Eq(c.StorageDir + "/" + p.Id)).
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

				res, err := s.RetrieveFile(p)

				eRes := &goseidon.RetrieveFileResult{
					File: binFile,
				}
				Expect(res).To(Equal(eRes))
				Expect(err).To(BeNil())
			})
		})
	})
})
