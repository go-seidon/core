package aws_s3_test

import (
	"fmt"
	"os"

	aws_s3 "github.com/go-seidon/core/pkg/aws-s3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Storage Option", func() {
	Context("With static credential option", func() {
		When("success apply option", func() {
			It("should return nil", func() {
				cfg := &aws_s3.AwsS3Config{}
				opt := aws_s3.WithStatisCredential(
					"mock-region", "mock-access-key",
					"mock-secret-key", "mock-bucket-name",
				)

				err := opt.Apply(cfg)

				Expect(err).To(BeNil())
			})
		})

		When("failed apply option", func() {
			It("should return nil", func() {
				os.Setenv("AWS_SDK_LOAD_CONFIG", "true")
				os.Setenv("AWS_S3_USE_ARN_REGION", "invalid_value")

				cfg := &aws_s3.AwsS3Config{}
				opt := aws_s3.WithStatisCredential(
					"mock-region", "mock-access-key",
					"mock-secret-key", "mock-bucket-name",
				)

				err := opt.Apply(cfg)

				Expect(err).To(Equal(fmt.Errorf("failed to load environment config, invalid value for environment variable, AWS_S3_USE_ARN_REGION=invalid_value, need true or false")))
			})
		})
	})
})
