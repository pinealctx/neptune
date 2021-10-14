package snowflake

import (
	"fmt"
	"github.com/pinealctx/neptune/tex"
	"strconv"
	"time"
)

/*
典型的雪花算法的信息布局
最高位b0   ---->  最低位b63
[b0]      [b1 ---- b41]    [b42 --- b51]     [b52 --- b63]
b0: 作为符号位，不用，设为0。
b1 -- b41: 41位毫秒时间戳，可用69年。
b42 -- b51: 10位节点信息，可设置1024个节点。
b52 - b63: 同一毫秒时的步长(为4096)，在同一毫秒产生ID时，通过增加步长来产生不同ID，
意味着1秒内可以产生4096000个ID，应付大部分场景足够了。

按实际情况，可以适当做些调整与变化:
1. 有些场景下节点信息放在最低位可能更合理，此时，整个结构就变成了
[b0]      [b1 ---- b41]    [b42 --- b53]     [b54 --- b63]
设为0      毫秒时间戳         步长信息           节点信息
2. 如果需要的节点数不用那么多，可以将节省下来的位数给时间戳用。
例如如果只支持256个节点，则毫秒数可以支持的时间则变为69*4 = 276年
3. 为了方便使用，将节点位数可以简化成3种制式:
Node1024 -- 支持1024个节点，时间戳支持为69年
Node512 -- 支持512个节点，时间戳支持为138年
Node256 -- 支持256个节点，时间戳支持为276年
同时，将节点信息存放位置也简化成2种制式:
NodeAtLowest:
true -- 节点信息放在最低位
false -- 步长信息放在最低位
*/

var (
	//全局的创世的毫秒时间戳
	_epoch int64 = 1609430400000
	//全局的node bits
	_nodeBits uint8 = 10
	//全局的n节点是否放在最低位
	_nodeAtLowest = false

	timeLoc, _ = time.LoadLocation("Asia/Shanghai")
)

const (
	//StepBits 步长永远设置为12 bits
	StepBits uint8 = 12

	Node1024 NodeBitsMode = 10
	Node512  NodeBitsMode = 9
	Node256  NodeBitsMode = 8

	//SDivMs 1 second = 1000 ms
	SDivMs = 1000
	//MsDivNs 1 ms = 1000000 ns
	MsDivNs = 1000000

	//TimeStrLen convert id to time style string, its length is fixed 24
	TimeStrLen = 24
)

//Node : generate id interface
type Node interface {
	Generate() int64
}

//NodeBitsMode node bit mode
type NodeBitsMode uint8

//_Option : snowflake option
type _Option struct {
	//创世的毫秒时间戳，假如以2021年开始时间为创世时间，则需要保证后面ID生成的过程中，不会出现早于2021年的系统时间
	epoch int64
	//节点位数
	nodeBits uint8
	//节点是否放在最低位
	nodeAtLowest bool
}

type Option func(o *_Option)

//UseEpoch : 设置创世时间
func UseEpoch(t time.Time) Option {
	return func(o *_Option) {
		o.epoch = t.UnixNano() / int64(time.Millisecond)
	}
}

//UseNodeMode : 设置节点位数模式
func UseNodeMode(m NodeBitsMode) Option {
	return func(o *_Option) {
		switch m {
		case Node256:
		case Node512:
		default:
			m = Node1024
		}
		o.nodeBits = uint8(m)
	}
}

//NodeAtLowest : 设置节点位数模式
func NodeAtLowest() Option {
	return func(o *_Option) {
		o.nodeAtLowest = true
	}
}

//Setup setup snowflake
func Setup(opts ...Option) {
	var o = &_Option{
		epoch:        _epoch,
		nodeBits:     _nodeBits,
		nodeAtLowest: _nodeAtLowest,
	}

	for _, opt := range opts {
		opt(o)
	}
	_epoch = o.epoch
	_nodeBits = o.nodeBits
	_nodeAtLowest = o.nodeAtLowest
}

// IDFields figure out time/node/step from id
// Return : ms timestamp, node, step
func IDFields(id int64) (timeF, node, step int64) {
	var nodeMax int64 = (1 << _nodeBits) - 1
	var stepMax int64 = (1 << StepBits) - 1
	var timeShift, nodeShift, stepShift = figureShift()

	timeF = id >> timeShift
	node = (id >> nodeShift) & nodeMax
	step = (id >> stepShift) & stepMax
	return timeF, node, step
}

