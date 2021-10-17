# neptune
golang busybox

## bitmap1024    
将16个uint64(128字节)组成一个位图，适合的场景类似于消息信箱的删除信息记录。   
在传统的模式中，如果一个消息需要扩散到别的地方时，往往需要在别的地方同样记录下这样的拷贝，这种case被称为写放大。
这种模式好的地方是拷贝的数据独立，逻辑简单。不好的地方是会写入大量数据。如果所有的数据都共享，只有在发生变化时才做记录，
类似于Unix fork的写时复制，会节约很多写入操作。但消息信箱是一个比较特殊的场景，它对数据的更改只发生在删除的时候，如果
我们只记录删除数据，如何通过记录的删除数据来推导实际上应有的数据呢？用位图是一个不错的选择。  
[文档](./bitmap1024/README.md)

## bytex   
go内置的bytes.Buffer功能很丰富，但它更多面向的是bytes这样的数据类型，bytex封装了bytes.Buffer，丰富了更多的基本类型操作。   
[文档](./bytex/README.md)

## cache
封装了Map与LRU
[文档](./cache/README.md)

## captcha
滑块或图片验证码
[文档](./cache/README.md) 

## compress
压缩相关的包
[文档](./compress/README.md)

## cryptx
加解密相关的包
[文档](./cryptx/README.md)

## ds
数据结构相关的包
[文档](./ds/README.md)

## idgen
生成id相关的包
[文档](./idgen/README.md)

## jsonx
更快的json包
[文档](./jsonx/README.md)

## remap
用于将健值分组的一个工具包
[文档](./remap/README.md)

## store
封装了一些数据中间件客户端代码
[文档](./store/README.md)

## strvali
封装了字符串校验工具函数
[文档](./strvali/README.md)

## syncx
与go routine互斥/同步控制相关的包
[文档](./syncx/README.md)

## tex
扩展了一些基本类型
[文档](./tex/README.md)

## timex
时间相关的包
[文档](./timex/README.md)

## ulog
封装了uber的一个日志库
[文档](./ulog/README.md)

## vcode
验证码缓存与验证的通用模块(手机验证码和邮件验证码)
[文档](./vcode/README.md)