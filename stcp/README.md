## TCP Server 框架（基于 IncomingHook）

该模块提供一个轻量的 TCP Server 框架：
- 通过 NewTCPSrv(address, incomingHook) 创建服务；address 为监听地址（如 :9000 或 127.0.0.1:9000）。
- IncomingHook 完全由业务方实现，负责从 net.Conn 读取数据并处理（定长、变长、心跳等策略均由外部控制）。
- 框架内部管理连接生命周期、并发与发送队列，暴露简单选项进行限流与退避配置。

---

- Server 构造
  - NewTCPSrv(address string, incomingHook IncomingHook) *Server
  - IncomingHook 定义：
    - type IncomingHook func(handler *ConnHandler, conn net.Conn, metaInfo MetaInfo) error
    - 框架会在每个连接上循环调用 IncomingHook；当返回非 nil 错误时，该连接将被关闭并触发退出钩子。

- 启动与关闭
  - 选项式启动：
    - RunWithOption(errChan chan<- error, opts ...Option)
    - 可选项：
      - WithMaxConn(int32): 最大并发连接数（默认 65535）
      - WithAccDelay(time.Duration): 初始 accept 退避
      - WithAccMaxDelay(time.Duration): 最大 accept 退避
      - WithAccMaxRetry(int): 最大 accept 重试次数
      - WithMaxSendQSize(int): 发送队列容量（默认 1024）
  - 配置式启动：RunWithCnf(errChan chan<- error, cnf *ServerAcceptCnf)
  - 关闭：Close() error（停止监听，已建立连接会按钩子流程退出）

- 发送与超时
  - 异步发送：handler.SendAsync([]byte) error（写入发送队列）
  - 写超时：SetWriteTimeout(d time.Duration)，默认 5s

- 日志与元信息
  - 默认 MetaInfo 为 BasicMetaInfo{RemoteAddr}
  - 可在握手阶段调用 handler.SetMetaInfo(...) 替换为自定义 MetaInfo（需实现 zapcore.ObjectMarshaler）以丰富日志字段。

---

快速上手示例

- 服务启动

```go
package main

import (
    "errors"
    "io"
    "net"
    "time"

    "github.com/pinealctx/neptune/stcp"
)

func main() {
    incoming := func(h *stcp.ConnHandler, c net.Conn, _ stcp.MetaInfo) error {
        // 这里根据你的协议读取并处理数据；返回非 nil 即关闭该连接
        // 示例：回显一行数据（按 \n 分割）
        buf := make([]byte, 1024)
        n, err := c.Read(buf)
        if err != nil {
            return err // io.EOF 或其他错误会触发退出
        }
        // 业务处理...
        _ = h.SendAsync(buf[:n])
        return nil
    }

    srv := stcp.NewTCPSrv(":9000", incoming)
    errCh := make(chan error, 1)
    srv.RunWithOption(errCh,
        stcp.WithMaxConn(1000),
        stcp.WithMaxSendQSize(1024),
    )

    // 可选：设置全局写超时
    stcp.SetWriteTimeout(5 * time.Second)

    if err := <-errCh; err != nil {
        panic(err)
    }
}
```

- 定长消息读取（例如每条 128 字节）

```go
import (
    "io"
    "net"

    "github.com/pinealctx/neptune/stcp"
)

func incomingFixed(h *stcp.ConnHandler, c net.Conn, _ stcp.MetaInfo) error {
    const size = 128
    buf := make([]byte, size)
    if _, err := io.ReadFull(c, buf); err != nil {
        return err // 读不足或连接关闭则退出
    }
    // 处理 buf...
    return nil
}
```

- 变长消息读取（前置 4 字节长度，大端）

```go
import (
    "encoding/binary"
    "io"
    "net"

    "github.com/pinealctx/neptune/stcp"
)

func incomingVarLen(h *stcp.ConnHandler, c net.Conn, _ stcp.MetaInfo) error {
    var hdr [4]byte
    if _, err := io.ReadFull(c, hdr[:]); err != nil {
        return err
    }
    n := binary.BigEndian.Uint32(hdr[:])
    if n == 0 || n > 10<<20 { // 简单防御：限制最大包 10MiB
        return io.ErrUnexpectedEOF
    }
    body := make([]byte, n)
    if _, err := io.ReadFull(c, body); err != nil {
        return err
    }
    // 处理 body...
    return nil
}
```

- 心跳与读超时控制

说明：框架不强制读超时；可在 IncomingHook 内按照协议心跳/最小报文间隔设置 ReadDeadline。

```go
import (
    "io"
    "net"
    "time"

    "github.com/pinealctx/neptune/stcp"
)

func incomingWithHeartbeat(h *stcp.ConnHandler, c net.Conn, _ stcp.MetaInfo) error {
    const heartbeat = 15 * time.Second      // 协议约定的最小消息间隔/心跳周期
    const grace = 5 * time.Second           // 容忍抖动
    if err := c.SetReadDeadline(time.Now().Add(heartbeat + grace)); err != nil {
        return err
    }

    // 如：先读 4 字节长度，再读内容
    var hdr [4]byte
    if _, err := io.ReadFull(c, hdr[:]); err != nil {
        return err
    }
    // 成功读到数据后，通常重置下一次的 ReadDeadline（可按需重复设置）
    if err := c.SetReadDeadline(time.Now().Add(heartbeat + grace)); err != nil {
        return err
    }
    // ...继续按协议读取与处理
    return nil
}
```

---

进阶：自定义日志元信息

- 在握手成功或识别到业务身份后，可设置更丰富的 MetaInfo，便于日志检索。

```go
import (
    "github.com/pinealctx/neptune/stcp"
    "go.uber.org/zap/zapcore"
)

type MyMeta struct {
    UserID    string
    Remote    string
}

func (m *MyMeta) MarshalLogObject(enc zapcore.ObjectEncoder) error {
    enc.AddString("userId", m.UserID)
    enc.AddString("remote", m.Remote)
    return nil
}

func incomingWithMeta(h *stcp.ConnHandler, c net.Conn, _ stcp.MetaInfo) error {
    // 例如：握手后获取 userId
    h.SetMetaInfo(&MyMeta{UserID: "u123", Remote: c.RemoteAddr().String()})
    // ...继续读取/处理
    return nil
}
```

---

运行时行为概览
- 每个连接一个接收循环（调用 IncomingHook）。当 IncomingHook 返回错误时，连接关闭并触发退出钩子，连接计数递减。
- 达到最大连接数 MaxConn 时，新连接会立即关闭并记录错误日志。
- accept 出错时按指数退避（受 AcceptDelay/AcceptMaxDelay/AcceptMaxRetry 控制）。
- 写操作在真实写入前设置 WriteDeadline（受 SetWriteTimeout 控制）。

如需更多示例或适配特定协议（定长/变长处理等），可在 IncomingHook 内按需实现。
