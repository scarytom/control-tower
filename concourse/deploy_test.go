package concourse_test

import (
	"fmt"
	"io"
	"os"

	. "bitbucket.org/engineerbetter/concourse-up/concourse"
	"bitbucket.org/engineerbetter/concourse-up/config"
	"bitbucket.org/engineerbetter/concourse-up/terraform"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Deploy", func() {
	It("Generates the correct terraform infrastructure", func() {
		var appliedTFConfig []byte

		client := &FakeConfigClient{}
		client.FakeLoadOrCreate = func(project string) (*config.Config, error) {
			return &config.Config{
				PublicKey:   "example-public-key",
				PrivateKey:  "example-private-key",
				Region:      "eu-west-1",
				Deployment:  fmt.Sprintf("concourse-up-%s", project),
				Project:     project,
				TFStatePath: "example-path",
			}, nil
		}

		applied := false
		cleanedUp := false

		clientFactory := func(config []byte, stdout, stderr io.Writer) (terraform.IClient, error) {
			appliedTFConfig = config
			return &FakeTerraformClient{
				FakeApply: func() error {
					applied = true
					return nil
				},
				FakeOutput: func() (*terraform.Metadata, error) {
					return &terraform.Metadata{
						BoshDBPort: terraform.MetadataStringValue{
							Value: "5432",
						},
					}, nil
				},
				FakeCleanup: func() error {
					cleanedUp = true
					return nil
				},
			}, nil
		}

		err := Deploy("happymeal", "eu-west-1", clientFactory, client, os.Stdout, os.Stderr)
		Expect(err).ToNot(HaveOccurred())

		Expect(string(appliedTFConfig)).To(ContainSubstring("concourse-up-happymeal"))
		Expect(applied).To(BeTrue())
		Expect(cleanedUp).To(BeTrue())
	})
})