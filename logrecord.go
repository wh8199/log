package log

import "os"

type LogRecord struct {
	format       string
	args         []interface{}
	module       string
	callerLevel  int
	enableCaller bool
	logLevel     LoggingLevel
	logger       *logging
}

func (l *LogRecord) print(args ...interface{}) {
	l.args = args
	l.logger.Write(l.logger.Formater(l))
}

func (l *LogRecord) printf(format string, args ...interface{}) {
	l.args = args
	l.format = format
	l.logger.Write(l.logger.Formater(l))
}

func (l *LogRecord) Trace(args ...interface{}) {
	if l.logLevel <= TRACE_LEVEL {
		l.args = args
		l.logger.Write(l.logger.Formater(l))
	}
}

func (l *LogRecord) Debug(args ...interface{}) {
	if l.logLevel <= DEBUG_LEVEL {
		l.args = args
		l.logger.Write(l.logger.Formater(l))
	}
}

func (l *LogRecord) Info(args ...interface{}) {
	if l.logLevel <= INFO_LEVEL {
		l.args = args
		l.logger.Write(l.logger.Formater(l))
	}
}

func (l *LogRecord) Warn(args ...interface{}) {
	if l.logLevel <= WARN_LEVEL {
		l.args = args
		l.logger.Write(l.logger.Formater(l))
	}
}

func (l *LogRecord) Error(args ...interface{}) {
	if l.logLevel <= ERROR_LEVEL {
		l.args = args
		l.logger.Write(l.logger.Formater(l))
	}
}

func (l *LogRecord) Fatal(args ...interface{}) {
	if l.logLevel <= FATAL_LEVEL {
		l.args = args
		l.logger.Write(l.logger.Formater(l))
		os.Exit(0)
	}
}

func (l *LogRecord) Tracef(format string, args ...interface{}) {
	if l.logLevel <= TRACE_LEVEL {
		l.args = args
		l.format = format
		l.logger.Write(l.logger.Formater(l))
	}
}

func (l *LogRecord) Debugf(format string, args ...interface{}) {
	if l.logLevel <= DEBUG_LEVEL {
		l.args = args
		l.format = format
		l.logger.Write(l.logger.Formater(l))
	}
}

func (l *LogRecord) Infof(format string, args ...interface{}) {
	if l.logLevel <= INFO_LEVEL {
		l.args = args
		l.format = format
		l.logger.Write(l.logger.Formater(l))
	}
}

func (l *LogRecord) Warnf(format string, args ...interface{}) {
	if l.logLevel <= WARN_LEVEL {
		l.args = args
		l.format = format
		l.logger.Write(l.logger.Formater(l))
	}
}

func (l *LogRecord) Errorf(format string, args ...interface{}) {
	if l.logLevel <= ERROR_LEVEL {
		l.args = args
		l.format = format
		l.logger.Write(l.logger.Formater(l))
	}
}

func (l *LogRecord) Fatalf(format string, args ...interface{}) {
	if l.logLevel <= FATAL_LEVEL {
		l.args = args
		l.format = format
		l.logger.Write(l.logger.Formater(l))
	}

	os.Exit(0)
}

func (l *LogRecord) Module(module string) *LogRecord {
	l.module = module
	return l
}

func (l *LogRecord) CallerLevel(callerLevel int) *LogRecord {
	l.callerLevel = callerLevel
	l.enableCaller = true
	return l
}
