package log

import (
	"bytes"
	"io"
	"os"
	"sync"
)

/* type LoggingInterface interface {
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

	MTrace(module Module1, args ...interface{})
	MDebug(module Module1, args ...interface{})
	MInfo(module Module1, args ...interface{})
	MWarn(module Module1, args ...interface{})
	MError(module Module1, args ...interface{})
	MFatal(module Module1, args ...interface{})

	MTracef(module Module1, format string, args ...interface{})
	MDebugf(module Module1, format string, args ...interface{})
	MInfof(module Module1, format string, args ...interface{})
	MWarnf(module Module1, format string, args ...interface{})
	MErrorf(module Module1, format string, args ...interface{})
	MFatalf(module Module1, format string, args ...interface{})

	Close()

	SetOutPut(w io.Writer)
	SetLevel(level LoggingLevel)
} */

type Module1 string

func (m Module1) String() string {
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

func NewLogging(name string, level LoggingLevel, callerLevel int) *logging {
	if level < TRACE_LEVEL || level > FATAL_LEVEL {
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

func NewLoggingWithFormater(level LoggingLevel, callerLevel int, formater Formatter) *logging {
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

	Name string
	// default log level
	Level        LoggingLevel
	Out          io.Writer
	Pool         *BufferPool
	EnableCaller bool
	// default caller level
	CallerLevel int
	Formater    func(logRecord *LogRecord) *bytes.Buffer
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
	l.mux.Lock()
	defer l.mux.Unlock()

	l.Level = level
}

func (l *logging) print(level LoggingLevel, args ...interface{}) {
	if l.Level <= level {
		record := &LogRecord{
			logLevel:     l.Level,
			callerLevel:  l.CallerLevel,
			enableCaller: l.EnableCaller,
			logger:       l,
		}

		record.print(args...)
	}
}

func (l *logging) printf(level LoggingLevel, format string, args ...interface{}) {
	if l.Level <= level {
		record := &LogRecord{
			logLevel:     l.Level,
			callerLevel:  l.CallerLevel,
			enableCaller: l.EnableCaller,
			logger:       l,
		}

		record.printf(format, args...)
	}
}

func (l *logging) Trace(args ...interface{}) {
	l.print(TRACE_LEVEL, args...)
}

func (l *logging) Debug(args ...interface{}) {
	l.print(DEBUG_LEVEL, args...)
}

func (l *logging) Info(args ...interface{}) {
	l.print(INFO_LEVEL, args...)
}

func (l *logging) Warn(args ...interface{}) {
	l.print(WARN_LEVEL, args...)
}

func (l *logging) Error(args ...interface{}) {
	l.print(ERROR_LEVEL, args...)
}

func (l *logging) Fatal(args ...interface{}) {
	l.print(ERROR_LEVEL, args...)
	os.Exit(1)
}

func (l *logging) Tracef(format string, args ...interface{}) {
	l.printf(TRACE_LEVEL, format, args...)
}

func (l *logging) Debugf(format string, args ...interface{}) {
	l.printf(DEBUG_LEVEL, format, args...)
}

func (l *logging) Infof(format string, args ...interface{}) {
	l.printf(INFO_LEVEL, format, args...)
}

func (l *logging) Warnf(format string, args ...interface{}) {
	l.printf(WARN_LEVEL, format, args...)
}

func (l *logging) Errorf(format string, args ...interface{}) {
	l.printf(ERROR_LEVEL, format, args...)
}

func (l *logging) Fatalf(format string, args ...interface{}) {
	l.printf(FATAL_LEVEL, format, args...)
	os.Exit(1)
}

func (l *logging) Module(moudle string) *LogRecord {
	return &LogRecord{
		module:       moudle,
		logLevel:     l.Level,
		callerLevel:  l.CallerLevel,
		enableCaller: l.EnableCaller,
		logger:       l,
	}
}

func (l *logging) Caller(level int) *LogRecord {
	return &LogRecord{
		callerLevel:  level,
		enableCaller: true,
		logger:       l,
	}
}

func (l *logging) LogLevel(logLevel LoggingLevel) *LogRecord {
	return &LogRecord{
		logLevel:     logLevel,
		callerLevel:  l.CallerLevel,
		enableCaller: l.EnableCaller,
		logger:       l,
	}
}
