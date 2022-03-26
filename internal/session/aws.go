package session

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/eventbridge"
)

func NewAPIGatewaySession(endpoint string) *apigatewaymanagementapi.ApiGatewayManagementApi {

	sess, _ := session.NewSession(&aws.Config{
		Endpoint: aws.String(endpoint),
	})
	return apigatewaymanagementapi.New(sess)
}

func NewDynamoDBSession() (*dynamodb.DynamoDB, error) {
	cfg := aws.NewConfig().WithLogLevel(aws.LogOff).WithRegion(os.Getenv("AWS_REGION"))
	sess, err := session.NewSession(cfg)
	return dynamodb.New(sess), err
}

func NewEventBridgeSession() (*eventbridge.EventBridge, error) {
	cfg := aws.NewConfig().WithLogLevel(aws.LogOff).WithRegion(os.Getenv("AWS_REGION"))
	sess, err := session.NewSession(cfg)

	eb := eventbridge.New(sess)

	return eb, err
}
