package registry_test

import (
	"context"
	"os"
	"time"

	"github.com/apivzero/watchtower/internal/actions/mocks"
	unit "github.com/apivzero/watchtower/pkg/registry"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Registry", func() {
	Describe("WarnOnAPIConsumption", func() {
		When("Given a container with an image from ghcr.io", func() {
			It("should want to warn", func() {
				Expect(testContainerWithImage("ghcr.io/containrrr/watchtower")).To(BeTrue())
			})
		})
		When("Given a container with an image implicitly from dockerhub", func() {
			It("should want to warn", func() {
				Expect(testContainerWithImage("docker:latest")).To(BeTrue())
			})
		})
		When("Given a container with an image explicitly from dockerhub", func() {
			It("should want to warn", func() {
				Expect(testContainerWithImage("index.docker.io/docker:latest")).To(BeTrue())
				Expect(testContainerWithImage("docker.io/docker:latest")).To(BeTrue())
			})
		})
		When("Given a container with an image from some other registry", func() {
			It("should not want to warn", func() {
				Expect(testContainerWithImage("docker.fsf.org/docker:latest")).To(BeFalse())
				Expect(testContainerWithImage("altavista.com/docker:latest")).To(BeFalse())
				Expect(testContainerWithImage("gitlab.com/docker:latest")).To(BeFalse())
			})
		})
	})
})

func testContainerWithImage(imageName string) bool {
	container := mocks.CreateMockContainer("", "", imageName, time.Now())
	return unit.WarnOnAPIConsumption(container)
}

var _ = Describe("GetPullOptions", func() {
	When("no authentication is configured", func() {
		It("should return empty pull options without error", func() {
			_ = os.Unsetenv("REPO_USER")
			_ = os.Unsetenv("REPO_PASS")

			tmpDir, err := os.MkdirTemp("", "watchtower-test-*")
			Expect(err).NotTo(HaveOccurred())
			defer os.RemoveAll(tmpDir)
			_ = os.Setenv("DOCKER_CONFIG", tmpDir)
			defer os.Unsetenv("DOCKER_CONFIG")

			opts, err := unit.GetPullOptions("docker.io/library/alpine:latest")
			Expect(err).NotTo(HaveOccurred())
			Expect(opts.RegistryAuth).To(BeEmpty())
			Expect(opts.PrivilegeFunc).To(BeNil())
		})
	})
	When("env auth credentials are provided", func() {
		It("should return pull options with auth and privilege func", func() {
			_ = os.Setenv("REPO_USER", "testuser")
			_ = os.Setenv("REPO_PASS", "testpass")
			defer func() {
				_ = os.Unsetenv("REPO_USER")
				_ = os.Unsetenv("REPO_PASS")
			}()

			opts, err := unit.GetPullOptions("docker.io/library/alpine:latest")
			Expect(err).NotTo(HaveOccurred())
			Expect(opts.RegistryAuth).NotTo(BeEmpty())
			Expect(opts.PrivilegeFunc).NotTo(BeNil())
		})
	})
})

var _ = Describe("DefaultAuthHandler", func() {
	It("should accept a context parameter and return empty string with no error", func() {
		result, err := unit.DefaultAuthHandler(context.Background())
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(BeEmpty())
	})
})
