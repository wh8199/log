package log_test

import (
	"runtime"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wh8199/log"
)

func TestDefaultFormatter(t *testing.T) {
	bufferPool := log.NewBufferPool()
	_, caller, line, _ := runtime.Caller(1)

	testMessage := "Test"
	buf := log.DefaultFormater("test", "Info", 2, bufferPool, testMessage)
	assert.Contains(t, buf.String(), caller+":"+strconv.Itoa(line)+" "+"Info"+" msg: "+testMessage+"\n")

	buf1 := log.DefaultFormater("test", "Info", 2, bufferPool, "Test %d", 1)
	s1 := "Test 1"
	assert.Contains(t, buf1.String(), caller+":"+strconv.Itoa(line)+" "+"Info"+" msg: "+s1+"\n")

	buf2 := log.DefaultFormater("test", "Info", 2, bufferPool, "Test %.1f", 1.2)
	s2 := "Test 1.2"
	assert.Contains(t, buf2.String(), caller+":"+strconv.Itoa(line)+" "+"Info"+" msg: "+s2+"\n")
}
