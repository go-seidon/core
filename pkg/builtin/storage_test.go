package builtin_test

import (
	"testing"

	goseidon "github.com/go-seidon/core"
	"github.com/go-seidon/core/pkg/builtin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBuiltin(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Builtin Package")
}

var _ = Describe("Storage", func() {
	Context("NewStorage function", func() {
		When("success create storage", func() {
			It("should return builtin storage", func() {
				s, err := builtin.NewStorage()

				Expect(s).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})
	})

	Context("UploadFile method", func() {
		var (
			s goseidon.Storage
			p goseidon.UploadFileParam
		)

		BeforeEach(func() {
			storage, _ := builtin.NewStorage()
			s = storage
			p = goseidon.UploadFileParam{}
		})

		When("failed upload file", func() {
			It("should return error result", func() {
				res, err := s.UploadFile(p)

				Expect(res).To(BeNil())
				Expect(err.Error()).To(Equal("failed upload file"))
			})
		})
	})

	Context("RetrieveFile method", func() {
		var (
			s goseidon.Storage
			p goseidon.RetrieveFileParam
		)

		BeforeEach(func() {
			storage, _ := builtin.NewStorage()
			s = storage
			p = goseidon.RetrieveFileParam{}
		})

		When("failed retrieve file", func() {
			It("should return error result", func() {
				res, err := s.RetrieveFile(p)

				Expect(res).To(BeNil())
				Expect(err.Error()).To(Equal("failed retrieve file"))
			})
		})
	})
})
