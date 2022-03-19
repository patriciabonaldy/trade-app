package websocket

import "context"

// Client handles subscriptions to the broker
type Client interface {
	Subscribe(ctx context.Context, tradingPairs []string, message chan interface{}) error
}
