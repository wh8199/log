package log

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestNewFileLogging(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	SetLogConfig(true, "", "", 0, 0)

	for i := 0; i < 1000; i++ {
		go func(i int) {
			logfile := NewLogging("test"+strconv.Itoa(i), INFO_LEVEL, 2)
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()
			for {
				<-ticker.C
				logfile.Info(" ---------test write log----------")

			}
		}(i)
	}
	globalConfig.Start()
	wg.Wait()
}
