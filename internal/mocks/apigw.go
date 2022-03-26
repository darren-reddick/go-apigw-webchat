package mocks

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

type TestApiGatewayManagementApiAPI struct {
	Connections map[string]bool
}

func (t TestApiGatewayManagementApiAPI) GetConnection(in *apigatewaymanagementapi.GetConnectionInput) (*apigatewaymanagementapi.GetConnectionOutput, error) {
	if ok := t.Connections[*in.ConnectionId]; ok {
		return &apigatewaymanagementapi.GetConnectionOutput{
			ConnectedAt: &time.Time{},
			Identity: &apigatewaymanagementapi.Identity{
				SourceIp:  nil,
				UserAgent: nil,
			},
			LastActiveAt: &time.Time{},
		}, nil
	}
	return nil, &apigatewaymanagementapi.GoneException{}
}

func (TestApiGatewayManagementApiAPI) DeleteConnection(_ *apigatewaymanagementapi.DeleteConnectionInput) (*apigatewaymanagementapi.DeleteConnectionOutput, error) {
	panic("not implemented") // TODO: Implement
}

func (TestApiGatewayManagementApiAPI) DeleteConnectionWithContext(_ aws.Context, _ *apigatewaymanagementapi.DeleteConnectionInput, _ ...request.Option) (*apigatewaymanagementapi.DeleteConnectionOutput, error) {
	panic("not implemented") // TODO: Implement
}

func (TestApiGatewayManagementApiAPI) DeleteConnectionRequest(_ *apigatewaymanagementapi.DeleteConnectionInput) (*request.Request, *apigatewaymanagementapi.DeleteConnectionOutput) {
	panic("not implemented") // TODO: Implement
}

func (TestApiGatewayManagementApiAPI) GetConnectionWithContext(_ aws.Context, _ *apigatewaymanagementapi.GetConnectionInput, _ ...request.Option) (*apigatewaymanagementapi.GetConnectionOutput, error) {
	panic("not implemented") // TODO: Implement
}

func (TestApiGatewayManagementApiAPI) GetConnectionRequest(_ *apigatewaymanagementapi.GetConnectionInput) (*request.Request, *apigatewaymanagementapi.GetConnectionOutput) {
	panic("not implemented") // TODO: Implement
}

func (TestApiGatewayManagementApiAPI) PostToConnection(_ *apigatewaymanagementapi.PostToConnectionInput) (*apigatewaymanagementapi.PostToConnectionOutput, error) {
	panic("not implemented") // TODO: Implement
}

func (TestApiGatewayManagementApiAPI) PostToConnectionWithContext(_ aws.Context, _ *apigatewaymanagementapi.PostToConnectionInput, _ ...request.Option) (*apigatewaymanagementapi.PostToConnectionOutput, error) {
	panic("not implemented") // TODO: Implement
}

func (TestApiGatewayManagementApiAPI) PostToConnectionRequest(_ *apigatewaymanagementapi.PostToConnectionInput) (*request.Request, *apigatewaymanagementapi.PostToConnectionOutput) {
	panic("not implemented") // TODO: Implement
}
