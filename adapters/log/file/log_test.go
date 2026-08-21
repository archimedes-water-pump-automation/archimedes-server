package file

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	file, err := os.CreateTemp(t.TempDir(), "logger-*.log")
	is.NoError(err)
	defer file.Close()

	l := NewLogger(file)

	is.NotNil(l)
	is.NotNil(l.log)
}

func TestLogger_Log(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	path := filepath.Join(t.TempDir(), "logger.log")
	file, err := os.Create(path)
	is.NoError(err)
	defer file.Close()

	l := NewLogger(file)
	l.Log("hello world")

	contents, err := os.ReadFile(path)
	is.NoError(err)
	is.True(strings.Contains(string(contents), "archimedes:"))
	is.True(strings.Contains(string(contents), "hello world"))
}
