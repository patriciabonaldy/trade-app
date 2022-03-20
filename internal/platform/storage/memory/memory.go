package memory

import (
	"context"
	"sync"

	"github.com/patriciabonaldy/zero/internal/model"
	"github.com/patriciabonaldy/zero/internal/platform/storage"
)

const defaultMaxSize = 200

// Memory is a memory Repository implementation.
type Memory struct {
	mux      sync.Mutex
	Data     map[string][]model.Data
	VwpaData map[string]model.VWpaData
	MaxSize  uint
}

// NewRepository initializes a memory implementation of storage.Repository.
func NewRepository() storage.Repository {
	return &Memory{
		Data:     make(map[string][]model.Data),
		VwpaData: make(map[string]model.VWpaData),
	}
}

// GetData implements the storage.Repository interface.
func (m *Memory) GetData(ctx context.Context, code string) ([]model.Data, error) {
	defer m.mux.Unlock()

	m.mux.Lock()
	data, ok := m.Data[code]
	if !ok {
		return nil, model.ErrCoinsNotFound
	}

	return data, nil
}

// SaveData implements the storage.Repository interface.
func (m *Memory) SaveData(ctx context.Context, code string, data model.Data) {
	defer m.mux.Unlock()

	m.mux.Lock()
	memoryData, ok := m.Data[code]
	if ok {
		memoryData = append(memoryData, data)
		m.Data[code] = memoryData
	} else {
		m.Data[code] = []model.Data{data}
	}
}

// GetVwpa implements the storage.Repository interface.
func (m *Memory) GetVwpa(ctx context.Context, code string) (model.VWpaData, error) {
	defer m.mux.Unlock()

	m.mux.Lock()
	vwpaData, ok := m.VwpaData[code]
	if !ok {
		return model.VWpaData{}, model.ErrCoinsNotFound
	}

	return vwpaData, nil
}

// SaveVwpa implements the storage.Repository interface.
func (m *Memory) SaveVwpa(ctx context.Context, code string, data model.VWpaData) {
	defer m.mux.Unlock()

	m.mux.Lock()
	wpa, ok := m.VwpaData[code]
	if ok {
		wpa.Price += data.Price
		wpa.Quantity += data.Quantity
		wpa.CalculateVwpa()
		m.VwpaData[code] = wpa
	} else {
		m.VwpaData[code] = data
	}
}
