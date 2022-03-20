// Code generated by mockery v2.10.0. DO NOT EDIT.

package storagemocks

import (
	context "context"

	model "github.com/patriciabonaldy/zero/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// GetData provides a mock function with given fields: ctx, code
func (_m *Repository) GetData(ctx context.Context, code string) ([]model.Data, error) {
	ret := _m.Called(ctx, code)

	var r0 []model.Data
	if rf, ok := ret.Get(0).(func(context.Context, string) []model.Data); ok {
		r0 = rf(ctx, code)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Data)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, code)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMapVWpa provides a mock function with given fields: ctx
func (_m *Repository) GetMapVWpa(ctx context.Context) ([]byte, error) {
	ret := _m.Called(ctx)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context) []byte); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetVwpa provides a mock function with given fields: ctx, code
func (_m *Repository) GetVwpa(ctx context.Context, code string) (model.VWpaData, error) {
	ret := _m.Called(ctx, code)

	var r0 model.VWpaData
	if rf, ok := ret.Get(0).(func(context.Context, string) model.VWpaData); ok {
		r0 = rf(ctx, code)
	} else {
		r0 = ret.Get(0).(model.VWpaData)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, code)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReplaceData provides a mock function with given fields: ctx, code, data
func (_m *Repository) ReplaceData(ctx context.Context, code string, data []model.Data) {
	_m.Called(ctx, code, data)
}

// SaveData provides a mock function with given fields: ctx, code, data
func (_m *Repository) SaveData(ctx context.Context, code string, data model.Data) {
	_m.Called(ctx, code, data)
}

// SaveVwpa provides a mock function with given fields: ctx, code, data
func (_m *Repository) SaveVwpa(ctx context.Context, code string, data model.Data) {
	_m.Called(ctx, code, data)
}

// UpdateVwpa provides a mock function with given fields: ctx, code, data
func (_m *Repository) UpdateVwpa(ctx context.Context, code string, data model.VWpaData) {
	_m.Called(ctx, code, data)
}
