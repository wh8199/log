package log

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestFileLoggingUse(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	SetLogConfig(true, "", "", "1s", "1s", 0, 0)
	globalConfig.Start()
	for i := 0; i < 9000; i++ {
		go func(i int) {
			logfile := NewLogging("test"+strconv.Itoa(i), INFO_LEVEL, 2)
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()
			for {
				<-ticker.C
				logfile.Info(" ---------github O-ll-O test write log----------")
			}
		}(i)
	}

	wg.Wait()
}
