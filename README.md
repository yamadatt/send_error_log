## このリポジトリは

CloudWatchLogsにメッセージを書き込む。

以下のリポジトリで構築したCloudWatchLogsを想定。

https://github.com/yamadatt/poc-sqs

## 準備

```bash
go get github.com/aws/aws-sdk-go-v2
go get github.com/aws/aws-sdk-go-v2/config
go get github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs
go get "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
```

定数で定義している以下を自分の環境に合わせて書き換える。

```
const (
	logGroupName  = "/aws/poc-sqs/log-group"                     // ロググループ
	logStreamName = "test-stream"                                // ログストリーム
	logMessage    = "This is an ERROR log for testing purposes." // メッセージ
)
```


