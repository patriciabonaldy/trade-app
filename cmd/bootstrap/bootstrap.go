package bootstrap

import (
	"context"

	"github.com/patriciabonaldy/zero/config"
	"github.com/patriciabonaldy/zero/internal"
	"github.com/patriciabonaldy/zero/internal/model"
	"github.com/patriciabonaldy/zero/internal/platform/storage/memory"
	"github.com/patriciabonaldy/zero/internal/platform/websocket/coinbase"
)

// Run application
func Run() error {
	config, err := config.NewConfig()
	if err != nil {
		return err
	}

	client, err := coinbase.New(config.ExchangeURL, config.Log)
	if err != nil {
		return err
	}

	repo := memory.NewRepository()
	s := internal.NewService(repo, client, config.Log, config.MaxSize)
	ctx := context.Background()

	config.Log.Info(string(model.Header))
	s.Trading(ctx, config.TradingPairs)

	go func() {
		for {
			<-s.ChErr
		}
	}()

	<-ctx.Done()
	return nil
}
