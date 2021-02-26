package log

import (
	"bytes"
	"fmt"
	"runtime"
)

var (
	logger LoggingInterface
)

func init() {
	logger = NewLoggingWithFormater("global", INFO_LEVEL, 2, globalLogFormatter)
}

func globalLogFormatter(name string, level LoggingLevel, callerLevel int, pool *BufferPool, format string, args ...interface{}) *bytes.Buffer {
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

	result := fmt.Sprintf("%s %s:%d %v msg: %s\n", time, caller, line, level, s)

	buf := pool.Get()
	buf.WriteString(result)

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
	logger.Info(args)
}

func Printf(format string, args ...interface{}) {
	logger.Infof(format, args)
}
