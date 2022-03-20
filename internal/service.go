package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/patriciabonaldy/zero/internal/model"
	"github.com/patriciabonaldy/zero/internal/platform/logger"
	"github.com/patriciabonaldy/zero/internal/platform/storage"
	"github.com/patriciabonaldy/zero/internal/platform/websocket"
	"github.com/patriciabonaldy/zero/internal/platform/websocket/coinbase"
)

// Service defines the contract for Trading events.
type Service struct {
	repository storage.Repository
	client     websocket.Client
	lg         logger.Logger
	maxSize    int
	data       chan []byte
	ChErr      chan error
}

// NewService create a new vwpa service
func NewService(repository storage.Repository, client websocket.Client, lg logger.Logger, maxSize int) *Service {
	return &Service{
		repository: repository,
		client:     client,
		lg:         lg,
		maxSize:    maxSize,
		data:       make(chan []byte),
		ChErr:      make(chan error),
	}
}

// Trading method generate a VWPA for a list of coins pair
func (s *Service) Trading(ctx context.Context, pairs []string) {
	message := make(chan interface{})
	cErr := make(chan error)

	s.client.Subscribe(ctx, pairs, message, s.ChErr)
	go func() {
		for {
			select {
			case err := <-cErr:
				s.lg.Errorf("error subscribe socket: %s", err)
				s.ChErr <- err

			case m := <-message:
				msg, ok := m.(coinbase.Message)
				if !ok {
					s.lg.Error("invalid Message type")
					continue
				}
				err := s.process(ctx, msg)
				if err != nil {
					s.lg.Errorf("error process data: %s", err)
					s.ChErr <- err
				}
			case result := <-s.data:
				err := s.showResult(result)
				if err != nil {
					s.lg.Errorf("error show data: %s", err)
					s.ChErr <- err
				}
			}
		}
	}()
}

func (s *Service) process(ctx context.Context, m coinbase.Message) error {
	if err := m.Validate(); err != nil {
		return err
	}

	switch m.ProductID {
	case model.BTCUSDPair,
		model.ETHUSDPair,
		model.ETHBTCPair:
		if err := s.processCoinsPair(ctx, m); err != nil {
			return err
		}
	default:
		return fmt.Errorf("%s: %s", model.ErrInvalidCoinsPair, m.ProductID)
	}

	return nil
}

func (s *Service) checkMaxSize(ctx context.Context, m coinbase.Message) error {
	code := m.ProductID
	data, _ := s.repository.GetData(ctx, code)

	if data == nil {
		return nil
	}

	if len(data) == s.maxSize {
		dataVwpa, err := s.repository.GetVwpa(ctx, code)
		if err != nil {
			return err
		}

		dataVwpa.Size -= data[0].Size
		dataVwpa.Price -= data[0].Price
		dataVwpa.CalculateVwpa()

		data = data[1:]
		s.repository.ReplaceData(ctx, code, data)
		s.repository.UpdateVwpa(ctx, code, dataVwpa)
	}

	if len(data) > s.maxSize {
		return model.ErrInvalidSize
	}

	return nil
}

func (s *Service) processCoinsPair(ctx context.Context, m coinbase.Message) error {
	err := s.checkMaxSize(ctx, m)
	if err != nil {
		return err
	}

	const bitSize = 64
	price, err := strconv.ParseFloat(m.Price, bitSize)
	if err != nil {
		return err
	}

	size, err := strconv.ParseFloat(m.Size, bitSize)
	if err != nil {
		return err
	}

	data := model.Data{
		Price: price,
		Size:  size,
	}

	s.repository.SaveData(ctx, m.ProductID, data)
	s.repository.SaveVwpa(ctx, m.ProductID, data)

	result, err := s.repository.GetMapVWpa(ctx)
	if err != nil {
		return err
	}

	go func() { s.data <- result }()

	return nil
}

func (s *Service) showResult(result []byte) error {
	data := make(map[string]model.VWpaData)
	err := json.Unmarshal(result, &data)
	if err != nil {
		return err
	}

	for k, v := range data {
		s.lg.Info("coins pair %s  VWPA: %v\n", k, v.Vwpa)
	}

	return nil
}
