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
	BuildLogConfig(true, "", "", "1s", "1s", 0, 0)
	for i := 0; i < 10000; i++ {
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
	Notify()
	Start()
	wg.Wait()
}

func TestAttachAndDatch(t *testing.T) {
	num := len(globalConfig.observers)
	//Automatically attach new log objects
	logger = NewLoggingWithFormater(INFO_LEVEL, 3, globalLogFormatter)
	if len(globalConfig.observers) != num+1 {
		t.Error("attacth fail")
	}
	globalConfig.Detach(logger)
	if len(globalConfig.observers) != num {
		t.Error("detach fail")
	}
}
