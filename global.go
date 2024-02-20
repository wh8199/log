package log

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"strconv"
)

var (
	logger *logging
)

func globalLogFormatter(logRecord *LogRecord) *bytes.Buffer {
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
		buf.WriteString("[")
		buf.WriteString(logRecord.module)
		buf.WriteString("] ")
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

func SetLogLevel(level LoggingLevel) {
	logger.SetLevel(level)
}

func Trace(args ...interface{}) {
	logger.Trace(args...)
}

func Tracef(format string, args ...interface{}) {
	logger.Tracef(format, args...)
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

func Module(module string) *LogRecord {
	return logger.Module(module)
}

func CallLevel(level int) *LogRecord {
	return logger.Caller(level)
}
