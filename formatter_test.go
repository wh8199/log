package log_test

import (
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/wh8199/log"
)

func TestDefaultFormatter(t *testing.T) {
	bufferPool := log.NewBufferPool()
	_, caller, line, _ := runtime.Caller(1)

	testMessage := "Test"
	buf := log.DefaultFormater("", log.INFO_LEVEL, 2, bufferPool, testMessage)

	t.Log(buf.String())

	if strings.Contains(buf.String(), caller+":"+strconv.Itoa(line)+" "+"Info"+" msg: "+testMessage+"\n") {
		t.Error("test failed for default formatter")
		return
	}

	buf1 := log.DefaultFormater("", log.INFO_LEVEL, 2, bufferPool, "Test %d", 1)
	s1 := "Test 1"
	if strings.Contains(buf1.String(), caller+":"+strconv.Itoa(line)+" "+"Info"+" msg: "+s1+"\n") {
		t.Error("test failed for default formatter")
		return
	}

	buf2 := log.DefaultFormater("", log.INFO_LEVEL, 2, bufferPool, "Test %.1f", 1.2)
	s2 := "Test 1.2"
	if strings.Contains(buf2.String(), caller+":"+strconv.Itoa(line)+" "+"Info"+" msg: "+s2+"\n") {
		t.Error("test failed for default formatter")
		return
	}
}
