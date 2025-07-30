package gormx

import (
	"database/sql"
	"time"

	// import mysql driver
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	DefaultMaxOpenCount = 40
	DefaultMaxIdleCount = 20
	DefaultMaxLifeTime  = time.Hour
)

// db option
type _Option struct {
	maxOpenConn int
	maxIdle     int
	maxLifeTime time.Duration
	log         bool
	logger      logger.Interface
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

func WithLogSwitch(on bool) Option {
	return func(o *_Option) {
		o.log = on
	}
}

// WithLogger : set logger
func WithLogger(logger logger.Interface) Option {
	return func(o *_Option) {
		o.logger = logger
	}
}

// New : new *gorm.DB
func New(dsn string, opts ...Option) (*gorm.DB, error) {
	var option = &_Option{
		maxOpenConn: DefaultMaxOpenCount,
		maxIdle:     DefaultMaxIdleCount,
		maxLifeTime: DefaultMaxLifeTime,
	}
	for _, opt := range opts {
		opt(option)
	}
	var config = &gorm.Config{Logger: option.logger}
	if config.Logger == nil && option.log {
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

// NewDBByDsn : new *gorm.DB by Dsn
func NewDBByDsn(dsn *Dsn, opts ...Option) (*gorm.DB, error) {
	return New(dsn.UseDefault(), opts...)
}

// NewGorm 获取一个db客户端 --maxOpenConn 最大打开连接数
func NewGorm(dsn string, maxOpenConn, maxIdle int, maxLifeTime time.Duration, log bool) (*gorm.DB, error) {
	return New(dsn,
		WithMaxOpenConn(maxOpenConn), WithMaxIdle(maxIdle), WithMaxLifeTime(maxLifeTime), WithLogSwitch(log))
}

// NewDBBySSH new db by ssh
func NewDBBySSH(sshCnf *SSHConfig, dsn *Dsn, opts ...Option) (*gorm.DB, error) {
	var sshCli, err = CreateSSHConn(sshCnf)
	if err != nil {
		return nil, err
	}
	var sshDialer = &SSHDialer{
		client: sshCli,
	}
	sshDialer.Register()

	var db *gorm.DB
	db, err = New(dsn.UseDefault(), opts...)
	if err != nil {
		return nil, err
	}
	return db, err
}
