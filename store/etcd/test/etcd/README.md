### 测试说明
- 编译etcdTestCli
```bash
./build.sh
```

- 下载etcd 
```bash
wget https://github.com/etcd-io/etcd/releases/download/v3.4.8/etcd-v3.4.8-darwin-amd64.zip
```

- 运行etcd
```bash
# 在etcd运行目录下新建目录存放etcd数据
mkdir -p data
# 运行etcd
./etcd --data-dir=./data
```

- 写入etcd数据
```bash
./etcdTestCli --url="localhost:2379" --root="/test" --key="x1" --content='{"a":"A1"}' --mode="update" --dev=-1
```

- 监控etcd数据
```bash
./etcdTestCli --url="localhost:2379" --root="/test" --key="x1" --content='{"a":"A"}' --mode="watchValue"
```

- 修改etcd数据
```bash
./etcdTestCli --url="localhost:2379" --root="/test" --key="x1" --content='{"a":"A2"}' --mode="update" --dev=-1
```
在监控窗口可以看到修改的数据变化