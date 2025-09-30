## TCP Server 框架（基于接口分层设计）

该模块提供一个轻量、高性能的 TCP Server 框架，采用接口分层架构实现关注点分离：
- 通过 NewTcpServer(cnf, readProcessor, connIOFactory) 创建服务，支持工厂模式注入不同的连接处理器
- ReadProcessor 完全由业务方实现，负责处理从连接读取的数据（定长、变长、心跳等策略均由外部控制）
- 框架内部管理连接生命周期、并发处理，通过ConnIOFactory支持灵活的连接实现策略，包括完整的配置和Hook机制

### 核心组件

- **TcpServer**: 主服务器，负责Accept循环和连接管理
- **ConnHandler**: 连接处理器，管理单个连接的生命周期，包含panic恢复机制
- **QSendConn**: 队列连接实现，提供异步消息发送能力（IConnIO的一个实现）
- **IConnSender**: 连接发送器接口，定义消息发送能力，所有方法（除内部loopSend）均为线程安全
- **IConnReader**: 连接读取器接口，定义数据读取能力，支持不同的帧读取策略
- **IConnIO**: 连接IO接口，组合发送器和读取器，提供完整的连接操作能力
- **ConnIOFactory**: 连接IO工厂，用于创建不同类型的连接实例

---

### 接口分层架构

框架采用三层接口设计，实现关注点分离：

#### IConnSender - 连接发送器接口
负责消息发送相关功能，定义发送能力的抽象：
```go
type IConnSender interface {
    // 核心方法
    Conn() net.Conn                    // 获取底层连接 (线程安全)
    SetMetaInfo(m MetaInfo)            // 设置元信息 (线程安全，重入安全)
    MetaInfo() MetaInfo                // 获取元信息 (线程安全)
    Close() error                      // 关闭连接 (线程安全，重入安全)
    
    // 发送方法 (至少实现一个，线程安全，重入安全)
    Put2Queue(bs []byte) error                    // 异步发送到队列
    Put2SendMap(key uint32, bs []byte) error      // 按整型键发送
    Put2SendSMap(key string, bs []byte) error     // 按字符串键发送
    Put2SendMaps(pairs []KeyIntBytesPair) error   // 批量按整型键发送
    Put2SendSMaps(pairs []KeyStrBytesPair) error  // 批量按字符串键发送
    
    // 内部方法 (仅由框架调用，非线程安全)
    loopSend()                         // 发送循环，由ConnHandler管理
}
```

#### IConnReader - 连接读取器接口
负责数据读取相关功能，支持不同的帧读取策略：
```go
type IConnReader interface {
    ReadFrame() ([]byte, error)        // 读取一帧数据 (线程安全)
}
```

#### IConnIO - 连接IO接口
组合发送器和读取器，提供完整的连接操作能力：
```go
type IConnIO interface {
    IConnSender                        // 继承发送能力
    IConnReader                        // 继承读取能力
}
```

#### 工厂接口
```go
type ConnReaderFactory func(conn net.Conn) IConnReader
type ConnIOFactory func(conn net.Conn) IConnIO
```

### API 接口

#### Server 构造
```go
func NewTcpServer(cnf *ServerAcceptCnf, readProcessor ReadProcessor, connIOFactory ConnIOFactory) *TcpServer
```

#### ReadProcessor 函数定义
```go
type ReadProcessor func(iConnIO IConnIO, buffer []byte) error
```
框架会在每个连接上循环调用 ReadProcessor；当返回非 nil 错误时，该连接将被关闭并触发退出钩子。

#### 配置结构
```go
type ServerAcceptCnf struct {
    Address        string         // 监听地址，如 ":9000" 或 "127.0.0.1:9000"
    AcceptDelay    timex.Duration // 初始 accept 退避延迟
    AcceptMaxDelay timex.Duration // 最大 accept 退避延迟
    AcceptMaxRetry int            // 最大 accept 重试次数
    MaxConn        int32          // 最大并发连接数
}

// 获取默认配置
func DefaultServerAcceptCnf() *ServerAcceptCnf
```

#### ConnReader 函数定义
```go
type ConnReaderFunc func(handler IConnSender, conn net.Conn) error
```
框架会在每个连接上循环调用 ConnReaderFunc；当返回非 nil 错误时，该连接将被关闭并触发退出钩子。

#### 启动与关闭
- 启动服务：`Run(errChan chan<- error)`
- 关闭服务：`Close() error`（停止监听，已建立连接会按钩子流程退出）

