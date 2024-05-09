// Code generated by mockery v2.36.1. DO NOT EDIT.

package gonvoy

import mock "github.com/stretchr/testify/mock"

// MockHttpFilterPhaseController is an autogenerated mock type for the HttpFilterPhaseController type
type MockHttpFilterPhaseController struct {
	mock.Mock
}

type MockHttpFilterPhaseController_Expecter struct {
	mock *mock.Mock
}

func (_m *MockHttpFilterPhaseController) EXPECT() *MockHttpFilterPhaseController_Expecter {
	return &MockHttpFilterPhaseController_Expecter{mock: &_m.Mock}
}

// Handle provides a mock function with given fields: c, proc
func (_m *MockHttpFilterPhaseController) Handle(c Context, proc HttpFilterProcessor) (HttpFilterAction, error) {
	ret := _m.Called(c, proc)

	var r0 HttpFilterAction
	var r1 error
	if rf, ok := ret.Get(0).(func(Context, HttpFilterProcessor) (HttpFilterAction, error)); ok {
		return rf(c, proc)
	}
	if rf, ok := ret.Get(0).(func(Context, HttpFilterProcessor) HttpFilterAction); ok {
		r0 = rf(c, proc)
	} else {
		r0 = ret.Get(0).(HttpFilterAction)
	}

	if rf, ok := ret.Get(1).(func(Context, HttpFilterProcessor) error); ok {
		r1 = rf(c, proc)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockHttpFilterPhaseController_Handle_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Handle'
type MockHttpFilterPhaseController_Handle_Call struct {
	*mock.Call
}

// Handle is a helper method to define mock.On call
//   - c Context
//   - proc HttpFilterProcessor
func (_e *MockHttpFilterPhaseController_Expecter) Handle(c interface{}, proc interface{}) *MockHttpFilterPhaseController_Handle_Call {
	return &MockHttpFilterPhaseController_Handle_Call{Call: _e.mock.On("Handle", c, proc)}
}

func (_c *MockHttpFilterPhaseController_Handle_Call) Run(run func(c Context, proc HttpFilterProcessor)) *MockHttpFilterPhaseController_Handle_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(Context), args[1].(HttpFilterProcessor))
	})
	return _c
}

func (_c *MockHttpFilterPhaseController_Handle_Call) Return(_a0 HttpFilterAction, _a1 error) *MockHttpFilterPhaseController_Handle_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockHttpFilterPhaseController_Handle_Call) RunAndReturn(run func(Context, HttpFilterProcessor) (HttpFilterAction, error)) *MockHttpFilterPhaseController_Handle_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockHttpFilterPhaseController creates a new instance of MockHttpFilterPhaseController. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockHttpFilterPhaseController(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockHttpFilterPhaseController {
	mock := &MockHttpFilterPhaseController{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
