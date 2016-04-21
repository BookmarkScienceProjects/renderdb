package db

import "github.com/stretchr/testify/mock"

type MockScenes struct {
	mock.Mock
}

// GetAll provides a mock function with given fields:
func (_m *MockScenes) GetAll() ([]*Scene, error) {
	ret := _m.Called()

	var r0 []*Scene
	if rf, ok := ret.Get(0).(func() []*Scene); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*Scene)
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
func (_m *MockScenes) Get(id int64) (*Scene, error) {
	ret := _m.Called(id)

	var r0 *Scene
	if rf, ok := ret.Get(0).(func(int64) *Scene); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Scene)
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

// Add provides a mock function with given fields: scene
func (_m *MockScenes) Add(scene *Scene) (int64, error) {
	ret := _m.Called(scene)

	var r0 int64
	if rf, ok := ret.Get(0).(func(*Scene) int64); ok {
		r0 = rf(scene)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*Scene) error); ok {
		r1 = rf(scene)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: sceneid
func (_m *MockScenes) Delete(sceneid int64) error {
	ret := _m.Called(sceneid)

	var r0 error
	if rf, ok := ret.Get(0).(func(int64) error); ok {
		r0 = rf(sceneid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