#### Hook 机制
- 设置连接启动钩子：`SetStartHooker(hooker ConnStartEvent)`
- 设置连接退出钩子：`SetExitHooker(hooker ConnExitEvent)`

Hook函数定义：
```go
type ConnStartEvent func(iConnIO IConnIO)
type ConnExitEvent func(iConnIO IConnIO)
```

**重要**: 框架对Hook函数提供panic恢复机制，Hook函数中的panic不会导致整个连接或服务器崩溃，错误会被记录到日志中。

#### 消息发送
- 异步发送：`iConnIO.Put2Queue([]byte) error`（写入发送队列）
- 连接关闭：`iConnIO.Close() error`（线程安全，支持重入调用）
- 写超时控制：`SetWriteTimeout(d time.Duration)`，默认 5s

#### 连接信息
- 获取监听地址：`Address() string`
- 获取当前连接数：`ConnCount() int32`

#### 日志与元信息
- 默认 MetaInfo 为 BasicMetaInfo{RemoteAddr}
- 可调用 `iConnIO.SetMetaInfo(...)` 替换为自定义 MetaInfo（需实现 zapcore.ObjectMarshaler）以丰富日志字段
- SetMetaInfo 方法线程安全，支持在任何时候更新连接元信息

### IConnIO 接口设计

#### 线程安全与重入性
`IConnIO` 接口及其组合的 `IConnSender`、`IConnReader` 接口的所有公开方法都经过精心设计，确保：

- **线程安全**: 所有方法都可以在多个 goroutine 中并发调用
- **重入安全**: `Close()` 方法支持多次调用，不会 panic，确保资源清理的幂等性
- **错误处理**: 方法可能返回错误，但保证不会 panic
- **panic恢复**: ConnHandler对关键操作提供panic恢复机制，确保单个连接的异常不影响整体服务

#### 方法分类
- **必需方法**: `Conn()`, `SetMetaInfo()`, `MetaInfo()`, `Close()`, `ReadFrame()`
- **可选方法**: `Put2Queue()`, `Put2SendMap()`, `Put2SendSMap()`, `Put2SendMaps()`, `Put2SendSMaps()`
- **内部方法**: `loopSend()` - 仅由框架内部调用，业务代码不应直接使用

#### 生命周期管理
通过 `IConnIO.Close()` 可以优雅地关闭连接：
```go
// 业务代码中的任何地方都可以安全调用
err := iConnIO.Close()  // 线程安全，重入安全
if err != nil {
    // 处理关闭错误，但不会 panic
}
```

调用 `Close()` 后会触发：
1. 发送队列关闭，`loopSend` goroutine 退出
2. 网络连接关闭，`loopReceive` goroutine 退出  
3. `ConnHandler.Exit()` 执行，触发退出钩子

---

## 快速上手示例

### 基础服务启动

```go
package main

import (
    "io"
    "net"
    "time"

    "github.com/pinealctx/neptune/stcp"
    "github.com/pinealctx/neptune/timex"
)

func main() {
    // ReadProcessor 函数：负责处理从连接读取的数据
    readProcessor := func(iConnIO stcp.IConnIO, buffer []byte) error {
        // 示例：回显数据
        // 业务处理...
        return iConnIO.Put2Queue(buffer)  // 异步发送响应
    }

    // 创建连接读取器工厂（实现自定义帧读取逻辑）
    readerFactory := func(conn net.Conn) stcp.IConnReader {
        // 这里可以返回自定义的IConnReader实现
        // 或者使用框架提供的基础实现
        return &MyConnReader{conn: conn}
    }

    // 创建连接IO工厂（使用队列连接实现）
    connIOFactory := func(conn net.Conn) stcp.IConnIO {
        return stcp.NewQSendConnHandler(conn, 1024, readerFactory)  // 队列容量1024
    }

    // 创建服务器配置
    cnf := stcp.DefaultServerAcceptCnf()
    cnf.Address = ":9000"
    cnf.MaxConn = 1000
    
    // 创建服务器
    srv := stcp.NewTcpServer(cnf, readProcessor, connIOFactory)
    
    // 设置连接钩子（可选）
    srv.SetStartHooker(func(iConnIO stcp.IConnIO) {
        // 连接建立时的处理
    })
    srv.SetExitHooker(func(iConnIO stcp.IConnIO) {
        // 连接断开时的处理
    })
    
    // 设置全局写超时（可选）
    stcp.SetWriteTimeout(5 * time.Second)

    // 启动服务器
    errCh := make(chan error, 1)
    srv.Run(errCh)

    // 等待服务器退出
    if err := <-errCh; err != nil {
        panic(err)
    }
}

// 示例：自定义连接读取器实现
type MyConnReader struct {
    conn net.Conn
}

func (r *MyConnReader) ReadFrame() ([]byte, error) {
    // 示例：读取一行数据
    buf := make([]byte, 1024)
    n, err := r.conn.Read(buf)
    if err != nil {
        return nil, err
    }
    return buf[:n], nil
}
```

