package options

import "github.com/stretchr/testify/mock"

import "github.com/ungerik/go3d/float64/vec3"

type MockFilterGeometryOption struct {
	mock.Mock
}

// Apply provides a mock function with given fields: bounds
func (_m *MockFilterGeometryOption) Apply(bounds []*vec3.Box) []int {
	ret := _m.Called(bounds)

	var r0 []int
	if rf, ok := ret.Get(0).(func([]*vec3.Box) []int); ok {
		r0 = rf(bounds)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]int)
		}
	}

	return r0
}
