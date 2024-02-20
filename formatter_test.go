package log

import (
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestDefaultFormatter(t *testing.T) {
	_, caller, line, _ := runtime.Caller(0)

	testMessage := "Test"
	buf := DefaultFormater(&LogRecord{
		args:         []interface{}{testMessage},
		callerLevel:  2,
		enableCaller: true,
		logLevel:     INFO_LEVEL,
	})

	if strings.Contains(buf.String(), caller+":"+strconv.Itoa(line)+" "+"Info"+" msg: "+testMessage+"\n") {
		t.Error("test failed for default formatter")
		return
	}

	buf1 := DefaultFormater(&LogRecord{
		format:       "Test %d",
		args:         []interface{}{1},
		callerLevel:  2,
		enableCaller: true,
		logLevel:     INFO_LEVEL,
	})
	s1 := "Test 1"
	if strings.Contains(buf1.String(), caller+":"+strconv.Itoa(line)+" "+"Info"+" msg: "+s1+"\n") {
		t.Error("test failed for default formatter")
		return
	}

	buf2 := DefaultFormater(&LogRecord{
		format:       "Test %.1f",
		args:         []interface{}{1.2},
		callerLevel:  2,
		enableCaller: true,
		logLevel:     INFO_LEVEL,
	})
	s2 := "Test 1.2"
	if strings.Contains(buf2.String(), caller+":"+strconv.Itoa(line)+" "+"Info"+" msg: "+s2+"\n") {
		t.Error("test failed for default formatter")
		return
	}
}
