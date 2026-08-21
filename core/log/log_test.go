package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeLogger struct {
	messages []string
}

func (f *fakeLogger) Log(msg string) {
	f.messages = append(f.messages, msg)
}

// resetLogger clears the package-level singleton between subtests, since
// SetLogger only assigns on the first call.
func resetLogger(t *testing.T) {
	t.Helper()
	logger = nil
	t.Cleanup(func() { logger = nil })
}

func TestLog(t *testing.T) {
	t.Run("no-op when no logger is set", func(t *testing.T) {
		resetLogger(t)
		is := assert.New(t)

		is.NotPanics(func() { Log("hello") })
	})

	t.Run("forwards messages to the configured logger", func(t *testing.T) {
		resetLogger(t)
		is := assert.New(t)

		fake := &fakeLogger{}
		SetLogger(fake)

		Log("first message")
		Log("second message")

		is.Equal([]string{"first message", "second message"}, fake.messages)
	})
}

func TestSetLogger(t *testing.T) {
	t.Run("sets the logger when none is configured", func(t *testing.T) {
		resetLogger(t)
		is := assert.New(t)

		fake := &fakeLogger{}
		SetLogger(fake)

		is.Equal(ILogger(fake), logger)
	})

	t.Run("does not overwrite an already configured logger", func(t *testing.T) {
		resetLogger(t)
		is := assert.New(t)

		first := &fakeLogger{}
		second := &fakeLogger{}

		SetLogger(first)
		SetLogger(second)

		is.Equal(ILogger(first), logger)
	})
}
