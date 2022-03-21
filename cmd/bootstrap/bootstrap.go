package bootstrap

import (
	"context"

	"github.com/patriciabonaldy/zero/config"
	"github.com/patriciabonaldy/zero/internal/platform/server"
	"github.com/patriciabonaldy/zero/internal/platform/storage/memory"
	"github.com/patriciabonaldy/zero/internal/platform/websocket/coinbase"
	"github.com/patriciabonaldy/zero/internal/trading"
)

// Run application
func Run() error {
	config, err := config.NewConfig()
	if err != nil {
		return err
	}

	client, err := coinbase.New(config.BrokerURL, config.Log)
	if err != nil {
		return err
	}

	repo := memory.NewRepository()
	s := trading.NewService(repo, client, config.Log, config.MaxSize)

	server := server.New(s)
	server.Run(context.Background(), config)

	return nil
}
