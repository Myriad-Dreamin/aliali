
# AliAli

配置文件:

### Usage

```yaml
version: aliyunpan/v1beta
ali-drive:
  refresh-token: 123
  drive-id: 1234
  root-path: 录播同步
  chunk-size: 10485760
```

以[BililiveRecorder][BililiveRecorder]为例，在配置文件中添加：

```json
{
  "global": {
    "WebHookUrlsV2": {
      "HasValue": true,
      "Value": "http://bilibili-notifier:10305/notifier/bilibili"
    }
  }
}
```

### Build

```shell
go build -o docker/bin/notifier ./cmd/notifier
```

### Package

+ [ali-drive](./pkg/ali-drive): 支持文件分块上传、哈希检测和查看的阿里云盘客户端
+ [ali-notifier](./pkg/ali-notifier): 提供Webhook和与[BililiveRecorder][BililiveRecorder]适配的Webhook
  + 调用后自动上传和清理本地存储
  + 基于[sqlite3](sqlite3)事务实现持久化和可靠事件处理
  + 基于Sha1 Hash检查保证上传数据可靠性
+ [suppress](./pkg/suppress): 实验性的Golang错误处理方案，方案讨论见个人博客

### Testing

+ Dispatcher、Database: mock-based testing
+ Dispatcher、IO: [fuzzing][gofuzz]

[BililiveRecorder]: https://github.com/Bililive/BililiveRecorder
[sqlite3]: https://github.com/mattn/go-sqlite3
[gofuzz]: https://go.dev/blog/fuzz-beta

### License

MIT
