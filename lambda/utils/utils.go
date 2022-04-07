package utils

import (
	"os"

	"github.com/darren-reddick/go-apigw-webchat/internal/event"
	"github.com/darren-reddick/go-apigw-webchat/internal/store"
	"github.com/darren-reddick/go-apigw-webchat/internal/websocket"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func BuildApi() *websocket.ApigwWsApi {
	cfg := zap.NewProductionConfig()
	if os.Getenv("LOG_LEVEL") == "DEBUG" {
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		ec := cfg.EncoderConfig
		ec.TimeKey = "timestamp"
		ec.FunctionKey = "function"
		ec.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.EncoderConfig = ec
	}
	cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)

	logger, _ := cfg.Build()

	return websocket.NewApigwWsApi(
		store.NewConnectionStoreDynamo(os.Getenv("DYNAMO_DB_TABLE")),
		os.Getenv("WEBSOCKET_URL"),
		event.NewEventBridgeBus(os.Getenv("CHAT_EVENT_BUS")),
		logger,
	)
}
