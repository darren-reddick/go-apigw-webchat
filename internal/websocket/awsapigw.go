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
}

func (api *ApigwWsApi) Connect(id string) error {
	err := api.Store.Add(id)

	if err != nil {
		return err
	}

	detail := event.ChatEvent{
		ConnectionId: id,
		Message:      fmt.Sprintf("Welcome! - connected on %s", id),
	}

	bytes, err := json.Marshal(detail)
	err = api.Bus.Put(string(bytes))
	if err != nil {
		return err
	}

	return nil
}

func (api *ApigwWsApi) Disconnect(id string) error {
	err := api.Store.Remove(id)

	if err != nil {
		return err
	}
	return nil
}

func (api *ApigwWsApi) PurgeGone() error {
	fmt.Printf("Purging gone connections from endpoint %s", api.Endpoint)
	conns := api.Store.List()

	var getInput apigatewaymanagementapi.GetConnectionInput

	for _, val := range conns {
		getInput.ConnectionId = &val
		fmt.Printf("Checking connection %s", val)
		err := api.GetConnection(val)
		if err != nil {
			t, ok := err.(*goneError)
			fmt.Printf("%+v\n", t)
			if ok {
				api.Store.Remove(val)
			}
		}

	}
	return nil
}

func (api *ApigwWsApi) GetConnection(id string) error {
	var getInput apigatewaymanagementapi.GetConnectionInput
	getInput.ConnectionId = &id
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
	_, err := api.Session.PostToConnection(&postInput)
	return err
}

func (api *ApigwWsApi) BroadcastMessage(msg string) error {
	fmt.Println("Broadcasting")
	conns := api.Store.List()

	for _, id := range conns {
		err := api.SendMessage(id, msg)
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

func (api *ApigwWsApi) SendConnectionList(id string) error {
	fmt.Println("Fetching connection list")
	conns := api.Store.List()

	msg := fmt.Sprintf("connections: %s", strings.Join(conns, " "))

	fmt.Println(msg)

	err := api.SendMessage(id, msg)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func NewApigwWsApi(store store.ConnectionStore, endpoint string, bus event.Bus) *ApigwWsApi {
	return &ApigwWsApi{
		Store:    store,
		Endpoint: endpoint,
		Session:  session.NewAPIGatewaySession(endpoint),
		Bus:      bus,
	}
}
