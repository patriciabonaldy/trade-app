package coinbase

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	ws "github.com/gorilla/websocket"
	"github.com/patriciabonaldy/zero/internal/platform/logger"
	"github.com/patriciabonaldy/zero/internal/platform/websocket"
)

const (
	typeSubscribe  = "subscribe"
	channelMatches = "matches"
)

type socket interface {
	WriteJSON(v interface{}) error
	Close() error
	ReadMessage() (messageType int, p []byte, err error)
}

type client struct {
	// The websocket connection.
	conn socket
}

// New function create a coinbase websocket client
func New(url string, lg logger.Logger) (websocket.Client, error) {
	conn, _, err := ws.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	lg.Info("connected to coinbase websocket, url: %s", url)

	return &client{
		conn: conn,
	}, nil
}

func (c *client) Subscribe(ctx context.Context, tradingPairs []string, messages chan interface{}, err chan error) {
	go func() {
		request := Request{
			Type:       typeSubscribe,
			ProductIds: tradingPairs,
			Channels:   []string{channelMatches},
		}
		cErr := c.conn.WriteJSON(request)
		if cErr != nil {
			err <- cErr
			return
		}

		for {
			select {
			case <-ctx.Done():
				cErr = c.conn.Close()
				if cErr != nil {
					err <- cErr
					return
				}
			default:
				var data []byte
				_, data, cErr = c.conn.ReadMessage()
				if cErr != nil {
					err <- cErr
					break
				}

				var message Message
				cErr = json.Unmarshal(data, &message)
				if cErr != nil {
					err <- cErr
					break
				}

				if message.Type == "error" {
					cErr = errors.New(message.Message)
					err <- cErr
					break
				}

				if strings.Contains(message.Type, "match") {
					messages <- message
				}
			}
		}
	}()
}
