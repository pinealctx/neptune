package remap

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/pinealctx/neptune/idgen/snowflake"
)

// 测试平分
func TestUint64Split(t *testing.T) {
	var x uint64 = math.MaxUint64
	t.Log(x)
	var numbs = DefaultPrime
	var y = x / numbs
	for i := uint64(0); i < numbs; i++ {
		fmt.Printf("%+v,\n", y*(uint64(i)+1))
	}
	t.Log(x - y*numbs)
	t.Log(x + 1)
	var z = -9223372036854775808
	z *= -1
	t.Log(z)
	t.Log(math.MinInt64)
	t.Log(math.MaxInt64)
}

// 测试Hash时间
func TestHashKey(t *testing.T) {
	var t1 = time.Now()
	for i := 0; i < 100000; i++ {
		XXHash(i)
	}
	var d = time.Since(t1)
	t.Log("use time:", d, "average:", d/100000)
}

func TestSearchIndex(t *testing.T) {
	var r = NewReMap()
	var vs = make([]uint64, len(r.nps))
	copy(vs, r.nps)
	for i := 0; i < len(vs); i++ {
		var v = r.SearchIndex(vs[i])
		if v != i {
			panic(i)
		}
		t.Log(v)
	}
	t.Log("")
	t.Log("")
	for i := 0; i < len(vs); i++ {
		var v = r.SearchIndex(vs[i] - 1)
		if v != i {
			panic(i)
		}
		t.Log(v)
	}
	t.Log("")
	t.Log("")
	for i := 0; i < len(vs); i++ {
		var v = r.SearchIndex(vs[i] + 1)
		if v != (i+1)%int(r.numbs) {
			panic(i)
		}
		t.Log(v)
	}

	t.Log(r.SearchIndex(0))
	var y uint64 = 0xFFFFFFFFFFFFFFFF
	t.Log(y)
	t.Log(r.SearchIndex(y))
}

func TestSearchIndexSite(t *testing.T) {
	var r = NewReMap()
	var m = make(map[int]int)
	for i := int32(1000010); i < int32(1000010+1000); i++ {
		var x = r.XHashIndex(i)
		m[x]++
	}
	t.Log(m)

	m = make(map[int]int)
	for i := int32(10000); i < int32(10000+1000); i++ {
		var x = r.XHashIndex(i)
		m[x]++
	}
	t.Log(m)

	var node, _ = snowflake.NewMonoNode(0)
	m = make(map[int]int)
	for i := 0; i < 1000; i++ {
		var x = r.XHashIndex(node.Generate())
		m[x]++
	}
	t.Log(m)

	m = make(map[int]int)
	for i := int32(1000010); i < int32(1000010+1000); i++ {
		var x = r.SimpleIndex(i)
		m[x]++
	}
	t.Log(m)

	m = make(map[int]int)
	for i := int32(10000); i < int32(10000+1000); i++ {
		var x = r.SimpleIndex(i)
		m[x]++
	}
	t.Log(m)

	m = make(map[int]int)
	for i := 0; i < 1000; i++ {
		var x = r.SimpleIndex(node.Generate())
		m[x]++
	}
	t.Log(m)
}
