## TCP Server 框架（基于组合式设计）

该模块提供一个轻量、高性能的 TCP Server 框架：
- 通过 NewTcpServer(cnf, connReader, connSenderFactory) 创建服务，支持工厂模式注入不同的连接发送器
- ConnReader 完全由业务方实现，负责从 net.Conn 读取数据并处理（定长、变长、心跳等策略均由外部控制）
- 框架内部管理连接生命周期、并发处理，通过ConnSenderFactory支持灵活的发送策略，包括完整的配置和Hook机制

### 核心组件

- **TcpServer**: 主服务器，负责Accept循环和连接管理
- **ConnHandler**: 连接处理器，管理单个连接的生命周期
- **QSender**: 队列发送器，提供异步消息发送能力（IConnSender的一个实现）
- **IConnSender**: 连接发送器接口，支持不同的发送策略，所有方法（除内部loopSend）均为线程安全
- **ConnSenderFactory**: 连接发送器工厂，用于创建不同类型的发送器实例

---

### API 接口

#### Server 构造
```go
func NewTcpServer(cnf *ServerAcceptCnf, connReader ConnReaderFunc, connSenderFactory ConnSenderFactory) *TcpServer
```

#### 连接发送器工厂
```go
type ConnSenderFactory func(conn net.Conn) IConnSender
```
工厂函数用于为每个新连接创建相应的发送器实例，支持不同的发送策略。

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

#### 消息发送
- 异步发送：`handler.Put2Queue([]byte) error`（写入发送队列）
- 连接关闭：`handler.Close() error`（线程安全，支持重入调用）
- 写超时控制：`SetWriteTimeout(d time.Duration)`，默认 5s

#### 连接信息
- 获取监听地址：`Address() string`
- 获取当前连接数：`ConnCount() int32`

#### 日志与元信息
- 默认 MetaInfo 为 BasicMetaInfo{RemoteAddr}
- 可调用 `handler.SetMetaInfo(...)` 替换为自定义 MetaInfo（需实现 zapcore.ObjectMarshaler）以丰富日志字段
- SetMetaInfo 方法线程安全，支持在任何时候更新连接元信息

### IConnSender 接口设计

#### 线程安全与重入性
`IConnSender` 接口的所有公开方法都经过精心设计，确保：

- **线程安全**: 所有方法都可以在多个 goroutine 中并发调用
- **重入安全**: `Close()` 方法支持多次调用，不会 panic，确保资源清理的幂等性
- **错误处理**: 方法可能返回错误，但保证不会 panic

#### 方法分类
- **必需方法**: `Conn()`, `SetMetaInfo()`, `MetaInfo()`, `Close()`
- **可选方法**: `Put2Queue()`, `Put2SendMap()`, `Put2SendSMap()`, `Put2SendMaps()`, `Put2SendSMaps()`
- **内部方法**: `loopSend()` - 仅由框架内部调用，业务代码不应直接使用

#### 生命周期管理
通过 `IConnSender.Close()` 可以优雅地关闭连接：
```go
// 业务代码中的任何地方都可以安全调用
err := sender.Close()  // 线程安全，重入安全
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
    // ConnReader 函数：负责从连接读取数据并处理
    connReader := func(sender stcp.IConnSender, conn net.Conn) error {
        // 示例：回显一行数据
        buf := make([]byte, 1024)
        n, err := conn.Read(buf)
        if err != nil {
            return err // io.EOF 或其他错误会触发退出
        }
        // 业务处理...
        return sender.Put2Queue(buf[:n])  // 异步发送响应
    }

    // 创建连接发送器工厂（使用队列发送器）
    senderFactory := func(conn net.Conn) stcp.IConnSender {
        return stcp.NewQSendConnHandler(conn, 1024)  // 队列容量1024
    }

    // 创建服务器配置
    cnf := stcp.DefaultServerAcceptCnf()
    cnf.Address = ":9000"
    cnf.MaxConn = 1000
    
    // 创建服务器
    srv := stcp.NewTcpServer(cnf, connReader, senderFactory)
    
    // 设置连接钩子（可选）
    srv.SetStartHooker(func(sender stcp.IConnSender) {
        // 连接建立时的处理
    })
    srv.SetExitHooker(func(sender stcp.IConnSender) {
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
```

### 高级特性：灵活的连接发送器工厂

框架支持通过ConnSenderFactory工厂模式注入不同的发送器实现，满足各种业务场景需求。所有发送器实现都必须遵循 `IConnSender` 接口的线程安全约定。

#### 不同队列容量的发送器

