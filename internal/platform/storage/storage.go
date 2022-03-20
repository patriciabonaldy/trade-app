package storage

import (
	"context"

	"github.com/patriciabonaldy/zero/internal/model"
)

// Repository defines the expected behaviour from a lana storage.
type Repository interface {
	GetData(ctx context.Context, code string) ([]model.Data, error)
	SaveData(ctx context.Context, code string, data model.Data)
	GetVwpa(ctx context.Context, code string) (model.VWpaData, error)
	SaveVwpa(ctx context.Context, code string, data model.VWpaData)
}

//go:generate mockery --case=snake --outpkg=storagemocks --output=storagemocks --name=Repository
