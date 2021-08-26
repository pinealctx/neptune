package ulog

import (
	"testing"
)

func TestSetLogLevelStr(t *testing.T) {
	Debug("debug - debug level")
	Info("info - debug level")
	Warn("warn - debug level")
	Error("error - debug level")

	SetLogLevel(InfoLevel)

	Debug("debug - info level")
	Info("info - info level")
	Warn("warn - info level")
	Error("error - info level")

	SetLogLevel(WarnLevel)

	Debug("debug - warn level")
	Info("info - warn level")
	Warn("warn - warn level")
	Error("error - warn level")

	SetLogLevel(ErrorLevel)

	Debug("debug - error level")
	Info("info - error level")
	Warn("warn - error level")
	Error("error - error level")
}