### 高级特性：灵活的连接IO工厂

框架支持通过ConnIOFactory工厂模式注入不同的连接实现，满足各种业务场景需求。所有连接实现都必须遵循 `IConnIO` 接口的线程安全约定。

#### 不同队列容量的连接实现

```go
// 大容量队列连接（适合高并发场景）
largeQueueFactory := func(conn net.Conn) stcp.IConnIO {
    readerFactory := func(c net.Conn) stcp.IConnReader {
        return &MyConnReader{conn: c}
    }
    return stcp.NewQSendConnHandler(conn, 10000, readerFactory)
}

// 小容量队列连接（适合内存敏感场景）
smallQueueFactory := func(conn net.Conn) stcp.IConnIO {
    readerFactory := func(c net.Conn) stcp.IConnReader {
        return &MyConnReader{conn: c}
    }
    return stcp.NewQSendConnHandler(conn, 100, readerFactory)
}

// 无限容量队列连接
unlimitedFactory := func(conn net.Conn) stcp.IConnIO {
    readerFactory := func(c net.Conn) stcp.IConnReader {
        return &MyConnReader{conn: c}
    }
    return stcp.NewQSendConnHandler(conn, 0, readerFactory)  // 0表示无限容量
}
```

#### 条件化连接选择

```go
// 根据连接来源选择不同的连接实现
conditionalFactory := func(conn net.Conn) stcp.IConnIO {
    remoteAddr := conn.RemoteAddr().String()
    
    readerFactory := func(c net.Conn) stcp.IConnReader {
        if strings.Contains(remoteAddr, "127.0.0.1") {
            return &LocalConnReader{conn: c}  // 本地连接使用特殊读取器
        }
        return &RemoteConnReader{conn: c}     // 远程连接使用标准读取器
    }
    
    if strings.Contains(remoteAddr, "127.0.0.1") {
        // 本地连接使用大队列
        return stcp.NewQSendConnHandler(conn, 10000, readerFactory)
    } else {
        // 外部连接使用小队列
        return stcp.NewQSendConnHandler(conn, 1000, readerFactory)
    }
}
```

#### 自定义连接实现

```go
// 如果你有自定义的IConnIO实现
customFactory := func(conn net.Conn) stcp.IConnIO {
    // return your custom implementation
    // return NewMyCustomConnIO(conn, customConfig)
    
    // 示例中仍使用QSendConn
    readerFactory := func(c net.Conn) stcp.IConnReader {
        return &MyConnReader{conn: c}
    }
    return stcp.NewQSendConnHandler(conn, 1024, readerFactory)
}
```

### 连接生命周期管理示例

#### 业务代码中主动关闭连接

```go
func businessLogic(iConnIO stcp.IConnIO, buffer []byte) error {
    // 检查是否需要关闭连接
    if shouldCloseConnection(buffer) {
        // 线程安全地关闭连接，触发完整的清理流程
        if err := iConnIO.Close(); err != nil {
            // 记录错误但不会 panic
            log.Printf("Close connection error: %v", err)
        }
        return io.EOF // 返回错误让 ReadProcessor 退出
    }
    
    // 继续处理消息
    return iConnIO.Put2Queue(processMessage(buffer))
}
```

#### 多goroutine环境下的安全调用

```go
func multiGoroutineExample(iConnIO stcp.IConnIO) {
    // Goroutine 1: 处理业务逻辑
    go func() {
        for {
            // 线程安全的元信息更新
            iConnIO.SetMetaInfo(&MyMetaInfo{
                UserID:    getCurrentUser(),
                Timestamp: time.Now().Unix(),
            })
            time.Sleep(time.Second)
        }
    }()
    
    // Goroutine 2: 发送心跳
    go func() {
        ticker := time.NewTicker(30 * time.Second)
        defer ticker.Stop()
        
        for range ticker.C {
            // 线程安全的消息发送
            heartbeat := []byte("PING")
            if err := iConnIO.Put2Queue(heartbeat); err != nil {
                // 连接可能已关闭，安全退出
                return
            }
        }
    }()
    
    // Goroutine 3: 条件关闭
    go func() {
        <-shutdownSignal
        // 多个goroutine可以安全地调用Close()
        iConnIO.Close() // 重入安全，不会panic
    }()
}
```

