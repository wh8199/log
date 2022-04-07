package log

import "sync"

func init() {
	globalConfig = &logConfig{
		isFile:      false,
		observers:   make([]LoggingInterface, 0, 2048),
		observersMu: sync.Mutex{},
		fileMu:      sync.RWMutex{},
	}
}
func init() {
	logger = NewLoggingWithFormater("global", INFO_LEVEL, 3, globalLogFormatter)
}
