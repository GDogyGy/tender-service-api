// Code generated by mockery v2.49.1. DO NOT EDIT.

package tender_test

import (
	model "TenderServiceApi/internal/model"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Mockservice is an autogenerated mock type for the service type
type Mockservice struct {
	mock.Mock
}

// CheckResponsibleTender provides a mock function with given fields: ctx, username, tenderId
func (_m *Mockservice) CheckResponsibleTender(ctx context.Context, username string, tenderId string) (bool, error) {
	ret := _m.Called(ctx, username, tenderId)

	if len(ret) == 0 {
		panic("no return value specified for CheckResponsibleTender")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (bool, error)); ok {
		return rf(ctx, username, tenderId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) bool); ok {
		r0 = rf(ctx, username, tenderId)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, username, tenderId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateTender provides a mock function with given fields: ctx, saveModel
func (_m *Mockservice) CreateTender(ctx context.Context, saveModel model.Tender) (model.Tender, error) {
	ret := _m.Called(ctx, saveModel)

	if len(ret) == 0 {
		panic("no return value specified for CreateTender")
	}

	var r0 model.Tender
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, model.Tender) (model.Tender, error)); ok {
		return rf(ctx, saveModel)
	}
	if rf, ok := ret.Get(0).(func(context.Context, model.Tender) model.Tender); ok {
		r0 = rf(ctx, saveModel)
	} else {
		r0 = ret.Get(0).(model.Tender)
	}

	if rf, ok := ret.Get(1).(func(context.Context, model.Tender) error); ok {
		r1 = rf(ctx, saveModel)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FetchList provides a mock function with given fields: ctx, serviceType
func (_m *Mockservice) FetchList(ctx context.Context, serviceType string) ([]model.Tender, error) {
	ret := _m.Called(ctx, serviceType)

	if len(ret) == 0 {
		panic("no return value specified for FetchList")
	}

	var r0 []model.Tender
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]model.Tender, error)); ok {
		return rf(ctx, serviceType)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []model.Tender); ok {
		r0 = rf(ctx, serviceType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Tender)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, serviceType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FetchListByUser provides a mock function with given fields: ctx, username
func (_m *Mockservice) FetchListByUser(ctx context.Context, username string) ([]model.Tender, error) {
	ret := _m.Called(ctx, username)

	if len(ret) == 0 {
		panic("no return value specified for FetchListByUser")
	}

	var r0 []model.Tender
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]model.Tender, error)); ok {
		return rf(ctx, username)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []model.Tender); ok {
		r0 = rf(ctx, username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Tender)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FetchTenderById provides a mock function with given fields: ctx, tenderId
func (_m *Mockservice) FetchTenderById(ctx context.Context, tenderId string) (model.Tender, error) {
	ret := _m.Called(ctx, tenderId)

	if len(ret) == 0 {
		panic("no return value specified for FetchTenderById")
	}

	var r0 model.Tender
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (model.Tender, error)); ok {
		return rf(ctx, tenderId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) model.Tender); ok {
		r0 = rf(ctx, tenderId)
	} else {
		r0 = ret.Get(0).(model.Tender)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, tenderId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockservice creates a new instance of Mockservice. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockservice(t interface {
	mock.TestingT
	Cleanup(func())
}) *Mockservice {
	mock := &Mockservice{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
