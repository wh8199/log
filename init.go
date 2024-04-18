package log

func init() {
	logger = NewLoggingWithFormater(INFO_LEVEL, 5, globalLogFormatter)
}
