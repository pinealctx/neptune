## compress

压缩模块，现在主流的压缩有snappy和zstd，通过对这两者的测试对比，snappy的综合表现不错，
在CPU时间和压缩比上有不错的平衡。

使用方法:

```go
//压缩
var cd = compress.Snappy.Compress(data)

//解压
var dd, err = compress.Snappy.DeCompress(data)

//带前缀压缩，有时我们希望压缩数据内容，但保留数据头，即数据头不用被压缩。
//这样的好处是在某些情况下，可以根据数据头做出一些逻辑，以免所有情况都需要把数据解压。
var hcd = compress.Snappy.CompressWithPrefix(data, yourHead)
```

