package ulog

import (
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	//default logger
	_DefaultLogger atomic.Value
)

// Debug : log debug
func Debug(msg string, fields ...zap.Field) {
	loadLogger().Debug(msg, fields...)
}

// Info : log info
func Info(msg string, fields ...zap.Field) {
	loadLogger().Info(msg, fields...)
}

// Warn : log warn
func Warn(msg string, fields ...zap.Field) {
	loadLogger().Warn(msg, fields...)
}

// Error : log error
func Error(msg string, fields ...zap.Field) {
	loadLogger().Error(msg, fields...)
}

// GetLevel : get log level
func GetLevel() zapcore.Level {
	return loadLogger().Level()
}

// GetLevelStr : get log level str
func GetLevelStr() string {
	return loadLogger().LevelStr()
}

// SetLogLevel : set log level
func SetLogLevel(level zapcore.Level) {
	loadLogger().SetLevel(level)
}

// SetLogLevelStr : set log level by string
func SetLogLevelStr(levelStr string) {
	loadLogger().SetLevelStr(levelStr)
}

// SetDefaultLogger : set default logger
func SetDefaultLogger(logger *Logger) {
	if logger == nil {
		return
	}
	if logger.Logger == nil {
		return
	}
	var cloneLogger = *logger
	//add caller skip 1 because actually there is a wrapped function.
	cloneLogger.Logger = logger.WithOptions(zap.AddCallerSkip(1))
	_DefaultLogger.Store(&cloneLogger)
}

// GetDefaultLogger : get default logger
func GetDefaultLogger() *Logger {
	return loadLogger()
}

// load logger
func loadLogger() *Logger {
	var logger = _DefaultLogger.Load()
	if logger == nil {
		return nil
	}
	var zapLogger, ok = logger.(*Logger)
	if ok {
		return zapLogger
	}
	return nil
}

func init() {
	var logger = NewSimpleLogger(DebugLevelStr, zap.AddCaller(), zap.AddCallerSkip(2))
	_DefaultLogger.Store(logger)
}
