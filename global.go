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

func globalLogFormatter(module Module, level LoggingLevel, callerLevel int, pool *BufferPool, format string, args ...interface{}) *bytes.Buffer {
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
		buf.WriteString(module.String())
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

func Trace(args ...interface{}) {
	logger.Trace(args...)
}

func Tracef(format string, args ...interface{}) {
	logger.Tracef(format, args...)
}

func MTrace(module Module, args ...interface{}) {
	logger.MTrace(module, args)
}

func MTracef(module Module, format string, args ...interface{}) {
	logger.MTracef(module, format, args...)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func MDebug(module Module, args ...interface{}) {
	logger.MDebug(module, args)
}

func MDebugf(moudle Module, format string, args ...interface{}) {
	logger.MDebugf(moudle, format, args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func MInfo(module Module, args ...interface{}) {
	logger.MInfo(module, args...)
}

func MInfof(module Module, format string, args ...interface{}) {
	logger.MInfof(module, format, args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func MWarn(module Module, args ...interface{}) {
	logger.MWarn(module, args...)
}

func MWarnf(module Module, format string, args ...interface{}) {
	logger.MWarnf(module, format, args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func MError(module Module, args ...interface{}) {
	logger.MError(module, args...)
}

func MErrorf(module Module, format string, args ...interface{}) {
	logger.MErrorf(module, format, args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

func MFatal(module Module, args ...interface{}) {
	logger.MFatal(module, args...)
}

func MFatalf(module Module, format string, args ...interface{}) {
	logger.MFatalf(module, format, args...)
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
