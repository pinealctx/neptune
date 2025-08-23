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
	fs1 := []zap.Field{
		zap.Duration("d", time.Second),
		zap.Int("i", 1),
		zap.Time("ti", time.Now()),
		zap.ByteString("bs", []byte("abc我")),
	}
	fs2 := []zap.Field{
		zap.Duration("d", time.Millisecond),
		zap.Int("i", 1),
		zap.Time("ti", time.Now()),
		zap.ByteString("bs", []byte("abc我")),
		zap.Stack("stack"),
	}
	fs3 := []zap.Field{
		zap.Duration("d", time.Microsecond),
		zap.Int("i", 1),
		zap.Time("ti", time.Now()),
		zap.ByteString("bs", []byte("abc我")),
	}
	fs4 := []zap.Field{
		zap.Duration("d", time.Nanosecond),
		zap.Int("i", 1),
		zap.Time("ti", time.Now()),
		zap.ByteString("bs", []byte("abc我")),
	}
	defer func() {
		var logger = NewSimpleLogger("info", zap.AddCaller(), zap.AddCallerSkip(1))
		logger.Debug("debug.defer", fs1...)
		logger.Info("info.defer", fs2...)
		logger.Warn("warn.defer", fs3...)
		logger.Error("err.defer", fs4...)
	}()
	var logger = NewSimpleLogger("info")
	logger.Debug("debug", fs1...)
	logger.Info("info", fs2...)
	logger.Warn("warn", fs3...)
	logger.Error("err", fs4...)
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
