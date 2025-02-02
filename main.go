package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	logs "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

// 定数の定義
const (
	logGroupName  = "/aws/poc-sqs/log-group"                     // ロググループ
	logStreamName = "test-stream"                                // ログストリーム
	logMessage    = "This is an ERROR log for testing purposes." // メッセージ
)

func main() {
	// コンテキストの作成
	ctx := context.Background()

	// AWS SDK の設定をロード
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	// CloudWatch Logs クライアントの作成
	client := logs.NewFromConfig(cfg)

	// ログストリームの存在確認
	describeStreamsInput := &logs.DescribeLogStreamsInput{
		LogGroupName:        aws.String(logGroupName),
		LogStreamNamePrefix: aws.String(logStreamName),
	}

	describeStreamsOutput, err := client.DescribeLogStreams(ctx, describeStreamsInput)
	if err != nil {
		log.Fatalf("failed to describe log streams, %v", err)
	}

	var sequenceToken *string

	if len(describeStreamsOutput.LogStreams) == 0 {
		// ログストリームが存在しない場合は作成
		_, err := client.CreateLogStream(ctx, &logs.CreateLogStreamInput{
			LogGroupName:  aws.String(logGroupName),
			LogStreamName: aws.String(logStreamName),
		})
		if err != nil {
			log.Fatalf("failed to create log stream, %v", err)
		}
		fmt.Println("Log stream created.")
	} else {
		// 最後のログストリームのシーケンストークンを取得
		stream := describeStreamsOutput.LogStreams[0]
		sequenceToken = stream.UploadSequenceToken
	}

	// ログイベントの作成
	timestamp := time.Now().UnixNano() / int64(time.Millisecond) // ミリ秒単位のタイムスタンプ

	logEvent := types.InputLogEvent{
		Timestamp: aws.Int64(timestamp),
		Message:   aws.String(logMessage),
	}

	// PutLogEvents リクエストの準備
	putLogEventsInput := &logs.PutLogEventsInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String(logStreamName),
		LogEvents:     []types.InputLogEvent{logEvent},
	}

	if sequenceToken != nil {
		putLogEventsInput.SequenceToken = sequenceToken
	}

	// ログイベントの送信
	putLogEventsOutput, err := client.PutLogEvents(ctx, putLogEventsInput)
	if err != nil {
		log.Fatalf("failed to put log events, %v", err)
	}

	fmt.Println("Log event sent successfully.")
	fmt.Printf("Next sequence token: %s\n", aws.ToString(putLogEventsOutput.NextSequenceToken))
}
