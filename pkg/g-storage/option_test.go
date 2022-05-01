package g_storage_test

import (
	"fmt"
	"os"

	"cloud.google.com/go/storage"
	g_storage "github.com/go-seidon/core/pkg/g-storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Storage Option", func() {

	Context("With google client option", func() {
		When("success apply option", func() {
			It("should return nil", func() {
				cfg := &g_storage.GoogleConfig{}
				opt := g_storage.WithGoogleClient("mock-bucket-name", &storage.Client{})
				err := opt.Apply(cfg)

				Expect(err).To(BeNil())
			})
		})
	})

	Context("With credential path option", func() {
		When("failed apply option", func() {
			It("should return error", func() {
				os.Setenv("STORAGE_EMULATOR_HOST", " ") //invalid host

				cfg := &g_storage.GoogleConfig{}
				opt := g_storage.WithCredentialPath("mock-bucket-name", "mock-credential-path")
				err := opt.Apply(cfg)

				Expect(err).To(Equal(fmt.Errorf("dialing: options.WithoutAuthentication is incompatible with any option that provides credentials")))
			})
		})
	})
})
