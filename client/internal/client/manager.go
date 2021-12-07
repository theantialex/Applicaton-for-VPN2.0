package client

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"

	"vpn2.0/app/client/internal/config"
	"vpn2.0/app/client/internal/logs"
)

const (
	IDUndefined = -1
)
type Manager struct {
	Config *config.Config
	ID     int
}

func SetUpClient() (*Manager, context.Context) {
	conf, err := config.New()
	if err != nil {
		panic(err)
	}

	logger := logs.BuildLogger(conf)
	ctx := ctxzap.ToContext(context.Background(), logger)
	logger.Info("client starting...")

	c := Manager{Config: conf}
	c.SetClientID(IDUndefined)

	return &c, ctx
}

func (c *Manager) SetClientID(id int) {
	c.ID = id
}