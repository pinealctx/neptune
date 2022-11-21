## stringx

### 封装了随机数相关函数

#### Random with time.Now().UnixNano() seed.

- RandInt [int.Min, int.Max] 包括负数
- RandUint [uint.Min, uint.Max] 不包括负数
- RandInt64 [int64.Min, int64.Max] 包括负数
- RandUint64 [uint64.Min, uint64.Max] 不包括负数
- RandInt32 [int32.Min, int32.Max] 包括负数
- RandUint32 [uint32.Min, uint32.Max] 不包括负数

### Random with system source(/dev/urandom), more security.

- RandIntSecure [int.Min, int.Max] 包括负数
- RandUintSecure [uint.Min, uint.Max] 不包括负数
- RandInt64Secure [int64.Min, int64.Max] 包括负数
- RandUint64Secure [uint64.Min, uint64.Max] 不包括负数
- RandInt32Secure [int32.Min, int32.Max] 包括负数
- RandUint32Secure [uint32.Min, uint32.Max] 不包括负数

### Random in range.
  Specific m and n(m <= n, otherwise will panic), random range in [m, n].

- RandBetween use time.Now().UnixNano() as each random source.
- SimpleRandBetween use global rand source, in each time call, time.Now().UnixNano() will feed in.
- RandBetweenSecure feed system source(/dev/urandom) as each call source.

### Shuffle
  Shuffle slice.
  - Shuffle function to shuffle a slice.
  - SetShuffleRand set Random in range function.   
    (Option1. RandBetween)  
    (Option2. SimpleRandBetween)   
    (Option3. RandBetweenSecure)   
