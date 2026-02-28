package log

import "testing"

func TestLog(t *testing.T) {
	SetLogLevel(LevelInfo)
	Error("error")
	Info("info ?", 1)
}

func TestLog_Error(t *testing.T) {
	SetLogLevel(LevelError)
	Error("error")
	Info("info")
}

func TestLog_Disabled(t *testing.T) {
	SetLogLevel(Disabled)
	Error("error")
	Info("info")
}
