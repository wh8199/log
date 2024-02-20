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
		observers:      make([]*logging, 0, 2048),
		observersMu:    sync.Mutex{},
		fileMu:         sync.RWMutex{},
	}

	logger = NewLoggingWithFormater(INFO_LEVEL, 5, globalLogFormatter)
}
