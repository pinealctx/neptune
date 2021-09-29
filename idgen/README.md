## idgen
id生成相关模块

### idgen/nano
使用time.Now().UnixNano()，当前纳秒时间戳作为id，为int64类型。
并维护了一个全局计数，如果在time.Now()时因为NTP或别的网络时间服务导致了时间回退，会通过检查全局计数的方式。
如果当前纳秒时间戳小于全局计数，则说明了id变小，此时生成的id不是当前纳秒时间戳，而是在全局计数基础上加1。
通过这种方法可以保证在有些场景下产生的id是递增的。另外，此id生成器支持在生成时加载一个指定的全局计数，
这样就算服务重启，也会通过读取数据库的方式来加载上次服务关闭前最大的id，从而保证服务重启也能维护id单调递增。

### idgen/random
提供生成随机id或字符串的函数

```go
//通过纳秒时间戳生成随机id，大部分场景够用，但安全级别不算高，效率很高。
//RandomI64 random generate int64, rand seed is current nano timestamp
func RandomI64() int64

//安全级别高(够随机)，但效率不算高。
//通过读取crypto/rand.Reader作为seed，在类Unix系统中为设备/dev/random或/dev/urandom 
//   /dev/random收集了硬件噪声，
//   /dev/urandom是/dev/random的拷贝，如果内容没有更新，可以重复读取内容并不会阻塞。
//SecRandomI64 random generate int64, read data from /dev/random or /dev/urandom
//actually, it's almost a random id
func SecRandomI64() int64


//生成随机字符串，安全级别不算高，效率高，以纳秒时间戳为seed
func GenNonceStr(baseStr string, length int) string {

//生成随机字符串，安全级别高，效率不算高，会读取crypto/rand.Reader作为seed。
func SecGenNonceStr(baseStr string, length int) string

//先生成uuid，在使用md5 hash，能用到的地方不多
//MD5UUID new uuid first, use md5 hash
func MD5UUID() string

//先生成uuid，在使用sha256 hash，能用到的地方不多
//SHA256UUID new uuid first, use sha256 hash
func SHA256UUID() string
```

### idgen/snowflake
雪花算法中节点信息也是关键的一环，不同机器或服务使用不同节点id会维护全局id不重复。另外，雪花算法中也没有时间回退问题，
因为使用了单调时间。
[文档](./snowflake/README.md)
