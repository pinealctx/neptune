## TCP Server 框架（基于组合式设计）

该模块提供一个轻量、高性能的 TCP Server 框架：
- 通过 NewTcpServer(cnf, connReader, sendQSize) 创建服务
- ConnReader 完全由业务方实现，负责从 net.Conn 读取数据并处理（定长、变长、心跳等策略均由外部控制）
- 框架内部管理连接生命周期、并发处理与异步发送队列，支持完整的配置和Hook机制

### 核心组件

- **TcpServer**: 主服务器，负责Accept循环和连接管理
- **ConnHandler**: 连接处理器，管理单个连接的生命周期
- **QSender**: 队列发送器，提供异步消息发送能力
- **IConnSender**: 连接发送器接口，支持不同的发送策略

---

### API 接口

#### Server 构造
```go
func NewTcpServer(cnf *ServerAcceptCnf, connReader ConnReaderFunc, sendQSize int) *TcpServer
```

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
- 写超时控制：`SetWriteTimeout(d time.Duration)`，默认 5s

#### 连接信息
- 获取监听地址：`Address() string`
- 获取当前连接数：`ConnCount() int32`

#### 日志与元信息
- 默认 MetaInfo 为 BasicMetaInfo{RemoteAddr}
- 可调用 `handler.SetMetaInfo(...)` 替换为自定义 MetaInfo（需实现 zapcore.ObjectMarshaler）以丰富日志字段

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

    // 创建服务器配置
    cnf := stcp.DefaultServerAcceptCnf()
    cnf.Address = ":9000"
    cnf.MaxConn = 1000
    
    // 创建服务器（sendQSize=1024表示发送队列容量）
    srv := stcp.NewTcpServer(cnf, connReader, 1024)
    
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
    srv := stcp.NewTcpServer(cnf, connReader, 1024)
    
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

---

## 架构特点

- **组合式设计**：避免复杂继承，各组件职责清晰
- **接口驱动**：IConnSender接口支持不同发送策略扩展
- **异步发送**：队列化发送避免阻塞接收处理
- **完善日志**：结构化日志支持，MetaInfo可自定义
- **Hook机制**：连接生命周期钩子便于监控和扩展
- **并发安全**：所有共享状态都有适当的同步保护

如需更多示例或适配特定协议，可在ConnReader内按需实现相应的读取和处理逻辑。
