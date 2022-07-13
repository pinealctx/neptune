## 针对某个key的读写锁

我们常常面对这样的情况，多个go routine可能要访问同一个资源，但这个资源并不一定常驻内存，
换句话说没有固定的资源锁来控制这多个go routine的同步。

### KeyLocker
如果将这些资源以某种唯一key的形式组织，那就可以利用一个map来保存这些资源锁了，在使用的时候会从map中获取，
如果没有就在map中新建。使用完了后再检查是否有别的go routine也将使用这个资源锁，如果没有，就可以将资源锁从
map中删除。

### KeyLocker组
生成一组KeyLocker，一般此数组长度为素数，通过将key做hash后取模的方式将key映射到不同的KeyLocker上，可以提供更大的并发效率。   
在测试中，KeyLocker组的效率比单个KeyLocker要高，单个KeyLocker平均操作时长在400ns左右，简单分组的KeyLocker平均操作时长在80ns左右，
按xxhash后的分组平均时长在120ns左右。

用法
```go

//新建单个KeyLocker
var locker = NewKeyLocker()
//新建简单分组的KeyLocker组
var locker = NewKeyLockeGrp()
//新建xxhash分组的KeyLocker组
var locker = NewXHashKeyLockeGrp()


//读锁
locker.RLock(resource_id)
defer locker.RUnlock(resource_id)

//写锁
locker.Lock(resource_id)
defer locker.Unlock(resource_id)


```
