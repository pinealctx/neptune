package timex

import (
	"github.com/urfave/cli/v2"
	"time"
)

// LocalTime parse time from command line
// c -- cli context
// name -- cli field name
// github.com/urfave/cli/v2 use UTC as time parse zone.
// actually in most cases, use time.Local zone is good choice.
// 获取命令行参数中相关的时间，此时间以系统时区为基准，而不是UTC时区
// 例如当前时区为北京时区，假如命令行输入为"2021-11-01"，此函数会将此输入解释成 "2021-11-01"对应的北京时间。
// 如果不用此函数，直接用*cli.Context的Timestamp来解析，就会解析成"2021-11-01"对应的UTC时间，实际上它比北京时间早8个小时。
func LocalTime(c *cli.Context, name string) *time.Time {
	var ts = c.Timestamp(name)
	if ts == nil {
		return nil
	}
	var tv = ts.Add(_localDiff).Local()
	return &tv
}
