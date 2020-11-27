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

	SetOutPut(w io.Writer)
}

type LoggingLevel int

const (
	DEBUG_LEVEL LoggingLevel = iota
	INFO_LEVEL
	WARN_LEVEL
	ERROR_LEVEL
	FATAL_LEVEL
)

func (level LoggingLevel) String() string {
	switch level {
	case DEBUG_LEVEL:
		return "Debug"
	case INFO_LEVEL:
		return "Info"
	case WARN_LEVEL:
		return "Warn"
	case ERROR_LEVEL:
		return "Error"
	case FATAL_LEVEL:
		return "Fatal"
	default:
		return "Unknown"
	}
}

func NewLogging(name string, level LoggingLevel, callerLevel int) LoggingInterface {
	if level < DEBUG_LEVEL || level > FATAL_LEVEL {
		level = INFO_LEVEL
	}

	return &logging{
		Name:         name,
		Level:        level,
		Out:          os.Stdout,
		Pool:         NewBufferPool(),
		EnableCaller: true,
		CallerLevel:  callerLevel,
		Formater:     DefaultFormater,
	}
}

func NewLoggingWithFormater(name string, level LoggingLevel, callerLevel int, formater Formatter) LoggingInterface {
	if level < DEBUG_LEVEL || level > FATAL_LEVEL {
		level = INFO_LEVEL
	}

	return &logging{
		Name:         name,
		Level:        level,
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
	Pool         *BufferPool
	EnableCaller bool
	CallerLevel  int
	Formater     func(string, LoggingLevel, int, *BufferPool, string, ...interface{}) *bytes.Buffer
}

func (l *logging) Close() {
	if closer, ok := l.Out.(io.WriteCloser); ok {
		closer.Close()
	}
}

func (l *logging) SetOutPut(w io.Writer) {
	l.Out = w
}

func (l *logging) Write(buf *bytes.Buffer) {
	l.Out.Write(buf.Bytes())
	l.Pool.Put(buf)
}

func (l *logging) Debug(args ...interface{}) {
	if l.Level <= DEBUG_LEVEL {
		s := fmt.Sprint(args...)
		l.Write(l.Formater(l.Name, DEBUG_LEVEL, l.CallerLevel, l.Pool, s))
	}
}

func (l *logging) Info(args ...interface{}) {
	if l.Level <= INFO_LEVEL {
		s := fmt.Sprint(args...)
		l.Write(l.Formater(l.Name, INFO_LEVEL, l.CallerLevel, l.Pool, s))
	}
}

func (l *logging) Warn(args ...interface{}) {
	if l.Level <= WARN_LEVEL {
		s := fmt.Sprint(args...)
		l.Write(l.Formater(l.Name, WARN_LEVEL, l.CallerLevel, l.Pool, s))
	}
}

func (l *logging) Error(args ...interface{}) {
	if l.Level <= ERROR_LEVEL {
		s := fmt.Sprint(args...)
		l.Write(l.Formater(l.Name, ERROR_LEVEL, l.CallerLevel, l.Pool, s))
	}
}

func (l *logging) Fatal(args ...interface{}) {
	if l.Level <= FATAL_LEVEL {
		s := fmt.Sprint(args...)
		l.Write(l.Formater(l.Name, FATAL_LEVEL, l.CallerLevel, l.Pool, s))
	}
	os.Exit(1)
}

func (l *logging) Debugf(format string, args ...interface{}) {
	if l.Level <= DEBUG_LEVEL {
		l.Write(l.Formater(l.Name, DEBUG_LEVEL, l.CallerLevel, l.Pool, format, args...))
	}
}

func (l *logging) Infof(format string, args ...interface{}) {
	if l.Level <= INFO_LEVEL {
		l.Write(l.Formater(l.Name, INFO_LEVEL, l.CallerLevel, l.Pool, format, args...))
	}
}

func (l *logging) Warnf(format string, args ...interface{}) {
	if l.Level <= WARN_LEVEL {
		l.Write(l.Formater(l.Name, WARN_LEVEL, l.CallerLevel, l.Pool, format, args...))
	}
}

func (l *logging) Errorf(format string, args ...interface{}) {
	if l.Level <= ERROR_LEVEL {
		l.Write(l.Formater(l.Name, ERROR_LEVEL, l.CallerLevel, l.Pool, format, args...))
	}
}

func (l *logging) Fatalf(format string, args ...interface{}) {
	if l.Level <= FATAL_LEVEL {
		l.Write(l.Formater(l.Name, FATAL_LEVEL, l.CallerLevel, l.Pool, format, args...))
	}
	os.Exit(1)
}
