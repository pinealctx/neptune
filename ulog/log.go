package ulog

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	//level string
	DebugLevelStr = "debug"
	InfoLevelStr  = "info"
	WarnLevelStr  = "warn"
	ErrorLevelStr = "error"

	//level define
	DebugLevel = zap.DebugLevel
	InfoLevel  = zap.InfoLevel
	WarnLevel  = zap.WarnLevel
	ErrorLevel = zap.ErrorLevel
)

var (
	//ErrZapEmpty empty zap logger
	ErrZapEmpty = errors.New("zap.pointer.empty")
)

// ZapConf zap configuration
type ZapConf struct {
	//文件名
	FileName string `json:"fileName" toml:"fileName"`
	//文件最大长度 (M bytes)
	MaxSize int `json:"maxSize" toml:"maxSize"`
	//文件最久天数
	MaxAge int `json:"maxAge" toml:"maxAge"`
	//最大备份文件数量
	MaxBackups int `json:"maxBackups" toml:"maxBackups"`

	//是否在备份文件格式中使用utc time，否则为local time
	UTCTime bool `json:"utcTime" toml:"utcTime"`
	//是否压缩旧的文件
	Compress bool `json:"compress" toml:"compress"`
	//Disable log caller
	DisableCaller bool `json:"disableCaller" toml:"disableCaller"`

	//log level
	LogLevel string `json:"logLevel" toml:"logLevel"`
}

// zap conf string
func (c ZapConf) String() string {
	return fmt.Sprintf("f:%s, max:%d, age:%d, backs:%d, utc:%+v, compress:%+v, disableCaller:%+v, level:%s",
		c.FileName, c.MaxSize, c.MaxAge, c.MaxBackups, c.UTCTime, c.Compress, c.DisableCaller, c.LogLevel)
}

// Logger 封装可以动态设置log level的logger
type Logger struct {
	*zap.Logger
	atomicLevel zap.AtomicLevel
}

// SetLevel set level
func (z *Logger) SetLevel(level zapcore.Level) {
	if z == nil {
		return
	}
	z.atomicLevel.SetLevel(level)
}

// SetLevelStr set level string
func (z *Logger) SetLevelStr(sLevel string) {
	if z == nil {
		return
	}
	var level = ParseLevel(sLevel)
	z.atomicLevel.SetLevel(level)
}

// Level get level
func (z *Logger) Level() zapcore.Level {
	if z == nil {
		//a strange number
		return zapcore.Level(127)
	}
	return z.atomicLevel.Level()
}

// LevelStr get level str
func (z *Logger) LevelStr() string {
	if z == nil {
		return "unknown:empty"
	}
	return Level2String(z.atomicLevel.Level())
}

/*------following add nil pointer protect-------*/

func (z *Logger) Debug(msg string, fields ...zap.Field) {
	if z == nil {
		return
	}
	z.Logger.Debug(msg, fields...)
}

func (z *Logger) Info(msg string, fields ...zap.Field) {
	if z == nil {
		return
	}
	z.Logger.Info(msg, fields...)
}

func (z *Logger) Warn(msg string, fields ...zap.Field) {
	if z == nil {
		return
	}
	z.Logger.Warn(msg, fields...)
}

func (z *Logger) Error(msg string, fields ...zap.Field) {
	if z == nil {
		return
	}
	z.Logger.Error(msg, fields...)
}

func (z *Logger) Sync() error {
	if z == nil {
		return ErrZapEmpty
	}
	return z.Logger.Sync()
}

// NewFileLogger new zap logger
func NewFileLogger(cnf ZapConf) *Logger {
	var hook = lumberjack.Logger{
		Filename:   cnf.FileName,
		MaxSize:    cnf.MaxSize,
		MaxAge:     cnf.MaxAge,
		MaxBackups: cnf.MaxBackups,
		//If you want to use local time in back file format, open the flag
		LocalTime: !cnf.UTCTime,
		//If you want to compress back file, open the flag
		Compress: cnf.Compress,
	}
	var writeSync = zapcore.AddSync(&hook)
	var encodeCnf = zap.NewProductionEncoderConfig()
	//If you want to use other time encoder, update code here.
	encodeCnf.EncodeTime = zapcore.ISO8601TimeEncoder
	//If you want to use other logger encoder(not json) here, update code here.
	var atomicLevel = zap.NewAtomicLevel()
	var core = zapcore.NewCore(zapcore.NewJSONEncoder(encodeCnf), writeSync, atomicLevel)
	atomicLevel.SetLevel(ParseLevel(cnf.LogLevel))

	var logger *zap.Logger
	if cnf.DisableCaller {
		logger = zap.New(core)
	} else {
		logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	}
	return &Logger{
		Logger:      logger,
		atomicLevel: atomicLevel,
	}
}

// NewSimpleLogger new zap logger without file log
func NewSimpleLogger(level string, options ...zap.Option) *Logger {
	var encodeCnf = zap.NewProductionEncoderConfig()
	encodeCnf.EncodeTime = zapcore.ISO8601TimeEncoder
	var atomicLevel = zap.NewAtomicLevel()
	var core = zapcore.NewCore(zapcore.NewJSONEncoder(encodeCnf), os.Stdout, atomicLevel)
	atomicLevel.SetLevel(ParseLevel(level))
	var logger = zap.New(core, options...)
	return &Logger{
		Logger:      logger,
		atomicLevel: atomicLevel,
	}
}

// ParseLevel parse zap level
func ParseLevel(level string) zapcore.Level {
	var lLevel = strings.ToLower(level)
	switch lLevel {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	default:
		return ErrorLevel
	}
}

// Level2String extract zap level
func Level2String(level zapcore.Level) string {
	switch level {
	case DebugLevel:
		return DebugLevelStr
	case InfoLevel:
		return InfoLevelStr
	case WarnLevel:
		return WarnLevelStr
	case ErrorLevel:
		return ErrorLevelStr
	default:
		return "unknown:" + strconv.Itoa(int(level))
	}
}