### 定长消息读取示例

```go
import (
    "io"
    "net"
    
    "github.com/pinealctx/neptune/stcp"
)

// 定长消息读取器实现
type FixedLengthReader struct {
    conn net.Conn
}

func (r *FixedLengthReader) ReadFrame() ([]byte, error) {
    const messageSize = 128
    buf := make([]byte, messageSize)
    if _, err := io.ReadFull(r.conn, buf); err != nil {
        return nil, err // 读不足或连接关闭则退出
    }
    return buf, nil
}

// ReadProcessor 处理定长消息
func fixedLengthProcessor(iConnIO stcp.IConnIO, buffer []byte) error {
    // 处理定长消息...
    return iConnIO.Put2Queue(buffer)
}
```

### 变长消息读取示例（前置长度头）

```go
import (
    "encoding/binary"
    "io"
    "net"
    
    "github.com/pinealctx/neptune/stcp"
)

// 变长消息读取器实现
type VarLengthReader struct {
    conn net.Conn
}

func (r *VarLengthReader) ReadFrame() ([]byte, error) {
    // 读取 4 字节长度头（大端序）
    var hdr [4]byte
    if _, err := io.ReadFull(r.conn, hdr[:]); err != nil {
        return nil, err
    }
    
    length := binary.BigEndian.Uint32(hdr[:])
    if length == 0 || length > 10<<20 { // 防御：限制最大包 10MB
        return nil, io.ErrUnexpectedEOF
    }
    
    // 读取消息体
    body := make([]byte, length)
    if _, err := io.ReadFull(r.conn, body); err != nil {
        return nil, err
    }
    
    return body, nil
}

// ReadProcessor 处理变长消息
func varLengthProcessor(iConnIO stcp.IConnIO, buffer []byte) error {
    // 处理变长消息...
    return iConnIO.Put2Queue(buffer)
}
```

### 心跳与读超时控制示例

```go
import (
    "io"
    "net"
    "time"
    
    "github.com/pinealctx/neptune/stcp"
)

// 心跳读取器实现
type HeartbeatReader struct {
    conn net.Conn
}

func (r *HeartbeatReader) ReadFrame() ([]byte, error) {
    const heartbeatInterval = 15 * time.Second
    const gracePeriod = 5 * time.Second
    
    // 设置读超时
    if err := r.conn.SetReadDeadline(time.Now().Add(heartbeatInterval + gracePeriod)); err != nil {
        return nil, err
    }

    // 读取消息
    var hdr [4]byte
    if _, err := io.ReadFull(r.conn, hdr[:]); err != nil {
        return nil, err
    }
    
    // 成功读到数据后，重置下一次的读超时
    if err := r.conn.SetReadDeadline(time.Now().Add(heartbeatInterval + gracePeriod)); err != nil {
        return nil, err
    }
    
    // 继续读取消息体...
    return hdr[:], nil
}

// ReadProcessor 处理心跳消息
func heartbeatProcessor(iConnIO stcp.IConnIO, buffer []byte) error {
    // 处理心跳或业务消息...
    return nil
}
```

### 自定义日志元信息

在握手成功或识别到业务身份后，可设置更丰富的 MetaInfo，便于日志检索：

```go
import (
    "github.com/pinealctx/neptune/stcp"
    "go.uber.org/zap/zapcore"
)

// 自定义MetaInfo结构
type MyMetaInfo struct {
    UserID     string
    RemoteAddr string
    SessionID  string
}

func (m *MyMetaInfo) MarshalLogObject(enc zapcore.ObjectEncoder) error {
    enc.AddString("userId", m.UserID)
    enc.AddString("remoteAddr", m.RemoteAddr)
    enc.AddString("sessionId", m.SessionID)
    return nil
}

func (m *MyMetaInfo) GetRemoteAddr() string {
    return m.RemoteAddr
}

// ReadProcessor 中设置自定义MetaInfo
func customMetaProcessor(iConnIO stcp.IConnIO, buffer []byte) error {
    // 例如：在握手后设置自定义MetaInfo
    customMeta := &MyMetaInfo{
        UserID:     "user123",
        RemoteAddr: iConnIO.Conn().RemoteAddr().String(),
        SessionID:  "session456",
    }
    iConnIO.SetMetaInfo(customMeta)
    
    // 继续处理消息...
    return iConnIO.Put2Queue(buffer)
}
```

