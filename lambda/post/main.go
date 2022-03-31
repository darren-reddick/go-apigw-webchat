package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	e "github.com/darren-reddick/go-apigw-webchat/internal/event"
	"github.com/darren-reddick/go-apigw-webchat/internal/websocket"
	"github.com/darren-reddick/go-apigw-webchat/lambda/utils"
)

func main() {
	lambda.Start(HandleConnect)
}

var api *websocket.ApigwWsApi

func init() {
	api = utils.BuildApi()
}

func HandleConnect(ctx context.Context, event events.CloudWatchEvent) (string, error) {

	detail := e.ChatEvent{}
	err := json.Unmarshal(event.Detail, &detail)

	if err != nil {
		api.Logger.Error(err.Error())
	}

	err = api.SendMessage(detail.ConnectionId, detail.Message)

	if err != nil {
		api.Logger.Error(err.Error())
		return "ERROR", err
	}

	return "SUCCESS", nil
}
