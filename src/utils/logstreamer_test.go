package utils

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockWriter struct{}

func (m *MockWriter) Write(p []byte) (n int, err error) {
	return
}

func TestFush(t *testing.T) {
	logger := log.New(&MockWriter{}, "stdout", log.Ldate|log.Ltime)
	l := NewLogstreamer(logger, "", true)

	l.Write([]byte("hello"))

	l.Flush()

	output := l.FlushRecord()
	assert.Equal(t, output, "hello")

	clearedOutput := l.FlushRecord()
	assert.Equal(t, clearedOutput, "")
}

func TestOutputLines(t *testing.T) {
	logger := log.New(&MockWriter{}, "stdout", log.Ldate|log.Ltime)
	l := NewLogstreamer(logger, "", true)

	l.Write([]byte("hello\n"))
	l.Write([]byte("world\n"))
	l.Write([]byte("not there"))

	output := l.FlushRecord()
	assert.Equal(t, "hello\nworld\n", output)
}
