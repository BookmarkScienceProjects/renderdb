package db

import "github.com/stretchr/testify/mock"

type MockLayers struct {
	mock.Mock
}

// GetAll provides a mock function with given fields:
func (_m *MockLayers) GetAll() ([]*Layer, error) {
	ret := _m.Called()

	var r0 []*Layer
	if rf, ok := ret.Get(0).(func() []*Layer); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*Layer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: layerid
func (_m *MockLayers) Get(layerid int64) (*Layer, error) {
	ret := _m.Called(layerid)

	var r0 *Layer
	if rf, ok := ret.Get(0).(func(int64) *Layer); ok {
		r0 = rf(layerid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Layer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(layerid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Add provides a mock function with given fields: layer
func (_m *MockLayers) Add(layer *Layer) (int64, error) {
	ret := _m.Called(layer)

	var r0 int64
	if rf, ok := ret.Get(0).(func(*Layer) int64); ok {
		r0 = rf(layer)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*Layer) error); ok {
		r1 = rf(layer)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: layerid
func (_m *MockLayers) Delete(layerid int64) error {
	ret := _m.Called(layerid)

	var r0 error
	if rf, ok := ret.Get(0).(func(int64) error); ok {
		r0 = rf(layerid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
