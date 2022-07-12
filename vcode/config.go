package vcode

import (
	"github.com/pinealctx/neptune/tex"
)

type Config struct {
	// LRU cache size.
	CacheSize int64 `json:"cache_size" toml:"cache_size"`
	// Mock mode, the verification code is the last `CodeLen` digits of the mobile phone.
	Mock bool `json:"mock" toml:"mock"`
	// Random code length.
	CodeLen int `json:"code_len" toml:"code_len"`
	// Code TTL.
	TTL tex.Duration `json:"ttl" toml:"ttl"`
	// Minimum interval for sending captcha.
	MinInterval tex.Duration `json:"min_interval" toml:"min_interval"`
	// Captcha count cycle.
	CounterDuration tex.Duration `json:"counter_duration" toml:"counter_duration"`
	// Verification code count upper limit.
	MaxCount int `json:"max_count" toml:"max_count"`
	// Upper limit of verification times of captcha.
	MaxVerifyCount int `json:"max_verify_count" toml:"max_verify_count"`
}
