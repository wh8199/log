package log

import (
	"bytes"
	"strings"
	"testing"
)

func TestInfo(t *testing.T) {
	buf := &bytes.Buffer{}
	testMessage := "Test Message"

	logging := NewLogging("test", INFO_LEVEL, 4)
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

	logging := NewLogging("test", INFO_LEVEL, 2)
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

func TestUpdateOutput(t *testing.T) {
	l := NewLogging("test", DEBUG_LEVEL, 2)

	buf := bytes.Buffer{}
	buf1 := bytes.Buffer{}

	l.SetOutPut(&buf)
	l.Info("buffer1")

	l.SetOutPut(&buf1)
	l.Info("buffer2")

	if !strings.Contains(buf.String(), "buffer1") {
		t.Error("update output failed")
		return
	}

	if !strings.Contains(buf1.String(), "buffer2") {
		t.Error("update output failed")
		return
	}
}

func TestPrintLevel(t *testing.T) {
	buf := &bytes.Buffer{}
	testMessage := "Test Message"

	logging := NewLogging("test", INFO_LEVEL, 2)
	logging.SetOutPut(buf)

	buf.Reset()
	logging.Info(testMessage)
	if buf.Len() == 0 {
		t.Error("test level failed")
		return
	}

	if !strings.Contains(buf.String(), "Info") {
		t.Error("test log level failed")
		return
	}
}
