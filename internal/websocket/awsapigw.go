package websocket

import (
	"fmt"
	"strings"

	"encoding/json"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi/apigatewaymanagementapiiface"
	"github.com/darren-reddick/go-apigw-webchat/internal/event"
	"github.com/darren-reddick/go-apigw-webchat/internal/session"
	"github.com/darren-reddick/go-apigw-webchat/internal/store"
	"go.uber.org/zap"
)

type goneError struct {
}

func (e *goneError) Error() string {
	return "Connection gone"
}

type ApigwWsApi struct {
	Store    store.ConnectionStore                                   `json:"store,omitempty"`
	Endpoint string                                                  `json:"endpoint,omitempty"`
	Session  apigatewaymanagementapiiface.ApiGatewayManagementApiAPI `json:"session,omitempty"`
	Bus      event.Bus                                               `json:"bus,omitempty"`
	Logger   *zap.Logger                                             `json:"logger,omitempty"`
}

func (api *ApigwWsApi) Connect(id string) error {
	api.Logger.Debug(fmt.Sprintf("Adding connection %s to store", id))
	err := api.Store.Add(id)

	if err != nil {
		api.Logger.Error(err.Error())
		return err
	}

	detail := event.ChatEvent{
		ConnectionId: id,
		Message:      fmt.Sprintf("Welcome! - connected on %s", id),
	}

	bytes, err := json.Marshal(detail)
	if err != nil {
		api.Logger.Error(err.Error())
		return err
	}

	api.Logger.Debug("Adding message to event bus")

	err = api.Bus.Put(string(bytes))
	if err != nil {
		api.Logger.Error(err.Error())
		return err
	}

	return nil
}

func (api *ApigwWsApi) Disconnect(id string) error {
	api.Logger.Debug(fmt.Sprintf("Removing connection %s from store", id))

	err := api.Store.Remove(id)

	if err != nil {
		api.Logger.Error(err.Error())
		return err
	}
	return nil
}

func (api *ApigwWsApi) PurgeGone() error {
	api.Logger.Debug(fmt.Sprintf("Purging gone connections from endpoint %s", api.Endpoint))
	api.Logger.Debug("Fetching connections from store")
	conns := api.Store.List()

	var getInput apigatewaymanagementapi.GetConnectionInput

	for _, val := range conns {
		getInput.ConnectionId = &val
		api.Logger.Debug(fmt.Sprintf("Checking connection %s", val))
		err := api.GetConnection(val)
		if err != nil {
			t, ok := err.(*goneError)
			api.Logger.Error(t.Error())
			if ok {
				api.Logger.Debug(fmt.Sprintf("Removing connection %s from store", val))
				api.Store.Remove(val)
			}
		}

	}
	return nil
}

func (api *ApigwWsApi) GetConnection(id string) error {
	var getInput apigatewaymanagementapi.GetConnectionInput
	getInput.ConnectionId = &id
	api.Logger.Debug(fmt.Sprintf("Getting connection %+v", getInput))
	_, err := api.Session.GetConnection(&getInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case apigatewaymanagementapi.ErrCodeGoneException:
				return &goneError{}

			default:
				return fmt.Errorf("no error code found.\nDefault error: %v", aerr.Error())
			}
		}

	}
	return nil
}

func (api *ApigwWsApi) SendMessage(id string, msg string) error {
	var postInput apigatewaymanagementapi.PostToConnectionInput
	postInput.ConnectionId = &id
	postInput.Data = []byte(msg)

	api.Logger.Debug(fmt.Sprintf("Posting message to connection: %s", id))
	_, err := api.Session.PostToConnection(&postInput)
	return err
}

func (api *ApigwWsApi) BroadcastMessage(msg string) error {
	api.Logger.Debug("Broadcasting")
	api.Logger.Debug("Fetching connections from store")
	conns := api.Store.List()

	for _, id := range conns {
		err := api.SendMessage(id, msg)
		if err != nil {
			api.Logger.Error(err.Error())
		}
	}

	return nil
}

func (api *ApigwWsApi) SendConnectionList(id string) error {
	api.Logger.Debug("Fetching connections from store")
	conns := api.Store.List()

	msg := fmt.Sprintf("connections: %s", strings.Join(conns, " "))

	api.Logger.Debug(fmt.Sprintf("Sending list of connections to: %s", id))
	err := api.SendMessage(id, msg)
	if err != nil {
		api.Logger.Error(err.Error())
	}

	return nil
}

func NewApigwWsApi(store store.ConnectionStore, endpoint string, bus event.Bus, logger *zap.Logger) *ApigwWsApi {
	return &ApigwWsApi{
		Store:    store,
		Endpoint: endpoint,
		Session:  session.NewAPIGatewaySession(endpoint),
		Bus:      bus,
		Logger:   logger,
	}
}
