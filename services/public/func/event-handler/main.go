package main

import (
	"bytes"
	"errors"
	"net/http"
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

type EnvVars struct {
	ChannelSecret        string
	ChannelToken         string
	pushMessageLambdaArn string
}

func getEnvironmentVariables() (envVars *EnvVars, err error) {
	channelSecret, ok := os.LookupEnv("CHANNEL_SECRET")
	if !ok {
		return nil, errors.New("CHANNEL_SECRET is not set")
	}

	channelToken, ok := os.LookupEnv("CHANNEL_TOKEN")
	if !ok {
		return nil, errors.New("CHANNEL_TOKEN is not set")
	}

	pushMessageLambdaArn, ok := os.LookupEnv("PUSH_MESSAGE_LAMBDA_ARN")
	if !ok {
		return nil, errors.New("PUSH_MESSAGE_LAMBDA_ARN is not set")
	}
	return &EnvVars{
		ChannelSecret:        channelSecret,
		ChannelToken:         channelToken,
		pushMessageLambdaArn: pushMessageLambdaArn,
	}, nil

}

func Handler(request *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	envVars, err := getEnvironmentVariables()
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, nil
	}

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

	// 模擬一个http.Request，為了parse bot request的type
	reqBody := bytes.NewBufferString(request.Body)
	req, err := http.NewRequest("POST", "", reqBody)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, nil
	}

	// 複製headers，特别是 X-Line-Signature是必須的
	req.Header = make(http.Header)
	for key, value := range request.Headers {
		if len(value) > 0 {
			req.Header.Set(key, value)
		}
	}

	// 初始化 LINE Bot
	bot, err := linebot.New(
		envVars.ChannelSecret,
		envVars.ChannelToken,
	)

	if err != nil {
		logger.WithError(err).Error("Initial line bot failed")
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, nil
	}

	messageEvents, err := bot.ParseRequest(req)
	if err != nil {
		logger.WithError(err).Error("Parse request events failed")
		if err == linebot.ErrInvalidSignature {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
			}, nil
		} else {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
			}, nil
		}
	}

	logger.Infof("This invoke deals with %d messages", len(messageEvents))

	// 處理 LINE 事件
	for _, event := range messageEvents {
		logger.WithFields(logrus.Fields{
			"event_type": event.Type,
			"user_id":    event.Source.UserID,
			"room_id":    event.Source.RoomID,
			"group_id":   event.Source.GroupID,
		}).Info("event handling")
		if event.Type == linebot.EventTypeMessage {
			// 根據事件類型回應訊息
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
					logger.WithError(err)
				}
			}
		}
	}

	return events.APIGatewayProxyResponse{
		Body:       "ok",
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
