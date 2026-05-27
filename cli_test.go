package main

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	profile "github.com/karetskiiVO/ProfileCLIUtility/pkg/profile"
	yaml "gopkg.in/yaml.v3"

	require "github.com/stretchr/testify/require"
)

func buildCLI(t *testing.T) string {
	t.Helper()

	out := filepath.Join(t.TempDir(), "profile-utility")
	if runtime.GOOS == "windows" {
		out += ".exe"
	}

	cmd := exec.Command("go", "build", "-o", out, ".")
	cmd.Dir = projectRoot(t)
	b, err := cmd.CombinedOutput()
	require.NoError(t, err, string(b))

	return out
}

func projectRoot(t *testing.T) string {
	t.Helper()

	wd, err := os.Getwd()
	require.NoError(t, err)
	return wd
}

type runResult struct {
	ExitCode int
	Output   string
}

func runCLI(t *testing.T, bin string, args ...string) runResult {
	t.Helper()

	cmd := exec.Command(bin, args...)
	cmd.Dir = projectRoot(t)
	b, err := cmd.CombinedOutput()

	res := runResult{Output: string(b)}
	if err == nil {
		res.ExitCode = 0
		return res
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		res.ExitCode = exitErr.ExitCode()
		return res
	}

	t.Fatalf("unexpected error running %v: %v\noutput:\n%s", append([]string{bin}, args...), err, res.Output)
	return runResult{}
}

func requireOK(t *testing.T, res runResult) {
	t.Helper()
	require.Equal(t, 0, res.ExitCode, res.Output)
}

func requireFail(t *testing.T, res runResult) {
	t.Helper()
	require.NotEqual(t, 0, res.ExitCode, res.Output)
}

func TestCLI(t *testing.T) {
	bin := buildCLI(t)

	t.Run("help", func(t *testing.T) {
		res := runCLI(t, bin, "--help")
		requireOK(t, res)
		require.Contains(t, res.Output, "profile-utility")
	})

	t.Run("profile lifecycle", func(t *testing.T) {
		dir := t.TempDir()

		res := runCLI(t, bin,
			"profile", "create",
			"--path", dir,
			"--name", "dev",
			"--user", "alice",
			"--project", "proj",
		)
		requireOK(t, res)

		yamlPath := filepath.Join(dir, "dev.yaml")
		data, err := os.ReadFile(yamlPath)
		require.NoError(t, err)

		var p profile.Struct
		require.NoError(t, yaml.Unmarshal(data, &p))
		require.Equal(t, "alice", p.User)
		require.Equal(t, "proj", p.Project)

		res = runCLI(t, bin, "profile", "get", "--path", dir, "--name", "dev")
		requireOK(t, res)
		require.Contains(t, res.Output, "Name")
		require.Contains(t, res.Output, "dev")
		require.Contains(t, res.Output, "alice")
		require.Contains(t, res.Output, "proj")

		res = runCLI(t, bin, "profile", "list", "--path", dir)
		requireOK(t, res)
		require.Contains(t, res.Output, "Total: 1")

		res = runCLI(t, bin, "profile", "delete", "--path", dir, "--name", "dev", "-v")
		requireOK(t, res)
		require.Contains(t, res.Output, "deleted successfully")

		res = runCLI(t, bin, "profile", "get", "--path", dir, "--name", "dev")
		requireFail(t, res)
		require.Contains(t, res.Output, "profile not found")
	})

	t.Run("create requires flags", func(t *testing.T) {
		dir := t.TempDir()

		res := runCLI(t, bin, "profile", "create", "--path", dir)
		requireFail(t, res)
		require.Contains(t, strings.ToLower(res.Output), "required")
	})

	t.Run("list strict fails on invalid yaml", func(t *testing.T) {
		dir := t.TempDir()
		bad := []byte("user: alice\nproject: proj\nextra: 1\n")
		require.NoError(t, os.WriteFile(filepath.Join(dir, "bad.yaml"), bad, 0o600))

		res := runCLI(t, bin, "profile", "list", "--path", dir, "--strict")
		requireFail(t, res)
		require.Contains(t, res.Output, "failed to parse profiles")
	})
}
