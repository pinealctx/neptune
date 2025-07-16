package ulog

import (
	"errors"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestNewZapLoggerX1(_ *testing.T) {
	test1()
}

func test1() {
	defer func() {
		var logger = NewSimpleLogger("info", zap.AddCaller(), zap.AddCallerSkip(1))
		logger.Debug("debug",
			zap.Duration("d", time.Second),
			zap.Int("i", 1),
			zap.Time("ti", time.Now()),
			zap.ByteString("bs", []byte("abc我")),
		)
		logger.Info("info",
			zap.Duration("d", time.Millisecond),
			zap.Int("i", 1),
			zap.Time("ti", time.Now()),
			zap.ByteString("bs", []byte("abc我")),
			zap.Stack("stack"),
		)
		logger.Warn("warn",
			zap.Duration("d", time.Second),
			zap.Int("i", 1),
			zap.Time("ti", time.Now()),
			zap.ByteString("bs", []byte("abc我")),
			zap.Error(nil),
		)
		logger.Error("err",
			zap.Duration("d", time.Second),
			zap.Int("i", 1),
			zap.Time("ti", time.Now()),
			zap.ByteString("bs", []byte("abc我")),
			zap.Error(errors.New("error1")),
		)
	}()
}

func TestNewZapLoggerX2(_ *testing.T) {
	test2()
}

func test2() {
	var logger = NewSimpleLogger("info")
	logger.Info("x", zap.Error(errors.New("error1")))
	logger.Info("x", zap.Strings("error", nil))
	var s = make([]string, 0)
	logger.Info("x", zap.Strings("error", s))

	s = []string{"1", "2"}
	logger.Info("x", zap.Strings("error", s))
	logger.Warn("x", zap.Strings("error", s))
	logger.Error("x", zap.Strings("error", s))
}
