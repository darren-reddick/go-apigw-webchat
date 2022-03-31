package utils

import (
	"os"

	"github.com/darren-reddick/go-apigw-webchat/internal/event"
	"github.com/darren-reddick/go-apigw-webchat/internal/store"
	"github.com/darren-reddick/go-apigw-webchat/internal/websocket"
)

// func BuildApi() *zap.Logger {
//
// 	cfg := zap.NewProductionConfig()
//
// 	if os.Getenv("LOG_LEVEL") == "DEBUG" {
// 		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
// 	}
//
// 	return cfg.Build()
//
// }

func BuildApi() *websocket.ApigwWsApi {
	return websocket.NewApigwWsApi(
		store.NewConnectionStoreDynamo(os.Getenv("DYNAMO_DB_TABLE")),
		os.Getenv("WEBSOCKET_URL"),
		event.NewEventBridgeBus(os.Getenv("CHAT_EVENT_BUS")),
	)
}
