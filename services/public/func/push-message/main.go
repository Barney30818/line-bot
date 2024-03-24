package main

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/sirupsen/logrus"
)

const (
	SEVERITY  = "severity"
	MESSAGE   = "message"
	TIMESTAMP = "timestamp"
)

type RequestBody struct {
	Message string `json:"message"`
}

func Handler(request *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  TIMESTAMP,
			logrus.FieldKeyLevel: SEVERITY,
			logrus.FieldKeyMsg:   MESSAGE,
		},
	})
	// stdout and stderr are sent to AWS CloudWatch Logs
	logger.Infof("Processing Lambda request, id: %s", request.RequestContext.RequestID)

	var reqBody RequestBody
	err := json.Unmarshal([]byte(request.Body), &reqBody)
	if err != nil {
		logger.WithError(err).Error("Error parsing request body")
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Bad Request",
		}, nil
	}

	// 初始化 LINE Bot 客户端
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		logger.Errorf("Error initializing LINE bot client: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, nil
	}

	// 推播訊息給所有用戶，這裡使用了廣播訊息的方法，它不需要指定接收者的 ID
	// 注意：根據你的使用情況，這可能會觸發配額限制，請參考 LINE 官方文件
	if _, err := bot.BroadcastMessage(linebot.NewTextMessage(reqBody.Message)).Do(); err != nil {
		logger.WithError(err).Error("Error broadcasting message")
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to send message",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       "ok",
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