### 优雅关闭示例

```go
import (
    "context"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    // 创建连接IO工厂
    connIOFactory := func(conn net.Conn) stcp.IConnIO {
        readerFactory := func(c net.Conn) stcp.IConnReader {
            return &MyConnReader{conn: c}
        }
        return stcp.NewQSendConnHandler(conn, 1024, readerFactory)
    }
    
    srv := stcp.NewTcpServer(cnf, readProcessor, connIOFactory)
    
    errCh := make(chan error, 1)
    srv.Run(errCh)
    
    // 监听系统信号
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    
    select {
    case err := <-errCh:
        if err != nil {
            panic(err)
        }
    case <-sigCh:
        // 优雅关闭
        if err := srv.Close(); err != nil {
            log.Printf("Server close error: %v", err)
        }
        // 等待现有连接处理完成
        time.Sleep(time.Second)
    }
}
```

---

## 运行时行为说明

### 连接处理流程
1. 服务器在指定地址监听TCP连接
2. 每个新连接会通过ConnIOFactory创建对应的IConnIO实例
3. ConnHandler启动两个goroutine管理连接：
   - 接收goroutine：循环调用IConnIO.ReadFrame()获取数据，然后调用ReadProcessor处理
   - 发送goroutine：从发送队列取数据并写入连接
4. 当ReadFrame()或ReadProcessor返回错误时，连接被关闭并触发退出钩子
5. ConnHandler提供panic恢复机制，确保单个连接异常不影响整体服务

### 连接数限制
- 达到MaxConn限制时，新连接会被立即关闭并记录日志
- 连接计数通过原子操作维护，在连接建立时递增，退出时递减

### 错误处理与重试
- Accept操作出错时采用指数退避策略重试
- 可配置最大重试次数和延迟时间
- 写操作支持超时控制，默认5秒
- Hook函数异常会被捕获并记录，不会影响连接正常运行

### 内存与性能
- 发送队列支持容量限制，防止内存无限增长
- 队列满时Put2Queue会返回错误，应用层应优雅处理
- 所有关键路径都经过并发安全设计
- `IConnIO` 接口方法的线程安全实现确保高并发场景下的稳定性
- `Close()` 方法使用 `sync.Once` 实现，避免重复资源释放的开销

### 设计原则与最佳实践

#### 接口分层原则
- **IConnSender**: 专注于发送能力，定义消息发送的抽象
- **IConnReader**: 专注于读取能力，支持不同的帧解析策略
- **IConnIO**: 组合两者，提供完整的连接操作能力
- 业务代码只需依赖相应接口，无需了解具体实现细节

#### 工厂模式设计
- **ConnReaderFactory**: 支持不同的帧读取策略（定长、变长、分隔符等）
- **ConnIOFactory**: 支持不同的连接实现（队列、直接发送、批处理等）
- 便于依赖注入和单元测试Mock，提高代码可测试性

#### 并发安全设计
- 所有公开方法都是线程安全的，可在多 goroutine 环境下安全使用
- `Close()` 方法支持重入调用，多次调用不会导致 panic
- 资源清理操作具有幂等性，确保系统的健壮性
- panic恢复机制确保单个连接异常不影响整体服务

#### 错误处理策略
- 方法可能返回错误，但保证不会 panic
- 错误信息结构化，便于日志记录和问题诊断
- 支持优雅降级，连接异常时不影响整体服务稳定性
- Hook函数异常会被隔离，不会传播到核心逻辑

---

## 架构特点

- **接口分层设计**：IConnSender/IConnReader/IConnIO三层接口实现关注点分离
- **工厂模式**：ConnReaderFactory和ConnIOFactory支持灵活的实现策略
- **接口驱动**：所有核心功能基于接口定义，便于扩展和测试
- **插拔式架构**：可以轻松替换和扩展连接读取器、发送器实现
- **异步发送**：队列化发送避免阻塞接收处理
- **完善日志**：结构化日志支持，MetaInfo可自定义
- **Hook机制**：连接生命周期钩子便于监控和扩展，具备panic恢复能力
- **并发安全**：所有共享状态都有适当的同步保护，支持高并发场景
- **重入安全**：关键方法如 Close() 支持多次调用，确保资源清理的可靠性
- **测试友好**：工厂模式便于单元测试时注入Mock实现
- **关注点分离**：读取、发送、连接管理职责清晰，便于维护和扩展
- **生命周期管理**：通过接口方法即可完整控制连接的创建、运行和销毁
- **错误隔离**：panic恢复机制确保单个连接异常不影响整体服务稳定性

