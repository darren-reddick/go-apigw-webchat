package websocket

type WsApi interface {
	Connect(id string) error
	Disconnect(id string) error
	PurgeGone() error
	GetConnection(id string) error
	SendMessage(id string, msg string) error
}
