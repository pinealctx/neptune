package semap

import (
	"fmt"
	"github.com/pinealctx/neptune/idgen/snowflake"
	"math"
	"testing"
	"time"
)

//测试平分
func TestUint64Split(t *testing.T) {
	var x uint64 = math.MaxUint64
	t.Log(x)
	var y = x / numbs
	for i := uint64(0); i < numbs; i++ {
		fmt.Printf("%+v,\n", y*(uint64(i)+1))
	}
	t.Log(x - y*numbs)
	t.Log(x + 1)
	var z int = -9223372036854775808
	z *= -1
	t.Log(z)
	t.Log(math.MinInt64)
	t.Log(math.MaxInt64)
}

//测试Hash时间
func TestHashKey(t *testing.T) {
	var t1 = time.Now()
	for i := 0; i < 100000; i++ {
		XXHash(i)
	}
	var d = time.Now().Sub(t1)
	t.Log("use time:", d, "average:", d/100000)
}

func TestSearchIndex(t *testing.T) {
	var vs = make([]uint64, len(nps))
	copy(vs, nps)
	for i := 0; i < len(vs); i++ {
		var v = SearchIndex(vs[i])
		if v != i {
			panic(i)
		}
		t.Log(v)
	}
	t.Log("")
	t.Log("")
	for i := 0; i < len(vs); i++ {
		var v = SearchIndex(vs[i] - 1)
		if v != i {
			panic(i)
		}
		t.Log(v)
	}
	t.Log("")
	t.Log("")
	for i := 0; i < len(vs); i++ {
		var v = SearchIndex(vs[i] + 1)
		if v != (i+1)%int(numbs) {
			panic(i)
		}
		t.Log(v)
	}

	t.Log(SearchIndex(0))
	var y uint64 = 0xFFFFFFFFFFFFFFFF
	t.Log(y)
	t.Log(SearchIndex(y))
}

func TestSearchIndexSite(t *testing.T) {
	var m = make(map[int]int)
	for i := int32(1000010); i < int32(1000010+1000); i++ {
		var x = XHashIndex(i)
		m[x]++
	}
	t.Log(m)

	m = make(map[int]int)
	for i := int32(10000); i < int32(10000+1000); i++ {
		var x = XHashIndex(i)
		m[x]++
	}
	t.Log(m)

	var node, _ = snowflake.NewNode(0)
	m = make(map[int]int)
	for i := 0; i < 1000; i++ {
		var x = XHashIndex(node.Generate())
		m[x]++
	}
	t.Log(m)

	m = make(map[int]int)
	for i := int32(1000010); i < int32(1000010+1000); i++ {
		var x = SimpleIndex(i)
		m[x]++
	}
	t.Log(m)

	m = make(map[int]int)
	for i := int32(10000); i < int32(10000+1000); i++ {
		var x = SimpleIndex(i)
		m[x]++
	}
	t.Log(m)

	m = make(map[int]int)
	for i := 0; i < 1000; i++ {
		var x = SimpleIndex(node.Generate())
		m[x]++
	}
	t.Log(m)
}
