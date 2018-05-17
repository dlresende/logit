package git_test

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	git "github.com/dlresende/logit/git"
)

var _ = Describe("Git", func() {
	var gitDir string

	BeforeEach(func() {
		gitDir, _ = ioutil.TempDir("", "")
	})

	It("should create a new git repo with inital commit when no repo exists", func() {
		git.Init(gitDir)

		output, err := run("git", "-C", gitDir, "log")

		Expect(err).To(Not(HaveOccurred()))
		Expect(output).To(ContainSubstring("Initial commit"))
	})

	It("should return the existing git repo when a repo already exists", func() {
		run("git", "-C", gitDir, "init")

		_, err := git.Init(gitDir)

		Expect(err).To(Not(HaveOccurred()))
	})

	It("should commit given a valid repo", func() {
		repo, _ := git.Init(gitDir)
		now := time.Now()

		repo.Commit("test commit", "ginkgo", "master", now)

		output, err := run("git", "-C", gitDir, "log")
		Expect(err).To(Not(HaveOccurred()))
		Expect(output).To(ContainSubstring("test commit"))
	})
})

func run(command string, arguments ...string) (output string, err error) {
	cmd := exec.Command(command, arguments...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	return out.String(), err
}
