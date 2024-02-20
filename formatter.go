package log

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
)

type Formatter func(logRecord *LogRecord) *bytes.Buffer

func DefaultFormater(logRecord *LogRecord) *bytes.Buffer {
	var (
		s string
	)

	if len(logRecord.format) == 0 {
		s = fmt.Sprint(logRecord.args...)
	} else {
		s = fmt.Sprintf(logRecord.format, logRecord.args...)
	}

	time := CacheTime()

	_, caller, line, _ := runtime.Caller(logRecord.callerLevel)
	buf := pool.Get()
	buf.Reset()

	if len(logRecord.module) != 0 {
		buf.WriteString("[ ")
		buf.WriteString(logRecord.module)
		buf.WriteString(" ] ")
	}

	buf.WriteString(time)
	buf.WriteString(" ")
	buf.WriteString(caller)
	buf.WriteString(":")
	buf.WriteString(strconv.Itoa(line))
	buf.WriteString(" ")
	buf.WriteString(logRecord.logLevel.String())
	buf.WriteString(" msg: ")
	buf.WriteString(s)
	buf.WriteString("\n")

	return buf
}