```go
// 大容量队列发送器（适合高并发场景）
largeQueueFactory := func(conn net.Conn) stcp.IConnSender {
    return stcp.NewQSendConnHandler(conn, 10000)
}

// 小容量队列发送器（适合内存敏感场景）
smallQueueFactory := func(conn net.Conn) stcp.IConnSender {
    return stcp.NewQSendConnHandler(conn, 100)
}

// 无限容量队列发送器
unlimitedFactory := func(conn net.Conn) stcp.IConnSender {
    return stcp.NewQSendConnHandler(conn, 0)  // 0表示无限容量
}
```

#### 条件化发送器选择

```go
// 根据连接来源选择不同的发送器
conditionalFactory := func(conn net.Conn) stcp.IConnSender {
    remoteAddr := conn.RemoteAddr().String()
    
    if strings.Contains(remoteAddr, "127.0.0.1") {
        // 本地连接使用大队列
        return stcp.NewQSendConnHandler(conn, 10000)
    } else {
        // 外部连接使用小队列
        return stcp.NewQSendConnHandler(conn, 1000)
    }
}
```

#### 自定义发送器实现

```go
// 如果你有自定义的IConnSender实现
customFactory := func(conn net.Conn) stcp.IConnSender {
    // return your custom implementation
    // return NewMyCustomSender(conn, customConfig)
    return stcp.NewQSendConnHandler(conn, 1024)  // 示例中仍使用QSender
}
```

### 连接生命周期管理示例

#### 业务代码中主动关闭连接

```go
func businessLogic(sender stcp.IConnSender, conn net.Conn) error {
    // 正常的业务处理...
    buf := make([]byte, 1024)
    n, err := conn.Read(buf)
    if err != nil {
        return err
    }
    
    // 检查是否需要关闭连接
    if shouldCloseConnection(buf[:n]) {
        // 线程安全地关闭连接，触发完整的清理流程
        if err := sender.Close(); err != nil {
            // 记录错误但不会 panic
            log.Printf("Close connection error: %v", err)
        }
        return io.EOF // 返回错误让 ConnReader 退出
    }
    
    // 继续处理消息
    return sender.Put2Queue(processMessage(buf[:n]))
}
```

#### 多goroutine环境下的安全调用

