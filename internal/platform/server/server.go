package server

import (
	"context"

	"github.com/patriciabonaldy/zero/config"
	"github.com/patriciabonaldy/zero/internal/trading"
)

// Server represents a server
type Server struct {
	// deps
	service *trading.Service
}

// New create a new server
func New(tradingService *trading.Service) Server {
	srv := Server{
		service: tradingService,
	}

	return srv
}

// Run method start the server
func (s *Server) Run(ctx context.Context, config *config.Config) {
	s.service.Trading(ctx, config.TradingPairs)

	go func() {
		for {
			<-s.service.ChErr
		}
	}()

	<-ctx.Done()
}