// IDParse figure out time/node/step from id
// Return : ms timestamp, node, step
func IDParse(id int64) (timeMs, node, step int64) {
	timeMs, node, step = IDFields(id)
	timeMs += _epoch
	return timeMs, node, step
}

// IDParseEx figure out time/node/step from id
// Return : time.Time, node, step
func IDParseEx(id int64) (t time.Time, node, step int64) {
	var ts int64
	ts, node, step = IDParse(id)
	t = time.Unix(ts/SDivMs, (ts%SDivMs)*MsDivNs).In(timeLoc)
	return t, node, step
}

// TimeIDRange figure out a specific time min and max id
// The calculation is based on second
func TimeIDRange(t time.Time) (min, max int64) {
	var timeShift = _nodeBits + StepBits
	var ts = t.Unix()
	var timeMs = ts*SDivMs - _epoch
	timeMs <<= timeShift
	min = timeMs
	var reMax int64 = (1 << timeShift) - 1
	max = timeMs | reMax
	return min, max
}

// TimeBetweenID figure out a specific [after time, before time]
// The calculation is based on second
func TimeBetweenID(begin time.Time, end time.Time) (min, max int64) {
	var timeShift = _nodeBits + StepBits
	var beginTS = begin.Unix()
	var beginMs = beginTS*SDivMs - _epoch
	min = beginMs << timeShift
	var reMax int64 = (1 << timeShift) - 1
	var endTS = end.Unix()
	var endMs = endTS*SDivMs - _epoch
	max = (endMs << timeShift) | reMax
	return min, max
}

// CnStyle
// A weird implement for chinese style
// 总共长度24
// 17位时间- 20210901 003859 000
// 7位NODE+STEP
func CnStyle(id int64) string {
	var timeShift = _nodeBits + StepBits
	var ms = (id >> timeShift) + _epoch
	var t = time.Unix(ms/SDivMs, (ms%SDivMs)*MsDivNs).In(timeLoc)
	var buf = tex.NewSizedBuffer(TimeStrLen)
	var mask int64 = (1 << timeShift) - 1
	var left = id & mask
	_, _ = buf.WriteString(fmt.Sprintf("%04d", t.Year()))
	_, _ = buf.WriteString(fmt.Sprintf("%02d", t.Month()))
	_, _ = buf.WriteString(fmt.Sprintf("%02d", t.Day()))
	_, _ = buf.WriteString(fmt.Sprintf("%02d", t.Hour()))
	_, _ = buf.WriteString(fmt.Sprintf("%02d", t.Minute()))
	_, _ = buf.WriteString(fmt.Sprintf("%02d", t.Second()))
	_, _ = buf.WriteString(fmt.Sprintf("%03d", t.Nanosecond()/MsDivNs))
	_, _ = buf.WriteString(fmt.Sprintf("%07d", left))
	return buf.String()
}

// FromChStyle from string to int64
func FromChStyle(v string) (int64, error) {
	var l = len(v)
	if l != TimeStrLen {
		return 0, fmt.Errorf("unspported.id.cn.len:%d", l)
	}
	var (
		year int
		es   [5]int
		ms   int
		left int
		err  error
	)
	year, err = strconv.Atoi(v[:4])
	if err != nil {
		return 0, err
	}
	for i := 0; i < 5; i++ {
		es[i], err = strconv.Atoi(v[4+i*2 : 4+(i+1)*2])
		if err != nil {
			return 0, err
		}
	}
	ms, err = strconv.Atoi(v[14:17])
	if err != nil {
		return 0, err
	}
	left, err = strconv.Atoi(v[17:])
	if err != nil {
		return 0, err
	}
	var t = time.Date(year, time.Month(es[0]), es[1], es[2], es[3], es[4], ms*MsDivNs, timeLoc)
	var ns = t.UnixNano()
	var tms = ns/MsDivNs - _epoch
	var timeShift = _nodeBits + StepBits
	var id = tms << timeShift
	id |= int64(left)
	return id, nil
}

// figureShift : calculate time shift, node shift, step shift.
func figureShift() (timeShift, nodeShift, stepShift uint8) {
	timeShift = _nodeBits + StepBits
	if _nodeAtLowest {
		nodeShift = 0
		stepShift = _nodeBits
	} else {
		stepShift = 0
		nodeShift = StepBits
	}
	return timeShift, nodeShift, stepShift
}
