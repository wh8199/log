package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

type LoggingInterface interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Close()
}

type LoggingLevel int

const (
	DEBUG_LEVEL LoggingLevel = iota
	INFO_LEVEL
	WARN_LEVEL
	ERROR_LEVEL
	FATAL_LEVEL
)

var levelMap map[LoggingLevel]string

func init() {
	levelMap = map[LoggingLevel]string{}

	levelMap[DEBUG_LEVEL] = "Debug"
	levelMap[INFO_LEVEL] = "Info"
	levelMap[WARN_LEVEL] = "Warn"
	levelMap[ERROR_LEVEL] = "Error"
	levelMap[FATAL_LEVEL] = "Fatal"
}

func NewLogging(name, outpath string, level LoggingLevel, callerLevel int) LoggingInterface {
	return &logging{
		Name:         name,
		Level:        level,
		OutPath:      outpath,
		Out:          os.Stdout,
		Pool:         NewBufferPool(),
		EnableCaller: true,
		CallerLevel:  callerLevel,
		Formater:     DefaultFormater,
	}
}

func NewLoggingWithFormater(name, outpath string, level LoggingLevel, callerLevel int, formater Formatter) LoggingInterface {
	return &logging{
		Name:         name,
		Level:        level,
		OutPath:      outpath,
		Out:          os.Stdout,
		Pool:         NewBufferPool(),
		EnableCaller: true,
		CallerLevel:  callerLevel,
		Formater:     formater,
	}
}

type logging struct {
	Name         string
	Level        LoggingLevel
	Out          io.Writer
	OutPath      string
	Pool         *BufferPool
	EnableCaller bool
	CallerLevel  int
	Formater     func(string, string, int, *BufferPool, string, ...interface{}) *bytes.Buffer
}

func (l *logging) Close() {
	if closer, ok := l.Out.(io.WriteCloser); ok {
		closer.Close()
	}
}

func (l *logging) Write(buf *bytes.Buffer) {
	l.Out.Write(buf.Bytes())
	l.Pool.Put(buf)
}

func (l *logging) Debug(args ...interface{}) {
	if l.Level <= DEBUG_LEVEL {
		s := fmt.Sprintln(args)
		l.Write(l.Formater(l.Name, levelMap[DEBUG_LEVEL], l.CallerLevel, l.Pool, s))
	}
}

func (l *logging) Info(args ...interface{}) {
	if l.Level <= INFO_LEVEL {
		s := fmt.Sprint(args...)
		l.Write(l.Formater(l.Name, levelMap[INFO_LEVEL], l.CallerLevel, l.Pool, s))
	}
}

func (l *logging) Warn(args ...interface{}) {
	if l.Level <= WARN_LEVEL {
		s := fmt.Sprintln(args)
		l.Write(l.Formater(l.Name, levelMap[WARN_LEVEL], l.CallerLevel, l.Pool, s))
	}
}

func (l *logging) Error(args ...interface{}) {
	if l.Level <= ERROR_LEVEL {
		s := fmt.Sprintln(args)
		l.Write(l.Formater(l.Name, levelMap[ERROR_LEVEL], l.CallerLevel, l.Pool, s))
	}
}

func (l *logging) Fatal(args ...interface{}) {
	if l.Level <= FATAL_LEVEL {
		s := fmt.Sprintln(args)
		l.Write(l.Formater(l.Name, levelMap[FATAL_LEVEL], l.CallerLevel, l.Pool, s))
	}
	os.Exit(1)
}

func (l *logging) Debugf(format string, args ...interface{}) {
	if l.Level <= DEBUG_LEVEL {
		l.Write(l.Formater(l.Name, levelMap[DEBUG_LEVEL], l.CallerLevel, l.Pool, format, args...))
	}
}

func (l *logging) Infof(format string, args ...interface{}) {
	if l.Level <= INFO_LEVEL {
		l.Write(l.Formater(l.Name, levelMap[INFO_LEVEL], l.CallerLevel, l.Pool, format, args...))
	}
}

func (l *logging) Warnf(format string, args ...interface{}) {
	if l.Level <= WARN_LEVEL {
		l.Write(l.Formater(l.Name, levelMap[WARN_LEVEL], l.CallerLevel, l.Pool, format, args...))
	}
}

func (l *logging) Errorf(format string, args ...interface{}) {
	if l.Level <= ERROR_LEVEL {
		l.Write(l.Formater(l.Name, levelMap[ERROR_LEVEL], l.CallerLevel, l.Pool, format, args...))
	}
}

func (l *logging) Fatalf(format string, args ...interface{}) {
	if l.Level <= FATAL_LEVEL {
		l.Write(l.Formater(l.Name, levelMap[FATAL_LEVEL], l.CallerLevel, l.Pool, format, args...))
	}
	os.Exit(1)
}
