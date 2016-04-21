package db

import "github.com/stretchr/testify/mock"

import "github.com/ungerik/go3d/float64/vec3"

type MockObject struct {
	mock.Mock
}

// ID provides a mock function with given fields:
func (_m *MockObject) ID() int64 {
	ret := _m.Called()

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// WorldID provides a mock function with given fields:
func (_m *MockObject) WorldID() int64 {
	ret := _m.Called()

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// LayerID provides a mock function with given fields:
func (_m *MockObject) LayerID() int64 {
	ret := _m.Called()

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// SceneID provides a mock function with given fields:
func (_m *MockObject) SceneID() int64 {
	ret := _m.Called()

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// Bounds provides a mock function with given fields:
func (_m *MockObject) Bounds() *vec3.Box {
	ret := _m.Called()

	var r0 *vec3.Box
	if rf, ok := ret.Get(0).(func() *vec3.Box); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*vec3.Box)
		}
	}

	return r0
}

// GeometryData provides a mock function with given fields:
func (_m *MockObject) GeometryData() []byte {
	ret := _m.Called()

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	return r0
}

// Metadata provides a mock function with given fields:
func (_m *MockObject) Metadata() interface{} {
	ret := _m.Called()

	var r0 interface{}
	if rf, ok := ret.Get(0).(func() interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}
