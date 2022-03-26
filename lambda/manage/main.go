package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/darren-reddick/go-apigw-webchat/internal/event"
	"github.com/darren-reddick/go-apigw-webchat/internal/store"
	"github.com/darren-reddick/go-apigw-webchat/internal/websocket"
)

func main() {
	lambda.Start(HandleConnect)
}

var api = websocket.NewApigwWsApi(
	store.NewConnectionStoreDynamo(os.Getenv("DYNAMO_DB_TABLE")),
	os.Getenv("WEBSOCKET_URL"),
	event.NewEventBridgeBus(os.Getenv("CHAT_EVENT_BUS")),
)

func HandleConnect(ctx context.Context, event events.CloudWatchEvent) (string, error) {

	err := api.PurgeGone()

	if err != nil {
		fmt.Println(err)
		return "ERROR", err
	}

	return "SUCCESS", nil
}
