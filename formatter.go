package log

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
)

type Formatter func(name, level string, callerLevel int, pool *BufferPool, format string, args ...interface{}) *bytes.Buffer

func DefaultFormater(name, level string, callerLevel int, pool *BufferPool, format string, args ...interface{}) *bytes.Buffer {
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

	result := "[ " + name + " ] " + time + " " + caller + ":" + strconv.Itoa(line) + " " + level + " msg: " + s + "\n"

	buf := pool.Get()
	buf.WriteString(result)

	return buf
}
