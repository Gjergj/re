// Code generated by mockery v2.43.2. DO NOT EDIT.

package db

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockPersistence is an autogenerated mock type for the Persistence type
type MockPersistence struct {
	mock.Mock
}

type MockPersistence_Expecter struct {
	mock *mock.Mock
}

func (_m *MockPersistence) EXPECT() *MockPersistence_Expecter {
	return &MockPersistence_Expecter{mock: &_m.Mock}
}

// FetchProduct provides a mock function with given fields: ctx
func (_m *MockPersistence) FetchProduct(ctx context.Context) (Product, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for FetchProduct")
	}

	var r0 Product
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (Product, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) Product); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(Product)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockPersistence_FetchProduct_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FetchProduct'
type MockPersistence_FetchProduct_Call struct {
	*mock.Call
}

// FetchProduct is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockPersistence_Expecter) FetchProduct(ctx interface{}) *MockPersistence_FetchProduct_Call {
	return &MockPersistence_FetchProduct_Call{Call: _e.mock.On("FetchProduct", ctx)}
}

func (_c *MockPersistence_FetchProduct_Call) Run(run func(ctx context.Context)) *MockPersistence_FetchProduct_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockPersistence_FetchProduct_Call) Return(_a0 Product, _a1 error) *MockPersistence_FetchProduct_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockPersistence_FetchProduct_Call) RunAndReturn(run func(context.Context) (Product, error)) *MockPersistence_FetchProduct_Call {
	_c.Call.Return(run)
	return _c
}

// InsertProduct provides a mock function with given fields: ctx, p
func (_m *MockPersistence) InsertProduct(ctx context.Context, p Product) error {
	ret := _m.Called(ctx, p)

	if len(ret) == 0 {
		panic("no return value specified for InsertProduct")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, Product) error); ok {
		r0 = rf(ctx, p)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockPersistence_InsertProduct_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InsertProduct'
type MockPersistence_InsertProduct_Call struct {
	*mock.Call
}

// InsertProduct is a helper method to define mock.On call
//   - ctx context.Context
//   - p Product
func (_e *MockPersistence_Expecter) InsertProduct(ctx interface{}, p interface{}) *MockPersistence_InsertProduct_Call {
	return &MockPersistence_InsertProduct_Call{Call: _e.mock.On("InsertProduct", ctx, p)}
}

func (_c *MockPersistence_InsertProduct_Call) Run(run func(ctx context.Context, p Product)) *MockPersistence_InsertProduct_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(Product))
	})
	return _c
}

func (_c *MockPersistence_InsertProduct_Call) Return(_a0 error) *MockPersistence_InsertProduct_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockPersistence_InsertProduct_Call) RunAndReturn(run func(context.Context, Product) error) *MockPersistence_InsertProduct_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockPersistence creates a new instance of MockPersistence. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockPersistence(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockPersistence {
	mock := &MockPersistence{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
