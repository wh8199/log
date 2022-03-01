package log

import (
	"sync"
	"testing"
	"time"
)

func TestNewFileLogging(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	logfile := NewFileLogging("test", INFO_LEVEL, 2, true, "", "", 0, 0)
	logfile.Start()
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