```go
func multiGoroutineExample(sender stcp.IConnSender) {
    // Goroutine 1: 处理业务逻辑
    go func() {
        for {
            // 线程安全的元信息更新
            sender.SetMetaInfo(&MyMetaInfo{
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
            if err := sender.Put2Queue(heartbeat); err != nil {
                // 连接可能已关闭，安全退出
                return
            }
        }
    }()
    
    // Goroutine 3: 条件关闭
    go func() {
        <-shutdownSignal
        // 多个goroutine可以安全地调用Close()
        sender.Close() // 重入安全，不会panic
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

func fixedLengthReader(sender stcp.IConnSender, conn net.Conn) error {
    const messageSize = 128
    buf := make([]byte, messageSize)
    if _, err := io.ReadFull(conn, buf); err != nil {
        return err // 读不足或连接关闭则退出
    }
    // 处理定长消息...
    return sender.Put2Queue(buf)
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

func varLengthReader(sender stcp.IConnSender, conn net.Conn) error {
    // 读取 4 字节长度头（大端序）
    var hdr [4]byte
    if _, err := io.ReadFull(conn, hdr[:]); err != nil {
        return err
    }
    
    length := binary.BigEndian.Uint32(hdr[:])
    if length == 0 || length > 10<<20 { // 防御：限制最大包 10MB
        return io.ErrUnexpectedEOF
    }
    
    // 读取消息体
    body := make([]byte, length)
    if _, err := io.ReadFull(conn, body); err != nil {
        return err
    }
    
    // 处理变长消息...
    return sender.Put2Queue(body)
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

func heartbeatReader(sender stcp.IConnSender, conn net.Conn) error {
    const heartbeatInterval = 15 * time.Second
    const gracePeriod = 5 * time.Second
    
    // 设置读超时
    if err := conn.SetReadDeadline(time.Now().Add(heartbeatInterval + gracePeriod)); err != nil {
        return err
    }

    // 读取消息
    var hdr [4]byte
    if _, err := io.ReadFull(conn, hdr[:]); err != nil {
        return err
    }
    
    // 成功读到数据后，重置下一次的读超时
    if err := conn.SetReadDeadline(time.Now().Add(heartbeatInterval + gracePeriod)); err != nil {
        return err
    }
    
    // 继续处理消息...
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

func customMetaReader(sender stcp.IConnSender, conn net.Conn) error {
    // 例如：在握手后设置自定义MetaInfo
    customMeta := &MyMetaInfo{
        UserID:     "user123",
        RemoteAddr: conn.RemoteAddr().String(),
        SessionID:  "session456",
    }
    sender.SetMetaInfo(customMeta)
    
    // 继续处理消息...
    return nil
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
    // 创建连接发送器工厂
    senderFactory := func(conn net.Conn) stcp.IConnSender {
        return stcp.NewQSendConnHandler(conn, 1024)
    }
    srv := stcp.NewTcpServer(cnf, connReader, senderFactory)
    
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
2. 每个新连接会创建一个ConnHandler来管理
3. ConnHandler启动两个goroutine：
   - 接收goroutine：循环调用ConnReader处理入站数据
   - 发送goroutine：从发送队列取数据并写入连接
4. 当ConnReader返回错误时，连接被关闭并触发退出钩子

### 连接数限制
- 达到MaxConn限制时，新连接会被立即关闭并记录日志
- 连接计数通过原子操作维护，在连接建立时递增，退出时递减

### 错误处理与重试
- Accept操作出错时采用指数退避策略重试
- 可配置最大重试次数和延迟时间
- 写操作支持超时控制，默认5秒

### 内存与性能
- 发送队列支持容量限制，防止内存无限增长
- 队列满时Put2Queue会返回错误，连接处理器应优雅退出
- 所有关键路径都经过并发安全设计
- `IConnSender` 接口方法的线程安全实现确保高并发场景下的稳定性
- `Close()` 方法使用 `sync.Once` 实现，避免重复资源释放的开销

### 设计原则与最佳实践

#### 接口隔离原则
- 业务代码只需依赖 `IConnSender` 接口，无需了解 `ConnHandler` 实现细节
- 通过接口可以完整管理连接生命周期，包括发送消息和关闭连接
- 支持依赖注入和单元测试 Mock

#### 并发安全设计
- 所有公开方法都是线程安全的，可在多 goroutine 环境下安全使用
- `Close()` 方法支持重入调用，多次调用不会导致 panic
- 资源清理操作具有幂等性，确保系统的健壮性

#### 错误处理策略
- 方法可能返回错误，但保证不会 panic
- 错误信息结构化，便于日志记录和问题诊断
- 支持优雅降级，连接异常时不影响整体服务稳定性

---

## 架构特点

- **组合式设计**：避免复杂继承，各组件职责清晰
- **工厂模式**：ConnSenderFactory支持灵活的发送器创建策略
- **接口驱动**：IConnSender接口支持不同发送策略扩展，所有方法线程安全
- **插拔式架构**：可以轻松替换和扩展连接发送器实现
- **异步发送**：队列化发送避免阻塞接收处理
- **完善日志**：结构化日志支持，MetaInfo可自定义
- **Hook机制**：连接生命周期钩子便于监控和扩展
- **并发安全**：所有共享状态都有适当的同步保护，支持高并发场景
- **重入安全**：关键方法如 Close() 支持多次调用，确保资源清理的可靠性
- **测试友好**：工厂模式便于单元测试时注入Mock实现
- **接口隔离**：业务代码只依赖 IConnSender 接口，降低耦合度
- **生命周期管理**：通过接口方法即可完整控制连接的创建、运行和销毁

如需更多示例或适配特定协议，可在ConnReader内按需实现相应的读取和处理逻辑。

---

## IConnSender 接口实现指南

如果需要实现自定义的 `IConnSender`，请遵循以下约定：

### 必须实现的方法
```go
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
```

### 实现要求

#### 1. 线程安全性
- 除 `loopSend()` 外的所有方法都必须是线程安全的
- 可以使用 `sync.Mutex`, `atomic.Value`, `sync.Once` 等同步原语

#### 2. 重入安全性
- `Close()` 方法必须支持多次调用
- 推荐使用 `sync.Once` 确保资源只清理一次
- 多次调用应该是无害的，可以返回错误但不能panic

#### 3. 错误处理
- 方法可以返回错误，但绝对不能panic
- 错误信息应该具有描述性，便于调试

#### 4. 资源管理
- 在 `Close()` 中确保所有资源得到正确清理
- 队列、连接、goroutine等都应该被适当关闭

### 实现示例模板
```go
type MyCustomSender struct {
    closeOnce sync.Once
    conn      net.Conn
    metaInfo  atomic.Value
    // 其他字段...
}

func (s *MyCustomSender) Close() error {
    var err error
    s.closeOnce.Do(func() {
        // 清理资源
        err = s.conn.Close()
        // 清理其他资源...
    })
    return err
}

func (s *MyCustomSender) SetMetaInfo(m MetaInfo) {
    s.metaInfo.Store(m)
}

func (s *MyCustomSender) MetaInfo() MetaInfo {
    if v := s.metaInfo.Load(); v != nil {
        return v.(MetaInfo)
    }
    return nil
}
```
