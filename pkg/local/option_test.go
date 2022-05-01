package local_test

import (
	"fmt"

	"github.com/go-seidon/core/pkg/local"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Storage Option", func() {
	Context("With normal storage dir option", func() {
		When("storage directory is invalid", func() {
			It("should return error", func() {
				cfg := &local.LocalConfig{}
				opt := local.WithNormalStorageDir("")
				err := opt.Apply(cfg)

				Expect(err).To(Equal(fmt.Errorf("invalid storage directory")))
			})
		})

		When("storage directory is valid", func() {
			It("should return nil", func() {
				cfg := &local.LocalConfig{}
				opt := local.WithNormalStorageDir("storage")
				err := opt.Apply(cfg)

				Expect(err).To(BeNil())
			})
		})

		When("storage directory contain capital word", func() {
			It("should be lowercased", func() {
				cfg := &local.LocalConfig{}
				opt := local.WithNormalStorageDir("STORAGE/Sub-Dir")
				err := opt.Apply(cfg)

				Expect(err).To(BeNil())
				Expect(cfg.StorageDir).To(Equal("storage/sub-dir"))
			})
		})

		When("storage directory contain trailing slashes", func() {
			It("should be removed", func() {
				cfg := &local.LocalConfig{}
				opt := local.WithNormalStorageDir("STORAGE/Sub-Dir/")
				err := opt.Apply(cfg)

				Expect(err).To(BeNil())
				Expect(cfg.StorageDir).To(Equal("storage/sub-dir"))
			})
		})
	})
})
