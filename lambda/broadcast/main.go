package main

import (
	"context"
	"encoding/json"
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

type BroadcastRequest struct {
	Message string `json:"message,omitempty"`
}

func HandleConnect(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	var bc BroadcastRequest
	json.Unmarshal([]byte(request.Body), &bc)

	err := api.BroadcastMessage(bc.Message)

	if err != nil {
		fmt.Println(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "OK",
	}, nil
}
