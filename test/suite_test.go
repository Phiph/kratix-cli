package integration_test

import (
	"os/exec"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var binaryPath string

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Suite")
}

var _ = BeforeSuite(func() {
	var err error
	binaryPath, err = gexec.Build("github.com/syntasso/kratix-cli")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

type runner struct {
	exitCode int
	dir      string
	flags    map[string]string
	timeout  time.Duration
}

func withExitCode(exitCode int) *runner {
	return &runner{exitCode: exitCode}
}

func (r *runner) run(args ...string) *gexec.Session {
	for key, value := range r.flags {
		if value == "" {
			args = append(args, key)
		} else {
			args = append(args, key, value)
		}
	}
	cmd := exec.Command(binaryPath, args...)
	cmd.Dir = r.dir
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	t := 1 * time.Second
	if r.timeout > 0 {
		t = r.timeout
	}
	EventuallyWithOffset(1, session).WithTimeout(t).Should(gexec.Exit(r.exitCode))
	return session
}
