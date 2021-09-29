## fjson

github.com/json-iterator/go是比标准json包更快的json序列化/反序列化包。

```go
//高效的json序列化函数
JSONFastMarshal

//高效的json反序列化函数
JSONFastUnmarshal

//高效的json序列化函数，并将json数据进行snappy压缩
JSONFastMarshalSnappy

//高效的json反序列化函数，可以反序列snappy压缩的json数据，也可以反序列化普通的json数据
JSONFastUnmarshalSnappy

//读取json文件并反序列化成对象或map
LoadJSONFile2Obj 
```
