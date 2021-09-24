package semap

import (
	"math"
)

var (
	numbs uint64   //切分的份数,取素数
	nps   []uint64 //切分uint64为211份,从小到大排列
)

func init() {
	//缺省使用509作为素数 -- 致敬谢耳朵最喜欢的素数
	SetupPrime(73)
}

//SetupPrime 设置分片的素数
func SetupPrime(prime uint64) {
	numbs = prime
	var x uint64 = math.MaxUint64
	var y = x / numbs
	nps = make([]uint64, numbs)
	for i := uint64(0); i < numbs; i++ {
		nps[i] = y * (uint64(i) + 1)
	}
	nps[numbs-1] = math.MaxUint64
}
