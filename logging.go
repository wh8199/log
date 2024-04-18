package log

import (
	"bytes"
	"io"
	"os"
	"sync"
	"time"
)

type LogRotateConfig struct {
	EnableLogFile bool   `json:"enableLogFile"`
	Prefix        string `json:"prefix"`
	FileDir       string `json:"fileDir"`
	MaxSize       int64  `json:"maxSize"`
	// unit is second
	MaxLogLife int64 `json:"maxLogLife"`
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
		level:        level,
		output:       os.Stdout,
		pool:         NewBufferPool(),
		enableCaller: true,
		callerLevel:  callerLevel,
		Formater:     DefaultFormater,
		exitChan:     make(chan struct{}, 2),
	}

	return logging
}

func NewLoggingWithFormater(level LoggingLevel, callerLevel int, formater Formatter) *logging {
	if level < DEBUG_LEVEL || level > FATAL_LEVEL {
		level = INFO_LEVEL
	}

	logging := &logging{
		level:        level,
		output:       os.Stdout,
		pool:         NewBufferPool(),
		enableCaller: true,
		callerLevel:  callerLevel,
		Formater:     formater,
		exitChan:     make(chan struct{}, 2),
	}

	return logging
}

type logging struct {
	mux  sync.Mutex
	name string
	// default log level
	level        LoggingLevel
	output       io.Writer
	pool         *BufferPool
	enableCaller bool
	// default caller level
	callerLevel int
	Formater    func(logRecord *LogRecord) *bytes.Buffer

	LogRotateConfig
	exitChan chan struct{}

	isStarted bool
}

func (l *logging) Close() {
	l.mux.Lock()
	defer l.mux.Unlock()

	if l.exitChan != nil {
		l.exitChan <- struct{}{}
	}
}

func (l *logging) UpdateConfig(cfg LogRotateConfig) {
	l.mux.Lock()
	defer l.mux.Unlock()

	l.LogRotateConfig = cfg
}

func (l *logging) SetOutPut(w io.Writer) {
	l.mux.Lock()
	defer l.mux.Unlock()
	l.output = w
}

func (l *logging) Write(buf *bytes.Buffer) {
	l.mux.Lock()
	defer l.mux.Unlock()

	if _, err := l.output.Write(buf.Bytes()); err != nil {
		panic(err)
	}

	l.pool.Put(buf)
}

func (l *logging) SetLevel(level LoggingLevel) {
	l.mux.Lock()
	defer l.mux.Unlock()

	l.level = level
}

func (l *logging) print(level LoggingLevel, args ...interface{}) {
	if l.level <= level {
		record := &LogRecord{
			logLevel:     level,
			callerLevel:  l.callerLevel,
			enableCaller: l.enableCaller,
			logger:       l,
		}

		record.print(args...)
	}
}

func (l *logging) printf(level LoggingLevel, format string, args ...interface{}) {
	if l.level <= level {
		record := &LogRecord{
			logLevel:     level,
			callerLevel:  l.callerLevel,
			enableCaller: l.enableCaller,
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
		logLevel:     l.level,
		callerLevel:  l.callerLevel,
		enableCaller: l.enableCaller,
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
		callerLevel:  l.callerLevel,
		enableCaller: l.enableCaller,
		logger:       l,
	}
}

func (l *logging) Start() {
	if !l.EnableLogFile {
		return
	}

	l.mux.Lock()
	if l.isStarted {
		l.mux.Unlock()
		return
	}
	l.mux.Unlock()

	output, err := NewFileOutput(l.LogRotateConfig)
	if err != nil {
		panic(err)
	}

	l.SetOutPut(output)

	go func() {
		ticker := time.NewTicker(time.Second * 1)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				l.mux.Lock()
				fileOutput := l.output.(OutPut)
				if err := fileOutput.Rotate(); err != nil {
					panic(err)
				}

				l.mux.Unlock()
			case <-l.exitChan:
				if closer, ok := l.output.(io.WriteCloser); ok {
					closer.Close()
				}
				return
			}
		}
	}()

	l.mux.Lock()
	l.isStarted = true
	l.mux.Unlock()
}
