package log_test

import (
	"testing"

	"github.com/wh8199/log"
)

func TestGloablPrint(t *testing.T) {
	log.SetLogLevel(log.DEBUG_LEVEL)
	log.Info("info")
	log.Debug("debug")
	log.Warn("warn")
	log.Error("error")
}
