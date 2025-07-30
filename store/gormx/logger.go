package gormx

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"

	"github.com/pinealctx/neptune/ulog"
)

type uberLogger struct {
}

// NewUberLogger 创建一个新的 GORM 日志适配器
func NewUberLogger() logger.Interface {
	return &uberLogger{}
}

// LogMode 设置日志级别
func (u *uberLogger) LogMode(level logger.LogLevel) logger.Interface {
	// 将 GORM 的日志级别映射到 ulog 的日志级别
	var ulogLevel string
	switch level {
	case logger.Silent:
		ulogLevel = ulog.ErrorLevelStr // 静默模式，设置为最高级别
	case logger.Error:
		ulogLevel = ulog.ErrorLevelStr
	case logger.Warn:
		ulogLevel = ulog.WarnLevelStr
	case logger.Info:
		ulogLevel = ulog.InfoLevelStr
	default:
		ulogLevel = ulog.InfoLevelStr
	}

	ulog.SetLogLevelStr(ulogLevel)
	return u
}

// Info 输出信息级别日志
func (u *uberLogger) Info(_ context.Context, msg string, data ...any) {
	if len(data) > 0 {
		msg = fmt.Sprintf(msg, data...)
	}
	ulog.Info(msg)
}

// Warn 输出警告级别日志
func (u *uberLogger) Warn(_ context.Context, msg string, data ...any) {
	if len(data) > 0 {
		msg = fmt.Sprintf(msg, data...)
	}
	ulog.Warn(msg)
}

// Error 输出错误级别日志
func (u *uberLogger) Error(_ context.Context, msg string, data ...any) {
	if len(data) > 0 {
		msg = fmt.Sprintf(msg, data...)
	}
	ulog.Error(msg)
}

// Trace 输出 SQL 跟踪日志
func (u *uberLogger) Trace(_ context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)

	if fc == nil {
		return
	}

	sql, rows := fc()

	fields := []zap.Field{
		zap.Duration("elapsed", elapsed),
		zap.String("sql", sql),
		zap.Int64("rows", rows),
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		ulog.Error("SQL execution failed", fields...)
	} else {
		// 根据执行时间决定日志级别
		if elapsed > 500*time.Millisecond {
			ulog.Warn("Slow SQL query detected", fields...)
		} else {
			ulog.Debug("SQL query executed", fields...)
		}
	}
}
