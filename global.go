package log

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"strconv"
)

var (
	logger LoggingInterface
)

func globalLogFormatter(level LoggingLevel, callerLevel int, pool *BufferPool, module, format string, args ...interface{}) *bytes.Buffer {
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
		buf.WriteString("[")
		buf.WriteString(module)
		buf.WriteString("] ")
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

func SetLogLevel(level LoggingLevel) {
	logger.SetLevel(level)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

func Println(args ...interface{}) {
	logger.Info(args...)
}

func Printf(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Fatalln(args ...interface{}) {
	logger.Fatal(args...)
}

func SetOutPut(w io.Writer) {
	logger.SetOutPut(w)
}