如需更多示例或适配特定协议，可通过实现IConnReader接口按需定制帧读取逻辑。

---

## IConnIO 接口实现指南

如果需要实现自定义的 `IConnIO`，请遵循以下约定：

### 接口定义
```go
type IConnIO interface {
    IConnSender    // 继承发送能力
    IConnReader    // 继承读取能力
}

type IConnSender interface {
    // 必需方法 - 必须实现
    Conn() net.Conn              // 返回底层连接，线程安全
    SetMetaInfo(m MetaInfo)       // 设置元信息，线程安全，重入安全
    MetaInfo() MetaInfo          // 获取元信息，线程安全
    Close() error                // 关闭连接，线程安全，重入安全，不能panic
    
    // 可选方法 - 至少实现一个Put2*方法
    Put2Queue(bs []byte) error
    Put2SendMap(key uint32, bs []byte) error
    Put2SendSMap(key string, bs []byte) error
    Put2SendMaps(pairs []KeyIntBytesPair) error
    Put2SendSMaps(pairs []KeyStrBytesPair) error
    
    // 内部方法 - 框架调用，不对外公开
    loopSend()                   // 发送循环，仅由ConnHandler调用
}

type IConnReader interface {
    ReadFrame() ([]byte, error)  // 读取一帧数据，线程安全
}
```

### 实现要求

#### 1. 线程安全性
- 除 `loopSend()` 外的所有方法都必须是线程安全的
- `ReadFrame()` 通常在单一goroutine中调用，但也应考虑线程安全
- 可以使用 `sync.Mutex`, `atomic.Value`, `sync.Once` 等同步原语

#### 2. 重入安全性
- `Close()` 方法必须支持多次调用
- 推荐使用 `sync.Once` 确保资源只清理一次
- 多次调用应该是无害的，可以返回错误但不能panic

#### 3. 错误处理
- 方法可以返回错误，但绝对不能panic
- 错误信息应该具有描述性，便于调试
- ReadFrame()返回错误会导致连接关闭

#### 4. 资源管理
- 在 `Close()` 中确保所有资源得到正确清理
- 队列、连接、goroutine等都应该被适当关闭
- 考虑错误包装，提供清晰的错误上下文

### 实现示例模板
```go
type MyCustomConnIO struct {
    closeOnce sync.Once
    conn      net.Conn
    metaInfo  atomic.Value
    reader    IConnReader  // 组合读取器
    // 其他字段...
}

func (c *MyCustomConnIO) Close() error {
    var err error
    c.closeOnce.Do(func() {
        // 清理资源
        err = c.conn.Close()
        // 清理其他资源...
    })
    return err
}

func (c *MyCustomConnIO) SetMetaInfo(m MetaInfo) {
    c.metaInfo.Store(m)
}

func (c *MyCustomConnIO) MetaInfo() MetaInfo {
    if v := c.metaInfo.Load(); v != nil {
        return v.(MetaInfo)
    }
    return nil
}

func (c *MyCustomConnIO) ReadFrame() ([]byte, error) {
    // 可以委托给内部reader，或直接实现
    return c.reader.ReadFrame()
}

// 实现其他必需的方法...
```

### 最佳实践

#### 组合优于继承
- 推荐组合现有的IConnReader实现而不是从头实现
- 可以重用框架提供的基础组件，专注于业务逻辑

#### 错误包装
- 在ReadFrame()中进行错误包装，提供更多上下文：
```go
func (c *MyCustomConnIO) ReadFrame() ([]byte, error) {
    buf, err := c.reader.ReadFrame()
    if err != nil {
        return nil, fmt.Errorf("MyCustomConnIO.ReadFrame: %w", err)
    }
    return buf, nil
}
```

#### 工厂函数设计
- 提供便于使用的工厂函数：
```go
func NewMyCustomConnIO(conn net.Conn, readerFactory ConnReaderFactory) IConnIO {
    return &MyCustomConnIO{
        conn:   conn,
        reader: readerFactory(conn),
        // 初始化其他字段...
    }
}
```
