package snowflake

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

//MonoNode use monotonic time compare epoch

// A MonoNode struct holds the basic information needed for a snowflake generator
// node
type MonoNode struct {
	mu    sync.Mutex
	epoch time.Time
	time  int64
	node  int64
	step  int64
}

// NewMonoNode returns a new snowflake node that can be used to generate snowflake
func NewMonoNode(node int64) (Node, error) {
	var nodeMax int64 = (1 << _nodeBits) - 1
	if node < 0 || node > nodeMax {
		return nil, errors.New("node.number.must.be.between.0.and." + strconv.FormatInt(nodeMax, 10))
	}
	var n = &MonoNode{}
	n.node = node
	var curTime = time.Now()
	// add time.Duration to curTime to make sure we use the monotonic clock if available
	n.epoch = curTime.Add(time.Unix(_epoch/SDivMs, (_epoch%SDivMs)*MsDivNs).Sub(curTime))
	return n, nil
}

// Generate creates and returns a unique snowflake ID
// To help guarantee uniqueness
// - Make sure your system is keeping accurate system time
// - Make sure you never have multiple nodes running with the same node ID
func (n *MonoNode) Generate() int64 {
	n.mu.Lock()
	defer n.mu.Unlock()

	var stepMax int64 = (1 << StepBits) - 1
	var now = time.Since(n.epoch).Nanoseconds() / MsDivNs
	if now == n.time {
		n.step = (n.step + 1) & stepMax
		if n.step == 0 {
			for now <= n.time {
				now = time.Since(n.epoch).Nanoseconds() / MsDivNs
			}
		}
	} else {
		n.step = 0
	}

	n.time = now
	var timeShift, nodeShift, stepShift = figureShift()
	var r = now<<timeShift | n.node<<nodeShift | n.step<<stepShift
	return r
}
