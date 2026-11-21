package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testConfig struct {
	Name  string
	Count int
}

func newTestConfig(t *testing.T, o ...Option) *Config {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("AppData", dir)
	t.Setenv("HOME", dir)

	c, err := New(append([]Option{DirectoryName("test")}, o...)...)
	require.NoError(t, err)
	return c
}

func TestSave_ReplacesExistingFile(t *testing.T) {
	c := newTestConfig(t)

	require.NoError(t, c.Save(testConfig{Name: "a much longer first value", Count: 1}))
	require.NoError(t, c.Save(testConfig{Name: "short", Count: 2}))

	var got testConfig
	require.NoError(t, c.Load(&got))
	assert.Equal(t, testConfig{Name: "short", Count: 2}, got)

	entries, err := os.ReadDir(c.Directory())
	require.NoError(t, err)
	assert.Len(t, entries, 1, "temporary files should not be left behind")
}

func TestSave_SetsConfiguredModeOnExistingFile(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("file modes are not meaningful on windows")
	}

	c := newTestConfig(t, Permissions(0700, 0600))
	path := filepath.Join(c.Directory(), "config.yml")

	require.NoError(t, os.MkdirAll(c.Directory(), 0700))
	require.NoError(t, os.WriteFile(path, []byte("name: old\n"), 0644))
	require.NoError(t, os.Chmod(path, 0644))

	require.NoError(t, c.Save(testConfig{Name: "new"}))

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())
}
