package db

import "github.com/stretchr/testify/mock"

type MockObjects struct {
	mock.Mock
}

// Add provides a mock function with given fields: o
func (_m *MockObjects) Add(o Object) (int64, error) {
	ret := _m.Called(o)

	var r0 int64
	if rf, ok := ret.Get(0).(func(Object) int64); ok {
		r0 = rf(o)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(Object) error); ok {
		r1 = rf(o)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMany provides a mock function with given fields: ids
func (_m *MockObjects) GetMany(ids []int64) (<-chan Object, <-chan error) {
	ret := _m.Called(ids)

	var r0 <-chan Object
	if rf, ok := ret.Get(0).(func([]int64) <-chan Object); ok {
		r0 = rf(ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan Object)
		}
	}

	var r1 <-chan error
	if rf, ok := ret.Get(1).(func([]int64) <-chan error); ok {
		r1 = rf(ids)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(<-chan error)
		}
	}

	return r0, r1
}

// GetAll provides a mock function with given fields:
func (_m *MockObjects) GetAll() (<-chan Object, <-chan error) {
	ret := _m.Called()

	var r0 <-chan Object
	if rf, ok := ret.Get(0).(func() <-chan Object); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan Object)
		}
	}

	var r1 <-chan error
	if rf, ok := ret.Get(1).(func() <-chan error); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(<-chan error)
		}
	}

	return r0, r1
}
