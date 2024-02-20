package log

import "testing"

func TestGlobalLog(t *testing.T) {
	Info("aaaaa")
	Infof("%d", 1)
}
