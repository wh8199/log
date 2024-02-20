package log

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
)

type Formatter func(module Module, level LoggingLevel, callerLevel int, pool *BufferPool, format string, args ...interface{}) *bytes.Buffer

func DefaultFormater(module Module, level LoggingLevel, callerLevel int, pool *BufferPool, format string, args ...interface{}) *bytes.Buffer {
	var (
		s string
	)

	if len(args) == 0 {
		s = format
	} else {
		s = fmt.Sprintf(format, args...)
	}

	time := CacheTime()

	_, caller, line, _ := runtime.Caller(callerLevel)
	buf := pool.Get()
	buf.Reset()

	if len(module) != 0 {
		buf.WriteString("[ ")
		buf.WriteString(module.String())
		buf.WriteString(" ] ")
	}

	buf.WriteString(time)
	buf.WriteString(" ")
	buf.WriteString(caller)
	buf.WriteString(":")
	buf.WriteString(strconv.Itoa(line))
	buf.WriteString(" ")
	buf.WriteString(level.String())
	buf.WriteString(" msg: ")
	buf.WriteString(s)
	buf.WriteString("\n")

	return buf
}
