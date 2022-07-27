## mpb

protobuf属于半自描述的序列化协议，收发两端(rpc/io)需要知晓消息的元信息。
在一些不依赖grpc的场景中使用protobuf时，如何识别消息就比较重要了，除非对应的场景只有一种消息格式。
比如消息队列中投递的格式是protobuf，如果一个主题上可以投递多种格式的消息，这时必须通过一种方式来重建消息格式。
这里使用的方式就是在传输/存储序列化后的protobuf信息时，在前面加上4个字节的标记。
这类标记最好使用代码生产工具来做，保证不会有重复的标记。

```go
// 使用示例
import github.com/pinealctx/neptune/mpb

// 注册消息标记
func init() {
	mpb.RegisterGenerator(func () proto.Message {return &YouDefinedMsg1{}})
	mpb.RegisterGenerator(func () proto.Message {return &YouDefinedMsg2{}})
	...
}

// 序列化消息
var data, err = mpb.MarshalMsg(youDefinedMsg)
// 序列化错误
var data, err = mpb.MarshalError(anErr)
// 反序列化
// 返回的message可能是*Status("google.golang.org/genproto/googleapis/rpc/status")
// 一种可以序列化的错误
var msg, err = mpb.UnmarshalMsg(data)
// 反序列化回报，一个典型的RPC可能返回消息，也可能返回错误
// msg -- RPC返回的消息
// msgErr -- RPC返回的错误
// err -- 反序列化本身的错误
var msg, msgErr, err = UnmarshalResponse(data)

// 如果有tag冲突的情况(正常情况不会有)，比如不同的包但消息命名完全一样。则可以使用不同的MsgPacker实例。
var msgPacker1 = mpb.NewMessagePacker()
var msgPacker2 = mpb.NewMessagePacker()

// 在msgPacker1和msgPacker2之间可以有冲突的消息标记
msgPacker1.RegisterGenerator(func () proto.Message {return &YouDefinedMsg1InPack1})
msgPacker2.RegisterGenerator(func () proto.Message {return &YouDefinedMsg1InPack2})
```
