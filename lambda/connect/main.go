package main

import (
	"context"

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

func HandleConnect(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	msg := "OK"

	id := request.RequestContext.ConnectionID

	err := api.Connect(id)

	if err != nil {
		msg = err.Error()
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       msg,
	}, nil
}
