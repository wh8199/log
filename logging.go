package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
)

type LoggingInterface interface {
	Trace(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	MTrace(module Module, args ...interface{})
	MDebug(module Module, args ...interface{})
	MInfo(module Module, args ...interface{})
	MWarn(module Module, args ...interface{})
	MError(module Module, args ...interface{})
	MFatal(module Module, args ...interface{})

	MTracef(module Module, format string, args ...interface{})
	MDebugf(module Module, format string, args ...interface{})
	MInfof(module Module, format string, args ...interface{})
	MWarnf(module Module, format string, args ...interface{})
	MErrorf(module Module, format string, args ...interface{})
	MFatalf(module Module, format string, args ...interface{})

	Close()

	SetOutPut(w io.Writer)
	SetLevel(level LoggingLevel)
}

type Module string

func (m Module) String() string {
	return string(m)
}

type LoggingLevel int

const (
	TRACE_LEVEL LoggingLevel = iota
	DEBUG_LEVEL
	INFO_LEVEL
	WARN_LEVEL
	ERROR_LEVEL
	FATAL_LEVEL
)

func (level LoggingLevel) String() string {
	switch level {
	case TRACE_LEVEL:
		return "Trace"
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

	logging := &logging{
		Level:        level,
		Out:          os.Stdout,
		Pool:         NewBufferPool(),
		EnableCaller: true,
		CallerLevel:  callerLevel,
		Formater:     DefaultFormater,
	}

	//check globalConfig file is or not nil
	if globalConfig.file != nil {
		logging.SetOutPut(globalConfig.file)
	}

	globalConfig.Attach(logging)

	return logging
}

func NewLoggingWithFormater(level LoggingLevel, callerLevel int, formater Formatter) LoggingInterface {
	if level < DEBUG_LEVEL || level > FATAL_LEVEL {
		level = INFO_LEVEL
	}

	logging := &logging{
		Level:        level,
		Out:          os.Stdout,
		Pool:         NewBufferPool(),
		EnableCaller: true,
		CallerLevel:  callerLevel,
		Formater:     formater,
	}
	//Check golbal file is or not  empty, not empty Change the output object at init logging
	if globalConfig.file != nil {
		logging.SetOutPut(globalConfig.file)
	}

	globalConfig.Attach(logging)

	return logging
}

type logging struct {
	mux sync.Mutex

	Name         string
	Level        LoggingLevel
	Out          io.Writer
	Pool         *BufferPool
	EnableCaller bool
	CallerLevel  int
	Formater     func(Module, LoggingLevel, int, *BufferPool, string, ...interface{}) *bytes.Buffer
}

func (l *logging) Close() {
	l.mux.Lock()
	defer l.mux.Unlock()
	if closer, ok := l.Out.(io.WriteCloser); ok {
		closer.Close()
	}
}

func (l *logging) SetOutPut(w io.Writer) {
	l.mux.Lock()
	defer l.mux.Unlock()
	l.Out = w
}

func (l *logging) Write(buf *bytes.Buffer) {
	l.mux.Lock()
	defer l.mux.Unlock()
	if _, err := l.Out.Write(buf.Bytes()); err != nil {
		panic(err)
	}

	l.Pool.Put(buf)
}

func (l *logging) SetLevel(level LoggingLevel) {
	l.Level = level
}

func (l *logging) Trace(args ...interface{}) {
	if l.Level <= TRACE_LEVEL {
		s := fmt.Sprint(args...)
		l.Write(l.Formater("", TRACE_LEVEL, l.CallerLevel, l.Pool, s))
	}
}

func (l *logging) Debug(args ...interface{}) {
	if l.Level <= DEBUG_LEVEL {
		s := fmt.Sprint(args...)
		l.Write(l.Formater("", DEBUG_LEVEL, l.CallerLevel, l.Pool, s))
	}
}

func (l *logging) Info(args ...interface{}) {
	if l.Level <= INFO_LEVEL {
		s := fmt.Sprint(args...)
		l.Write(l.Formater("", INFO_LEVEL, l.CallerLevel, l.Pool, s))
	}
}

func (l *logging) Warn(args ...interface{}) {
	if l.Level <= WARN_LEVEL {
		s := fmt.Sprint(args...)
		l.Write(l.Formater("", WARN_LEVEL, l.CallerLevel, l.Pool, s))
	}
}

func (l *logging) Error(args ...interface{}) {
	if l.Level <= ERROR_LEVEL {
		s := fmt.Sprint(args...)
		l.Write(l.Formater("", ERROR_LEVEL, l.CallerLevel, l.Pool, s))
	}
}

func (l *logging) Fatal(args ...interface{}) {
	if l.Level <= FATAL_LEVEL {
		s := fmt.Sprint(args...)
		l.Write(l.Formater("", FATAL_LEVEL, l.CallerLevel, l.Pool, s))
	}
	os.Exit(1)
}

func (l *logging) Tracef(format string, args ...interface{}) {
	if l.Level <= TRACE_LEVEL {
		l.Write(l.Formater("", TRACE_LEVEL, l.CallerLevel, l.Pool, format, args...))
	}
}

func (l *logging) Debugf(format string, args ...interface{}) {
	if l.Level <= DEBUG_LEVEL {
		l.Write(l.Formater("", DEBUG_LEVEL, l.CallerLevel, l.Pool, format, args...))
	}
}

func (l *logging) Infof(format string, args ...interface{}) {
	if l.Level <= INFO_LEVEL {
		l.Write(l.Formater("", INFO_LEVEL, l.CallerLevel, l.Pool, format, args...))
	}
}

func (l *logging) Warnf(format string, args ...interface{}) {
	if l.Level <= WARN_LEVEL {
		l.Write(l.Formater("", WARN_LEVEL, l.CallerLevel, l.Pool, format, args...))
	}
}

func (l *logging) Errorf(format string, args ...interface{}) {
	if l.Level <= ERROR_LEVEL {
		l.Write(l.Formater("", ERROR_LEVEL, l.CallerLevel, l.Pool, format, args...))
	}
}

func (l *logging) Fatalf(format string, args ...interface{}) {
	if l.Level <= FATAL_LEVEL {
		l.Write(l.Formater("", FATAL_LEVEL, l.CallerLevel, l.Pool, format, args...))
	}
	os.Exit(1)
}

func (l *logging) MTrace(module Module, args ...interface{}) {
	if l.Level <= TRACE_LEVEL {
		s := fmt.Sprint(args...)
		l.Write(l.Formater(module, TRACE_LEVEL, l.CallerLevel, l.Pool, s))
	}
}

func (l *logging) MDebug(module Module, args ...interface{}) {
	if l.Level <= DEBUG_LEVEL {
		s := fmt.Sprint(args...)
		l.Write(l.Formater(module, DEBUG_LEVEL, l.CallerLevel, l.Pool, s))
	}
}

func (l *logging) MInfo(module Module, args ...interface{}) {
	if l.Level <= INFO_LEVEL {
		s := fmt.Sprint(args...)
		l.Write(l.Formater(module, INFO_LEVEL, l.CallerLevel, l.Pool, s))
	}
}

func (l *logging) MWarn(module Module, args ...interface{}) {
	if l.Level <= WARN_LEVEL {
		s := fmt.Sprint(args...)
		l.Write(l.Formater(module, WARN_LEVEL, l.CallerLevel, l.Pool, s))
	}
}

func (l *logging) MError(module Module, args ...interface{}) {
	if l.Level <= ERROR_LEVEL {
		s := fmt.Sprint(args...)
		l.Write(l.Formater(module, ERROR_LEVEL, l.CallerLevel, l.Pool, s))
	}
}

func (l *logging) MFatal(module Module, args ...interface{}) {
	if l.Level <= FATAL_LEVEL {
		s := fmt.Sprint(args...)
		l.Write(l.Formater(module, FATAL_LEVEL, l.CallerLevel, l.Pool, s))
	}
	os.Exit(1)
}

func (l *logging) MTracef(module Module, format string, args ...interface{}) {
	if l.Level <= TRACE_LEVEL {
		l.Write(l.Formater(module, TRACE_LEVEL, l.CallerLevel, l.Pool, format, args...))
	}
}

func (l *logging) MDebugf(module Module, format string, args ...interface{}) {
	if l.Level <= DEBUG_LEVEL {
		l.Write(l.Formater(module, DEBUG_LEVEL, l.CallerLevel, l.Pool, format, args...))
	}
}

func (l *logging) MInfof(module Module, format string, args ...interface{}) {
	if l.Level <= INFO_LEVEL {
		l.Write(l.Formater(module, INFO_LEVEL, l.CallerLevel, l.Pool, format, args...))
	}
}

func (l *logging) MWarnf(module Module, format string, args ...interface{}) {
	if l.Level <= WARN_LEVEL {
		l.Write(l.Formater(module, WARN_LEVEL, l.CallerLevel, l.Pool, format, args...))
	}
}

func (l *logging) MErrorf(module Module, format string, args ...interface{}) {
	if l.Level <= ERROR_LEVEL {
		l.Write(l.Formater(module, ERROR_LEVEL, l.CallerLevel, l.Pool, format, args...))
	}
}

func (l *logging) MFatalf(module Module, format string, args ...interface{}) {
	if l.Level <= FATAL_LEVEL {
		l.Write(l.Formater(module, FATAL_LEVEL, l.CallerLevel, l.Pool, format, args...))
	}
	os.Exit(1)
}
