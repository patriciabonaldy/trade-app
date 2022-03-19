package coinbase

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/patriciabonaldy/zero/internal/platform/logger"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	_, err := New("wss://ws-feed.exchange.coinbase.com", logger.New())
	assert.NoError(t, err)
}

type MockSocket struct {
	wantReadError    bool
	wantWriteError   bool
	wantCloseError   bool
	wantErrorReceive bool
}

func (m MockSocket) WriteJSON(v interface{}) error {
	if m.wantWriteError {
		return errors.New("unknown error")
	}

	return nil
}

func (m MockSocket) Close() error {
	if m.wantCloseError {
		return errors.New("unknown error")
	}

	return nil
}

func (m MockSocket) ReadMessage() (messageType int, p []byte, err error) {
	if m.wantReadError {
		return 0, nil, errors.New("unknown error")
	}

	if m.wantErrorReceive {
		return 0, []byte("{\"type\":\"error\", \"message\":\"unknown error\"}\n"), nil
	}

	return 0, []byte("{\"type\":\"subscriptions\",\"sequence\":0,\"product_id\":\"\",\"price\":\"\",\"open_24h\":\"\",\"volume_24h\":\"\",\"low_24h\":\"\",\"high_24h\":\"\",\"volume_30d\":\"\",\"best_bid\":\"\",\"best_ask\":\"\",\"side\":\"\",\"time\":\"0001-01-01T00:00:00Z\",\"trade_id\":0,\"last_size\":\"\"}\n"), nil
}

var _ socket = &MockSocket{}

func TestClient_Subscribe(t *testing.T) {
	cases := []struct {
		name       string
		pairs      []string
		mockSocket socket
		wantError  bool
		expected   Message
		err        error
	}{
		{
			name:       "error write json",
			pairs:      []string{"BTC-USD"},
			mockSocket: &MockSocket{wantWriteError: true},
			wantError:  true,
			err:        errors.New("unknown error"),
		},
		{
			name:       "error read json",
			pairs:      []string{"BTC-USD"},
			mockSocket: &MockSocket{wantReadError: true},
			wantError:  true,
			err:        errors.New("unknown error"),
		},
		{
			name:       "validPairs",
			pairs:      []string{"BTC-USD"},
			mockSocket: &MockSocket{},
			expected: Message{
				Type: "subscriptions",
				Time: time.Time{},
			},
		},
		{
			name:       "invalidPairs",
			pairs:      []string{"xxxx-USD"},
			mockSocket: &MockSocket{wantErrorReceive: true},
			wantError:  true,
			err:        errors.New("unknown error"),
		},
	}

	channel := make(chan interface{})
	err := make(chan error)
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			client := &client{
				conn: tt.mockSocket,
				lg:   logger.New(),
			}

			client.Subscribe(context.Background(), tt.pairs, channel, err)
			if tt.wantError {
				gotErr := <-err
				assert.Equal(t, tt.err, gotErr)
				return
			}

			message := <-channel
			m, ok := message.(Message)
			if !ok {
				log.Fatal("Expected response of Message type")
			}

			assert.Equal(t, tt.expected, m)
		})
	}
}
