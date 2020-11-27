package log

import (
	"bytes"
	"fmt"
	"runtime"
)

type Formatter func(name string, level LoggingLevel, callerLevel int, pool *BufferPool, format string, args ...interface{}) *bytes.Buffer

func DefaultFormater(name string, level LoggingLevel, callerLevel int, pool *BufferPool, format string, args ...interface{}) *bytes.Buffer {
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

	result := fmt.Sprintf("[ %s ] %s %s:%d %v msg: %s\n",
		name, time, caller, line, level, s)

	buf := pool.Get()
	buf.WriteString(result)

	return buf
}
