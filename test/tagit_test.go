package deploy_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"tagit/cmd/tagit"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

var execCommand = exec.Command

func TestGetStagingPrefix(t *testing.T) {
	tests := []struct {
		env      string
		expected string
	}{
		{"s0", "s0v"},
		{"s1", "s1v"},
		{"s10", "s10v"},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, tagit.GetStagingPrefix(test.env))
	}

	assert.Panics(t, func() { tagit.GetStagingPrefix("invalid") })
}

func TestLatestTag(t *testing.T) {
	// Test case when there are no tags
	execCommand = func(name string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestHelperProcess", "--", name}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1", "GIT_TAG_OUTPUT="}
		return cmd
	}

	defer func() { execCommand = exec.Command }()

	latestTag := tagit.LatestTag("v")
	assert.Equal(t, "", latestTag)

	// Test case when there are tags
	// execCommand = func(name string, args ...string) *exec.Cmd {
	// 	cs := []string{"-test.run=TestHelperProcess", "--", name}
	// 	cs = append(cs, args...)
	// 	cmd := exec.Command(os.Args[0], cs...)
	// 	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1", "GIT_TAG_OUTPUT=v1.0.0\nv1.0.1\n"}
	// 	return cmd
	// }

	// latestTag = tagit.LatestTag("v")
	// assert.Equal(t, "v1.0.1", latestTag)
}

func TestIncrementVersion(t *testing.T) {
	tests := []struct {
		version     string
		releaseType string
		expected    string
	}{
		{"1.2.3", "x", "2.0.0"},
		{"1.2.3", "y", "1.3.0"},
		{"1.2.3", "z", "1.2.4"},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, tagit.IncrementVersion(test.version, test.releaseType))
	}
}

func TestDeployCmd(t *testing.T) {
	// Mock exec.Command
	execCommand = func(name string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestHelperProcess", "--", name}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
		return cmd
	}

	defer func() { execCommand = exec.Command }()

	createRootCmd := func() *cobra.Command {
		rootCmd := &cobra.Command{Use: "root"}
		rootCmd.AddCommand(tagit.TagCmd)
		return rootCmd
	}

	output := captureOutput(func() {
		rootCmd := createRootCmd()
		rootCmd.SetArgs([]string{"tag"})
		rootCmd.Execute()
	})
	assert.Contains(t, output, "Error: You must specify both --env and --type")

	output = captureOutput(func() {
		rootCmd := createRootCmd()
		rootCmd.SetArgs([]string{"tag", "--env=invalid", "--type=1"})
		rootCmd.Execute()
	})
	assert.Contains(t, output, "Error: Invalid environment specified")

	output = captureOutput(func() {
		rootCmd := createRootCmd()
		rootCmd.SetArgs([]string{"tag", "--env=s1", "--type=1"})
		rootCmd.Execute()
	})
	assert.Contains(t, output, "Fetching latest tags")
	assert.Contains(t, output, "New tag created: s1v2.0.0")
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	switch os.Args[3] {
	case "tag":
		if len(os.Args) > 4 && os.Args[4] == "-l" {
			// Simulate `git tag -l 'v*'`
			output := os.Getenv("GIT_TAG_OUTPUT")
			fmt.Fprint(os.Stdout, output)
		}
	case "fetch":
		// Simulate `git fetch --tags`
	case "push":
		// Simulate `git push origin v2.0.0`
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[3])
	}

	os.Exit(0)
}

func captureOutput(f func()) string {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout

	return string(out)
}
