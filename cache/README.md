## cache

### LRU: 最近没使用的换出缓存，go routine安全

接口LRUFacade

```go
    // Get returns a value from the cache, and marks the entry as most recently used.
    Get(key interface{}) (v Value, ok bool)
    // Peek returns a value from the cache without changing the LRU order.
    Peek(key interface{}) (v Value, ok bool)
    // Set sets a value in the cache.
    Set(key interface{}, value Value)
    // Delete removes an entry from the cache, and returns if the entry existed.
    Delete(key interface{}) bool
```

```go
//创建容量为10000的单个LRU Cache
var s = NewSingleLRUCache(10000)

//创建容量为20000的多路LRU Cache
var m = NeWideLRUCache(20000)

//创建容量为20000的多路LRU Cache
//并通过remap.WithPrime(211)明确指定将Cache分成211组
//需要传入一个素数让分组更均匀
var m = NeWideLRUCache(20000, remap.WithPrime(211))

//创建容量为30000的多路LRU Cache，其中在对key做映射时，
//对key做xxhash运算然后用求的的uint64值来计算它对应多路Cache中具体哪一个。
var x = NeWideLRUCache(20000)

```

### Map: 对Map进行了读写锁的封装，go routine安全

接口MapFacade
```go
//MapFacade an interface to define a Map
type MapFacade interface {
    //Set : set key-value
    Set(key interface{}, value interface{})
    //Get : get value
    Get(key interface{}) (interface{}, bool)
    //Exist : return true if key in map
    Exist(key interface{}) bool
    //Delete : delete a key
    Delete(key interface{})
}
```

```go
//创建单个Map
var s = NewSingleMap()

//创建容多路Map
var m = NeWideMap()

//创建多路Map
//并通过remap.WithPrime(211)明确指定将Cache分成211组
//需要传入一个素数让分组更均匀
var m = NeWideMap(remap.WithPrime(211))

//创建多路Map，其中在对key做映射时，
//对key做xxhash运算然后用求的的uint64值来计算它对应多路Cache中具体哪一个。
var x = NewWideXHashMap()

```

### 特别说明
对LRU和Map进行分组能提高至少一倍的效率，但就算使用单个Map或LRU本身的效率也非常高，可以根据具体的场景来取舍。   

LRU的读写测试大概耗时为300ns，分组LRU大概耗时为130ns。    
Map的读写测试大概耗时为150ns，分组Map大概耗时为50ns。

