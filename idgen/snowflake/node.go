package snowflake

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

// A HardNode struct holds the basic information needed for a snowflake generator
// node
type HardNode struct {
	mu    sync.Mutex
	epoch int64
	time  int64
	node  int64
	step  int64
}

// NewNode returns a new snowflake node that can be used to generate snowflake
func NewNode(node int64, min int64) (Node, error) {
	var nodeMax int64 = (1 << _nodeBits) - 1
	if node < 0 || node > nodeMax {
		return nil, errors.New("node.number.must.be.between.0.and." + strconv.FormatInt(nodeMax, 10))
	}
	var n = &HardNode{}
	n.node = node
	n.epoch = time.Unix(_epoch/SDivMs, (_epoch%SDivMs)*MsDivNs).UnixNano() / MsDivNs
	n.time, _, n.step = IDFields(min)
	return n, nil
}

// Generate creates and returns a unique snowflake ID
// To help guarantee uniqueness
func (n *HardNode) Generate() int64 {
	n.mu.Lock()
	defer n.mu.Unlock()
	var now = _HookNow().UnixNano()/MsDivNs - n.epoch

	if now > n.time {
		n.step = 0
		n.time = now
	} else {
		var stepMax int64 = (1 << StepBits) - 1
		n.step = (n.step + 1) & stepMax
		if n.step == 0 {
			n.time++
		}
	}

	var timeShift, nodeShift, stepShift = figureShift()
	var r = n.time<<timeShift | n.node<<nodeShift | n.step<<stepShift
	return r
}

var (
	//just for test hook
	_HookNow = time.Now
)
