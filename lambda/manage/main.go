package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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

	err := api.PurgeGone()

	if err != nil {
		fmt.Println(err)
		return "ERROR", err
	}

	return "SUCCESS", nil
}
