package log

import (
	"sync"
	"time"
)

func init() {

	globalConfig = &logConfig{
		isFile:         false,
		prefix:         "log",
		maxSize:        20,
		maxSecond:      5,
		splitDuration:  time.Second * 5,
		deleteDuration: time.Second * 5,
		observers:      make([]LoggingInterface, 0, 2048),
		observersMu:    sync.Mutex{},
		fileMu:         sync.RWMutex{},
	}

	logger = NewLoggingWithFormater("global", INFO_LEVEL, 3, globalLogFormatter)

}
