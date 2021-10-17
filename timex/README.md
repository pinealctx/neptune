## timex
时间相关的包

### LocalDiff()
获取当前系统的时区与UTC时区的差异时间
```go
//如果当前系统设定时区为北京时间，则返回8小时，因为北京时间在+8区
var diff = LocalDiff()
```

### 天相关的时间函数
计算机中有一个闰秒的时间概念，从time包中我们可以找到微秒/毫秒/秒/分/小时这些时间单位对应的纳秒数常量，但为什么没有Day相关的常量呢？这是因为用24小时实际上并不能准备地表示出1天的时间，因为闰秒的关系。  
大部分时候用24小时来表示一天没有太大的问题。为了更严谨地处理时间中的天，这里使用了time.Time这个构造函数，此构造函数可以准确地处理时间中与天相关的部分。

```go
//TodayBegin : today time begin
//获取今天起始时间
func TodayBegin() time.Time

//TodayDeltaBegin : one day begin with specific delta with today
//指定天数差异，可以获取早于今天n天的一天起始时间或晚于今天n天的一天起始时间
//其中由n为正或负来决定
//例如，昨天的起始时间可以由TodayDeltaBegin(-1)来获取，
//明天的起始时间可以由TodayDeltaBegin(1)来获取。
func TodayDeltaBegin(n int) time.Time

//TodayDeltaDayBegins : a series day begin list which have been specified delta day delta list
//与TodayDeltaBegin类似，不同的是通过指定一系列的天数差异来获取不同差值天数的起始时间
//例如，TodayDeltaDayBegins(-2, -1, 1, 2)可以分别获取前天/昨天/明天/后天各自的起始时间
func TodayDeltaDayBegins(diffs ...int) []time.Time

//DayBegin : input a time to figure its day begin
//指定一个时间，获取此时间所在的那天起始时间
func DayBegin(at time.Time) time.Time

//DayDeltaBegin : input a time to figure its delta day begin
//指定一个时间和天数差异，获取早于此时间n天或晚于此时间n天的那天的起始时间
//其中由n为正或负来决定
//例如，比输入时间早一天的起始时间可以由DayDeltaBegin(t, -1)来获取，
//比输入时间晚一天的起始时间可以由DayDeltaBegin(t, 1)来获取。
func DayDeltaBegin(at time.Time, n int) time.Time

//DayDeltaBegins : input a time to calculate a delta list day begins
//与DayDeltaBegin类似，不同的是通过指定一系列的天数差异来获取不同差值天数的起始时间
//例如，DayDeltaBegins(t, -2, -1, 1, 2)可以分别获取比t早2天/早1天/晚1天/晚2天的起始时间
func DayDeltaBegins(at time.Time, diffs ...int) []time.Time
```

### 命令行参数时间
github.com/urfave/cli/v2是我们经常使用的命令行参数包，实际上它提供了时间参数，可以自行设定时间格式，其原理是通过time.Parse来解析一个字符串并生成对应的时间。但此包用的时间解析是基于UTC的时间解析，如果在参数输入时间时需要换算UTC时间，在使用时往往不那么直观。   
对于以中国时区为基准的很多应用来说，往往我们输入的时间其实上就是指的北京时间，这里实现了一个函数:   
#### func LocalTime(c *cli.Context, name string) *time.Time   
此函数即将命令行中相关的时间参数解析成系统时区对应的时间，例如当前时区为北京时区，假如命令行输入为"2021-11-01"，此函数会将此输入解释成 "2021-11-01"对应的北京时间。如果不用此函数，直接用*cli.Context的Timestamp来解析，就会解析成"2021-11-01"对应的UTC时间，实际上它比北京时间早8个小时。
