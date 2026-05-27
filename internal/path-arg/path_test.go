package patharg

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	require "github.com/stretchr/testify/require"
)

func chdir(t *testing.T, dir string) {
	t.Helper()
	cwd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() {
		_ = os.Chdir(cwd)
	})
}

func TestExistsPath_Set_ExistingFile_SetsValueAndReturnsAbsPath(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	rel := "file.txt"
	require.NoError(t, os.WriteFile(rel, []byte("x"), 0o600))

	var ep ExistsPath
	require.NoError(t, ep.Set(rel))
	require.Equal(t, rel, ep.String())
	require.Equal(t, "exists_path", ep.Type())

	expected, err := filepath.Abs(rel)
	require.NoError(t, err)
	require.Equal(t, expected, ep.Path())
}

func TestExistsPath_Set_NonExistent_ReturnsWrappedErrNotExist(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	var ep ExistsPath
	err := ep.Set("missing.txt")
	require.Error(t, err)
	require.ErrorIs(t, err, os.ErrNotExist)
}

func TestExistsDir_Set_ExistingDir_SetsValueAndReturnsAbsPath(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	rel := "subdir"
	require.NoError(t, os.Mkdir(rel, 0o700))

	var ed ExistsDir
	require.NoError(t, ed.Set(rel))
	require.Equal(t, rel, ed.String())
	require.Equal(t, "exists_dir", ed.Type())

	expected, err := filepath.Abs(rel)
	require.NoError(t, err)
	require.Equal(t, expected, ed.Path())
}

func TestExistsDir_Set_PathIsFile_ReturnsNotDirectoryError(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	rel := "file.txt"
	require.NoError(t, os.WriteFile(rel, []byte("x"), 0o600))

	var ed ExistsDir
	err := ed.Set(rel)
	require.Error(t, err)
	require.EqualError(t, err, "convert \"file.txt\" to dir: not a directory")
}

func TestExistsDir_Set_NonExistent_ReturnsWrappedErrNotExist(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	var ed ExistsDir
	err := ed.Set("missing")
	require.Error(t, err)
	require.True(t, errors.Is(err, os.ErrNotExist))
}

func TestExistsDir_Set_InvalidPath_ReturnsErrorNotPanic(t *testing.T) {
	var ed ExistsDir
	bad := string([]byte{0})

	require.NotPanics(t, func() {
		err := ed.Set(bad)
		require.Error(t, err)
	})
}
