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

	MTrace(module string, args ...interface{})
	MDebug(module string, args ...interface{})
	MInfo(module string, args ...interface{})
	MWarn(module string, args ...interface{})
	MError(module string, args ...interface{})
	MFatal(module string, args ...interface{})

	MTracef(module, format string, args ...interface{})
	MDebugf(module, format string, args ...interface{})
	MInfof(module, format string, args ...interface{})
	MWarnf(module, format string, args ...interface{})
	MErrorf(module, format string, args ...interface{})
	MFatalf(module, format string, args ...interface{})

	Close()

	SetOutPut(w io.Writer)
	SetLevel(level LoggingLevel)
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
	Formater     func(LoggingLevel, int, *BufferPool, string, string, ...interface{}) *bytes.Buffer
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
		l.Write(l.Formater(TRACE_LEVEL, l.CallerLevel, l.Pool, "", s))
	}
}

func (l *logging) Debug(args ...interface{}) {
	if l.Level <= DEBUG_LEVEL {
		s := fmt.Sprint(args...)
		l.Write(l.Formater(DEBUG_LEVEL, l.CallerLevel, l.Pool, "", s))
	}
}

func (l *logging) Info(args ...interface{}) {
	if l.Level <= INFO_LEVEL {
		s := fmt.Sprint(args...)
		l.Write(l.Formater(INFO_LEVEL, l.CallerLevel, l.Pool, "", s))
	}
}

func (l *logging) Warn(args ...interface{}) {
	if l.Level <= WARN_LEVEL {
		s := fmt.Sprint(args...)
		l.Write(l.Formater(WARN_LEVEL, l.CallerLevel, l.Pool, "", s))
	}
}

func (l *logging) Error(args ...interface{}) {
	if l.Level <= ERROR_LEVEL {
		s := fmt.Sprint(args...)
		l.Write(l.Formater(ERROR_LEVEL, l.CallerLevel, l.Pool, "", s))
	}
}

func (l *logging) Fatal(args ...interface{}) {
	if l.Level <= FATAL_LEVEL {
		s := fmt.Sprint(args...)
		l.Write(l.Formater(FATAL_LEVEL, l.CallerLevel, l.Pool, "", s))
	}
	os.Exit(1)
}

func (l *logging) Tracef(format string, args ...interface{}) {
	if l.Level <= TRACE_LEVEL {
		l.Write(l.Formater(TRACE_LEVEL, l.CallerLevel, l.Pool, format, "", args...))
	}
}

func (l *logging) Debugf(format string, args ...interface{}) {
	if l.Level <= DEBUG_LEVEL {
		l.Write(l.Formater(DEBUG_LEVEL, l.CallerLevel, l.Pool, format, "", args...))
	}
}

func (l *logging) Infof(format string, args ...interface{}) {
	if l.Level <= INFO_LEVEL {
		l.Write(l.Formater(INFO_LEVEL, l.CallerLevel, l.Pool, format, "", args...))
	}
}

func (l *logging) Warnf(format string, args ...interface{}) {
	if l.Level <= WARN_LEVEL {
		l.Write(l.Formater(WARN_LEVEL, l.CallerLevel, l.Pool, format, "", args...))
	}
}

func (l *logging) Errorf(format string, args ...interface{}) {
	if l.Level <= ERROR_LEVEL {
		l.Write(l.Formater(ERROR_LEVEL, l.CallerLevel, l.Pool, format, "", args...))
	}
}

func (l *logging) Fatalf(format string, args ...interface{}) {
	if l.Level <= FATAL_LEVEL {
		l.Write(l.Formater(FATAL_LEVEL, l.CallerLevel, l.Pool, format, "", args...))
	}
	os.Exit(1)
}

func (l *logging) MTrace(module string, args ...interface{}) {
	if l.Level <= TRACE_LEVEL {
		s := fmt.Sprint(args...)
		l.Write(l.Formater(TRACE_LEVEL, l.CallerLevel, l.Pool, module, s))
	}
}

func (l *logging) MDebug(module string, args ...interface{}) {
	if l.Level <= DEBUG_LEVEL {
		s := fmt.Sprint(args...)
		l.Write(l.Formater(DEBUG_LEVEL, l.CallerLevel, l.Pool, module, s))
	}
}

func (l *logging) MInfo(module string, args ...interface{}) {
	if l.Level <= INFO_LEVEL {
		s := fmt.Sprint(args...)
		l.Write(l.Formater(INFO_LEVEL, l.CallerLevel, l.Pool, module, s))
	}
}

func (l *logging) MWarn(module string, args ...interface{}) {
	if l.Level <= WARN_LEVEL {
		s := fmt.Sprint(args...)
		l.Write(l.Formater(WARN_LEVEL, l.CallerLevel, l.Pool, module, s))
	}
}

func (l *logging) MError(module string, args ...interface{}) {
	if l.Level <= ERROR_LEVEL {
		s := fmt.Sprint(args...)
		l.Write(l.Formater(ERROR_LEVEL, l.CallerLevel, l.Pool, module, s))
	}
}

func (l *logging) MFatal(module string, args ...interface{}) {
	if l.Level <= FATAL_LEVEL {
		s := fmt.Sprint(args...)
		l.Write(l.Formater(FATAL_LEVEL, l.CallerLevel, l.Pool, module, s))
	}
	os.Exit(1)
}

func (l *logging) MTracef(module, format string, args ...interface{}) {
	if l.Level <= TRACE_LEVEL {
		l.Write(l.Formater(TRACE_LEVEL, l.CallerLevel, l.Pool, module, format, args...))
	}
}

func (l *logging) MDebugf(module, format string, args ...interface{}) {
	if l.Level <= DEBUG_LEVEL {
		l.Write(l.Formater(DEBUG_LEVEL, l.CallerLevel, l.Pool, module, format, args...))
	}
}

func (l *logging) MInfof(module, format string, args ...interface{}) {
	if l.Level <= INFO_LEVEL {
		l.Write(l.Formater(INFO_LEVEL, l.CallerLevel, l.Pool, module, format, args...))
	}
}

func (l *logging) MWarnf(module, format string, args ...interface{}) {
	if l.Level <= WARN_LEVEL {
		l.Write(l.Formater(WARN_LEVEL, l.CallerLevel, l.Pool, module, format, args...))
	}
}

func (l *logging) MErrorf(module, format string, args ...interface{}) {
	if l.Level <= ERROR_LEVEL {
		l.Write(l.Formater(ERROR_LEVEL, l.CallerLevel, l.Pool, module, format, args...))
	}
}

func (l *logging) MFatalf(module, format string, args ...interface{}) {
	if l.Level <= FATAL_LEVEL {
		l.Write(l.Formater(FATAL_LEVEL, l.CallerLevel, l.Pool, module, format, args...))
	}
	os.Exit(1)
}
