## 封装了TCP Server的框架代码

- 支持Option，如果不传入，则使用缺省的Option
- 使用方式，使用时需要传入实现好的IConnHandler接口。

IConnHandler接口有两个方法：
1. ConnCount 获取当前连接数数量，由管理连接的容器负责，或者使用一个原子计数
2. Do(conn net.Conn) 对当前连接进行处理

```go

func start() {
	var s = NewSTCPSrv(SERVER_ADDRESS, Your_IConnMgr)
	//假设最大连接数为100
	var eh = s.Start(WithMaxConn(100))
	
	... //do other things
	
	select {
	case e := <- eh:
	    ... //handle error 
    }   
}
```

- 使用NewTCPSrvX这个函数更简单，只需要实现接口IMsgReader，此接口负责读取数据处理
  其中Session的函数 Read表示读取指定的切片大小的数据，如果失败返回error

```go

//实现IMsgReader接口
type Your_IMsgReader {
    ...
}

func (y *Your_IMsgReader) Read(s *Session) error {
	var head [16]byte
	var data [64]byte
	var err = s.Read(head[:])
	if err != nil {
	    return err
	}
	err = s.Read(data[:])
	if err != nil {
	    return errr
	}   
	//now head and data all be read
	...
}

func start() {
	var s = NewSTCPSrvX(SERVER_ADDRESS, Your_IMsgReader)
	//假设最大连接数为100
	var eh = s.Start(YOUR_HANDLER, WithMaxConn(100))
	...
}
```
