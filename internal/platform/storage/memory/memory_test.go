package memory

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/patriciabonaldy/zero/internal/model"
)

func TestMemory_GetData(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		want    []model.Data
		wantErr bool
	}{
		{
			name:    "Not found",
			code:    "XXX-YYY",
			wantErr: true,
		},
		{
			name: "Found",
			code: "BTC-USD",
			want: []model.Data{
				{
					Price: 10,
					Size:  5,
				},
			},
		},
	}

	m := NewRepository()
	m.SaveData(context.Background(), "BTC-USD", model.Data{
		Price: 10,
		Size:  5,
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := m.GetData(context.Background(), tt.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVwpa() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMemory_GetVwpa(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		want    model.VWpaData
		wantErr bool
	}{
		{
			name:    "Not found",
			code:    "XXX-YYY",
			wantErr: true,
		},
		{
			name: "Found",
			code: "BTC-USD",
			want: model.VWpaData{
				PQ:   50,
				Size: 5,
				Vwpa: 10,
			},
		},
	}

	m := NewRepository()
	m.SaveVwpa(context.Background(), "BTC-USD", model.Data{
		Price: 10,
		Size:  5,
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := m.GetVwpa(context.Background(), tt.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVwpa() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetVwpa() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemory_SaveVwpa(t *testing.T) {
	data := model.Data{
		Price: 10,
		Size:  5,
	}

	m := NewRepository()
	m.SaveVwpa(context.Background(), "BTC-USD", data)
	m.SaveVwpa(context.Background(), "BTC-USD", data)

	got, err := m.GetVwpa(context.Background(), "BTC-USD")
	require.NoError(t, err)

	expected := model.VWpaData{
		PQ:   100,
		Size: 10,
		Vwpa: 10,
	}
	assert.Equal(t, expected, got)
}
