package log_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wh8199/log"
)

func TestInfo(t *testing.T) {
	buf := &bytes.Buffer{}
	testMessage := "Test Message"

	logging := log.NewLogging("test", log.INFO_LEVEL, 2)
	logging.SetOutPut(buf)

	buf.Reset()
	logging.Info(testMessage)
	assert.Contains(t, buf.String(), testMessage)

	buf.Reset()
	logging.Infof(testMessage)
	assert.Contains(t, buf.String(), testMessage)

	buf.Reset()
	logging.Infof("Test %s", "Message")
	assert.Contains(t, buf.String(), testMessage)
}

func TestLevel(t *testing.T) {
	buf := &bytes.Buffer{}
	testMessage := "Test Message"

	logging := log.NewLogging("test", log.INFO_LEVEL, 2)
	logging.SetOutPut(buf)

	buf.Reset()
	logging.Info(testMessage)
	assert.Contains(t, buf.String(), testMessage)

	buf.Reset()
	logging.Debug(testMessage)
	assert.Equal(t, buf.String(), "")
}

func BenchmarkInfo(b *testing.B) {
	testMessage := "Test Message"
	buf := &bytes.Buffer{}

	logging := log.NewLogging("test", log.INFO_LEVEL, 2)
	logging.SetOutPut(buf)

	for i := 0; i < b.N; i++ {
		buf.Reset()
		logging.Info(testMessage)
	}
}
