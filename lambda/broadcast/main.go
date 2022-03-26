package main

import (
	"context"
	"encoding/json"
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
