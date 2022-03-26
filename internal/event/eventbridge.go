package event

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/aws/aws-sdk-go/service/eventbridge/eventbridgeiface"
	"github.com/darren-reddick/go-apigw-webchat/internal/session"
)

type ChatEvent struct {
	ConnectionId string `json:"connection_id,omitempty"`
	Message      string `json:"message,omitempty"`
}

type EventBridgeBus struct {
	BusName string
	Session eventbridgeiface.EventBridgeAPI
}

func (e *EventBridgeBus) Put(i interface{}) error {

	entries := make([]*eventbridge.PutEventsRequestEntry, 0, 1)
	entries = append(entries, &eventbridge.PutEventsRequestEntry{
		Detail:       aws.String(i.(string)),
		DetailType:   aws.String("welcome"),
		EventBusName: &e.BusName,
		Source:       aws.String("welcome"),
	})

	events := eventbridge.PutEventsInput{Entries: entries}
	output, err := e.Session.PutEvents(&events)

	fmt.Printf("%+v\n", output)

	if err != nil {
		fmt.Printf("%s", err)
		return err
	}

	return nil
}

func NewEventBridgeBus(name string) *EventBridgeBus {
	sess, _ := session.NewEventBridgeSession()

	return &EventBridgeBus{
		BusName: name,
		Session: sess,
	}
}
