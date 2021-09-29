## 针对某个key的读写锁

我们常常面对这样的情况，多个go routine可能要访问同一个资源，但这个资源并不一定常驻内存，
换句话说没有固定的资源锁来控制这多个go routine的同步。

### KeyLocker
如果将这些资源以某种唯一key的形式组织，那就可以利用一个map来保存这些资源锁了，在使用的时候会从map中获取，
如果没有就在map中新建。使用完了后再检查是否有别的go routine也将使用这个资源锁，如果没有，就可以将资源锁从
map中删除。

用法
```go

var locker = NewKeyLocker()


//读锁
locker.RLock(resource_id)
defer locker.RUnlock(resource_id)

//写锁
locker.Lock(resource_id)
defer locker.Unlock(resource_id)
```
