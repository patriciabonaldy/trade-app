package coinbase

import (
	"context"
	"encoding/json"
	"errors"

	ws "github.com/gorilla/websocket"
	"github.com/patriciabonaldy/zero/internal/platform/logger"
	"github.com/patriciabonaldy/zero/internal/platform/websocket"
)

type socket interface {
	WriteJSON(v interface{}) error
	Close() error
	ReadMessage() (messageType int, p []byte, err error)
}

type client struct {
	// The websocket connection.
	conn socket
	lg   logger.Logger
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
		lg:   lg,
	}, nil
}

func (c *client) Subscribe(ctx context.Context, tradingPairs []string, messages chan interface{}, err chan error) {
	go func() {
		request := Request{
			Type:       "subscribe",
			ProductIds: tradingPairs,
			Channels:   []string{"matches"},
		}
		cErr := c.conn.WriteJSON(request)
		if cErr != nil {
			c.lg.Error(cErr.Error())
			err <- cErr
			return
		}

		for {
			select {
			case <-ctx.Done():
				cErr = c.conn.Close()
				if cErr != nil {
					c.lg.Errorf("error closing ws connection: %s", err)
					err <- cErr
					return
				}
			default:
				var data []byte
				_, data, cErr = c.conn.ReadMessage()
				if cErr != nil {
					c.lg.Errorf("error receiving message: %s", err)
					err <- cErr
					break
				}

				var message Message
				cErr = json.Unmarshal(data, &message)
				if cErr != nil {
					c.lg.Errorf("error receiving message: %s", err)
					err <- cErr
					break
				}

				if message.Type == "error" {
					cErr = errors.New(message.Message)
					c.lg.Errorf("error receiving message %s", cErr)
					err <- cErr
					break
				}

				messages <- message
			}
		}
	}()
}
