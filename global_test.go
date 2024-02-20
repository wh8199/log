package log_test

import (
	"testing"

	"github.com/wh8199/log"
)

func TestGlobalPrint(t *testing.T) {
	log.SetLogLevel(log.DEBUG_LEVEL)
	log.Info("info")
	log.Debug("debug")
	log.Warn("warn")
	log.Error("error")
}

func TestGlobalModulePrint(t *testing.T) {
	log.SetLogLevel(log.DEBUG_LEVEL)
	log.MInfo("video", "info")
	log.MDebug("video", "debug")
	log.MWarn("video", "warn")
	log.MError("video", "error")
}

func TestGlobalModuleFormatPrint(t *testing.T) {
	log.SetLogLevel(log.TRACE_LEVEL)
	log.MTracef("video", "test %s", "trace")
	log.MInfof("video", "test %s", "info")
	log.MDebugf("video", "test %s", "debug")
	log.MWarnf("video", "test %s", "warn")
	log.MErrorf("video", "test %s", "error")
}
