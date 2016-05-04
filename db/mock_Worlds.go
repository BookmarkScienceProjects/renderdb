package db

import "github.com/stretchr/testify/mock"

type MockWorlds struct {
	mock.Mock
}

// GetAll provides a mock function with given fields:
func (_m *MockWorlds) GetAll() ([]*World, error) {
	ret := _m.Called()

	var r0 []*World
	if rf, ok := ret.Get(0).(func() []*World); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*World)
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

// Get provides a mock function with given fields: id
func (_m *MockWorlds) Get(id int64) (*World, error) {
	ret := _m.Called(id)

	var r0 *World
	if rf, ok := ret.Get(0).(func(int64) *World); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*World)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Add provides a mock function with given fields: world
func (_m *MockWorlds) Add(world *World) (int64, error) {
	ret := _m.Called(world)

	var r0 int64
	if rf, ok := ret.Get(0).(func(*World) int64); ok {
		r0 = rf(world)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*World) error); ok {
		r1 = rf(world)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: worldid
func (_m *MockWorlds) Delete(worldid int64) error {
	ret := _m.Called(worldid)

	var r0 error
	if rf, ok := ret.Get(0).(func(int64) error); ok {
		r0 = rf(worldid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
