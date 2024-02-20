package log_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/wh8199/log"
)

func TestInfo(t *testing.T) {
	buf := &bytes.Buffer{}
	testMessage := "Test Message"

	logging := log.NewLogging("test", log.INFO_LEVEL, 4)
	logging.SetOutPut(buf)

	buf.Reset()
	logging.Info(testMessage)

	if !strings.Contains(buf.String(), testMessage) {
		t.Error("test info failed")
		return
	}

	buf.Reset()
	logging.Infof(testMessage)
	if !strings.Contains(buf.String(), testMessage) {
		t.Error("test info failed")
		return
	}

	buf.Reset()
	logging.Infof("Test %s", "Message")

	if !strings.Contains(buf.String(), testMessage) {
		t.Error("test info failed")
		return
	}
}

func TestLevel(t *testing.T) {
	buf := &bytes.Buffer{}
	testMessage := "Test Message"

	logging := log.NewLogging("test", log.INFO_LEVEL, 2)
	logging.SetOutPut(buf)

	buf.Reset()
	logging.Info(testMessage)
	if buf.Len() == 0 {
		t.Error("test level failed")
		return
	}

	buf.Reset()
	logging.Debug(testMessage)
	if buf.Len() != 0 {
		t.Error("test level failed")
		return
	}
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
