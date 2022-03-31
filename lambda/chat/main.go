package main

import (
	"context"
	"encoding/json"
	"errors"

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

type ChatRequest struct {
	Message    string `json:"message,omitempty"`
	Connection string `json:"connection,omitempty"`
}

func HandleConnect(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	var c ChatRequest
	err := json.Unmarshal([]byte(request.Body), &c)

	if err != nil {
		api.Logger.Error(err.Error())
	}

	if c.Connection == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "No connection specified for message",
		}, errors.New("No connection specified for message")
	}

	err = api.SendMessage(c.Connection, c.Message)

	if err != nil {
		api.Logger.Error(err.Error())
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "OK",
	}, nil
}
