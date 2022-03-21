package trading

import (
	"context"
	"testing"
	"time"

	"github.com/patriciabonaldy/zero/internal/model"

	"github.com/stretchr/testify/require"

	"github.com/patriciabonaldy/zero/internal/platform/logger"
	"github.com/patriciabonaldy/zero/internal/platform/storage"
	"github.com/patriciabonaldy/zero/internal/platform/storage/memory"
	"github.com/patriciabonaldy/zero/internal/platform/websocket/coinbase"
)

func Test_service_Trading(t *testing.T) {
	tests := []struct {
		name       string
		pairs      []string
		repository func() storage.Repository
		wantErr    bool
	}{
		{
			name:       "error in subscriber",
			pairs:      []string{"wsd-USD"},
			repository: memory.NewRepository,
			wantErr:    true,
		},
		{
			name:       "invalid coins",
			pairs:      []string{"QSP-USD"},
			repository: memory.NewRepository,
			wantErr:    true,
		},
		{
			name:  "error size of repo data is greater that default",
			pairs: []string{"ETH-USD"},
			repository: func() storage.Repository {
				repo := memory.NewRepository()
				data := []model.Data{
					{
						Price: 10,
						Size:  5,
					},
					{
						Price: 20,
						Size:  10,
					},
					{
						Price: 10,
						Size:  10,
					},
				}

				repo.ReplaceData(context.Background(), "ETH-USD", data)

				return repo
			},
			wantErr: true,
		},
		{
			name:       "success",
			pairs:      []string{"ETH-USD"},
			repository: memory.NewRepository,
		},
		{
			name:  "repo data is full",
			pairs: []string{"ETH-USD"},
			repository: func() storage.Repository {
				repo := memory.NewRepository()
				data := []model.Data{
					{
						Price: 10,
						Size:  5,
					},
					{
						Price: 20,
						Size:  10,
					},
				}

				repo.ReplaceData(context.Background(), "ETH-USD", data)
				for _, d := range data {
					repo.SaveVwpa(context.Background(), "ETH-USD", d)
				}

				return repo
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lg := logger.New()

			client, err := coinbase.New("wss://ws-feed.exchange.coinbase.com", lg)
			require.NoError(t, err)

			ctx := context.Background()
			repo := tt.repository()
			s := NewService(repo, client, lg, 2)
			s.Trading(ctx, tt.pairs)

			if tt.wantErr {
				err = <-s.ChErr
				if (err != nil) != tt.wantErr {
					t.Errorf("Trading() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				return
			}

			<-time.After(2 * time.Second)
		})
	}
}
