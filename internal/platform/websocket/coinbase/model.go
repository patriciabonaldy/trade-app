package coinbase

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Request struct represents a Ticker channel websocket request
type Request struct {
	Type       string   `json:"type"`
	ProductIds []string `json:"product_ids"`
	Channels   []string `json:"channels"`
}

// Message struct represents a Ticker channel websocket response
type Message struct {
	Type      string    `json:"type"`
	Sequence  int       `json:"sequence"`
	ProductID string    `json:"product_id"`
	Price     string    `json:"price"`
	Open24H   string    `json:"open_24h"`
	Volume24H string    `json:"volume_24h"`
	Low24H    string    `json:"low_24h"`
	High24H   string    `json:"high_24h"`
	Volume30D string    `json:"volume_30d"`
	BestBid   string    `json:"best_bid"`
	BestAsk   string    `json:"best_ask"`
	Side      string    `json:"side"`
	Time      time.Time `json:"time"`
	TradeID   int       `json:"trade_id"`
	Size      string    `json:"size"`
	Message   string    `json:"message,omitempty"`
}

// Validate validates the closed struct.
func (m *Message) Validate() error {
	err := validation.ValidateStruct(m,
		validation.Field(&m.Type, validation.Required),
		validation.Field(&m.ProductID, validation.Required),
		validation.Field(&m.Price, validation.Required),
		validation.Field(&m.Size, validation.Required),
	)

	return err
}
