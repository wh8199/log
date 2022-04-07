package log

import (
	"sync"
	"testing"
	"time"
)

func TestNewFileLogging1(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	SetLogConfig(true, "", "", 0, 0)
	globalConfig.Start()
	logfile := NewLogging("test", INFO_LEVEL, 2)
	go func() {
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				logfile.Info("test write log")
			}
		}
	}()
	wg.Wait()
}
