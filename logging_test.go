package log_test

import (
	"bytes"
	"testing"

	"github.com/wh8199/log"
)

func TestInfo(t *testing.T) {
	buf := &bytes.Buffer{}
	testMessage := "Test Message"

	logging := log.NewLogging("test", log.INFO_LEVEL, 2)
	logging.SetOutPut(buf)

	buf.Reset()
	logging.Info(testMessage)

	buf.Reset()
	logging.Infof(testMessage)

	buf.Reset()
	logging.Infof("Test %s", "Message")
}

func TestLevel(t *testing.T) {
	buf := &bytes.Buffer{}
	testMessage := "Test Message"

	logging := log.NewLogging("test", log.INFO_LEVEL, 2)
	logging.SetOutPut(buf)

	buf.Reset()
	logging.Info(testMessage)

	buf.Reset()
	logging.Debug(testMessage)
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
