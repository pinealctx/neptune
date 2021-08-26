package gormx

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

const (
	DefaultMaxOpenCount = 40
	DefaultMaxIdleCount = 20
	DefaultMaxLifeTime  = time.Hour
)

//db option
type _Option struct {
	maxOpenConn int
	maxIdle     int
	maxLifeTime time.Duration
	log         bool
}

type Option func(*_Option)

func WithMaxOpenConn(maxOpenConn int) Option {
	return func(o *_Option) {
		o.maxOpenConn = maxOpenConn
	}
}

func WithMaxIdle(maxIdle int) Option {
	return func(o *_Option) {
		o.maxIdle = maxIdle
	}
}

func WithMaxLifeTime(maxLifeTime time.Duration) Option {
	return func(o *_Option) {
		o.maxLifeTime = maxLifeTime
	}
}

func WithLog() Option {
	return func(o *_Option) {
		o.log = true
	}
}

func New(dsn string, opts ...Option) (*gorm.DB, error) {
	var option = &_Option{
		maxOpenConn: DefaultMaxOpenCount,
		maxIdle:     DefaultMaxIdleCount,
		maxLifeTime: DefaultMaxLifeTime,
	}
	for _, opt := range opts {
		opt(option)
	}
	var config = &gorm.Config{}
	if option.log {
		config.Logger = logger.Default.LogMode(logger.Info)
	}
	var gormDB, err = gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		return gormDB, err
	}
	if option.log {
		gormDB = gormDB.Debug()
	}
	var db *sql.DB
	db, err = gormDB.DB()
	if err != nil {
		return nil, err
	}
	if option.maxOpenConn != 0 {
		db.SetMaxOpenConns(option.maxOpenConn)
	}
	if option.maxIdle != 0 {
		db.SetMaxIdleConns(option.maxIdle)
	}
	if option.maxLifeTime != 0 {
		db.SetConnMaxLifetime(option.maxLifeTime)
	}
	return gormDB, err
}

//获取一个db客户端 --maxOpenConn 最大打开连接数
func NewGorm(dsn string, maxOpenConn, maxIdle int, maxLifeTime time.Duration, log bool) (*gorm.DB, error) {
	var config = &gorm.Config{}
	if log {
		config.Logger = logger.Default.LogMode(logger.Info)
	}
	var gormDB, err = gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		return gormDB, err
	}
	if log {
		gormDB = gormDB.Debug()
	}
	var db *sql.DB
	db, err = gormDB.DB()
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxOpenConn)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(maxLifeTime)
	return gormDB, err
}
