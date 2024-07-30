// Code generated by mockery v2.41.0. DO NOT EDIT.

package mocks

import (
	context "context"
	http "net/http"

	mock "github.com/stretchr/testify/mock"

	ordinals "github.com/babylonchain/staking-api-service/internal/clients/ordinals"

	types "github.com/babylonchain/staking-api-service/internal/types"
)

// OrdinalsClientInterface is an autogenerated mock type for the OrdinalsClientInterface type
type OrdinalsClientInterface struct {
	mock.Mock
}

// FetchUTXOInfos provides a mock function with given fields: ctx, utxos
func (_m *OrdinalsClientInterface) FetchUTXOInfos(ctx context.Context, utxos []types.UTXOIdentifier) ([]ordinals.OrdinalsOutputResponse, *types.Error) {
	ret := _m.Called(ctx, utxos)

	if len(ret) == 0 {
		panic("no return value specified for FetchUTXOInfos")
	}

	var r0 []ordinals.OrdinalsOutputResponse
	var r1 *types.Error
	if rf, ok := ret.Get(0).(func(context.Context, []types.UTXOIdentifier) ([]ordinals.OrdinalsOutputResponse, *types.Error)); ok {
		return rf(ctx, utxos)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []types.UTXOIdentifier) []ordinals.OrdinalsOutputResponse); ok {
		r0 = rf(ctx, utxos)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]ordinals.OrdinalsOutputResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []types.UTXOIdentifier) *types.Error); ok {
		r1 = rf(ctx, utxos)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*types.Error)
		}
	}

	return r0, r1
}

// GetBaseURL provides a mock function with given fields:
func (_m *OrdinalsClientInterface) GetBaseURL() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetBaseURL")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetDefaultRequestTimeout provides a mock function with given fields:
func (_m *OrdinalsClientInterface) GetDefaultRequestTimeout() int {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetDefaultRequestTimeout")
	}

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// GetHttpClient provides a mock function with given fields:
func (_m *OrdinalsClientInterface) GetHttpClient() *http.Client {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetHttpClient")
	}

	var r0 *http.Client
	if rf, ok := ret.Get(0).(func() *http.Client); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Client)
		}
	}

	return r0
}

// NewOrdinalsClientInterface creates a new instance of OrdinalsClientInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOrdinalsClientInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *OrdinalsClientInterface {
	mock := &OrdinalsClientInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
