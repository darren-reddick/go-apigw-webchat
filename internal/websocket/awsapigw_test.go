package websocket

import (
	"testing"

	"github.com/darren-reddick/go-apigw-webchat/internal/event"
	"github.com/darren-reddick/go-apigw-webchat/internal/mocks"
	"github.com/darren-reddick/go-apigw-webchat/internal/store"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	logger, _ = zap.NewProduction()
}

func EqualStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

type testStore struct {
	items map[string]bool
}

func (t *testStore) List() []string {
	keys := make([]string, 0, len(t.items))

	for k := range t.items {
		keys = append(keys, k)
	}
	return keys

}
func (t *testStore) Add(id string) error {
	t.items[id] = true
	return nil

}
func (t *testStore) Remove(id string) error {
	delete(t.items, id)
	return nil
}

type testBus struct {
}

func (t *testBus) Put(interface{}) error {
	return nil
}

func TestApigwWsApi_Connect(t *testing.T) {
	type fields struct {
		Store    store.ConnectionStore
		Endpoint string
		Session  mocks.TestApiGatewayManagementApiAPI
		Bus      event.Bus
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		expect []string
	}{
		{
			"Simple connect",
			fields{
				Store:    &testStore{items: map[string]bool{"id1": true}},
				Endpoint: "",
				Session:  mocks.TestApiGatewayManagementApiAPI{},
				Bus:      &testBus{},
			},
			args{id: "id2"},
			[]string{"id1", "id2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &ApigwWsApi{
				Store:    tt.fields.Store,
				Endpoint: tt.fields.Endpoint,
				Session:  tt.fields.Session,
				Bus:      tt.fields.Bus,
				Logger:   logger,
			}
			if err := api.Connect(tt.args.id); err != nil {
				t.Errorf("ApigwWsApi.Connect() error = %v", err)
			}
			items := tt.fields.Store.List()
			if !EqualStringSlice(items, tt.expect) {
				t.Errorf("Wanted %v but got %s", items, tt.expect)

			}
		})
	}
}

func TestApigwWsApi_GetConnection(t *testing.T) {
	type fields struct {
		Store    store.ConnectionStore
		Endpoint string
		Session  mocks.TestApiGatewayManagementApiAPI
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Connection gone",
			fields: fields{
				Store:    nil,
				Endpoint: "",
				Session: mocks.TestApiGatewayManagementApiAPI{
					Connections: map[string]bool{
						"id1": true,
						"id2": true,
					},
				},
			},
			args: args{
				id: "id3",
			},
			wantErr: true,
		},
		{
			name: "Connection present",
			fields: fields{
				Store:    nil,
				Endpoint: "",
				Session: mocks.TestApiGatewayManagementApiAPI{
					Connections: map[string]bool{
						"id1": true,
						"id2": true,
					},
				},
			},
			args: args{
				id: "id2",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &ApigwWsApi{
				Store:    tt.fields.Store,
				Endpoint: tt.fields.Endpoint,
				Session:  tt.fields.Session,
				Logger:   logger,
			}
			if err := api.GetConnection(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("ApigwWsApi.GetConnection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
