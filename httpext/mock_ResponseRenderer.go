package httpext

import "github.com/stretchr/testify/mock"

import "net/http"

type MockResponseRenderer struct {
	mock.Mock
}

// WriteEmpty provides a mock function with given fields: w, statusCode
func (_m *MockResponseRenderer) WriteEmpty(w http.ResponseWriter, statusCode int) {
	_m.Called(w, statusCode)
}

// WriteObject provides a mock function with given fields: w, statusCode, val
func (_m *MockResponseRenderer) WriteObject(w http.ResponseWriter, statusCode int, val interface{}) {
	_m.Called(w, statusCode, val)
}

// WriteError provides a mock function with given fields: w, err
func (_m *MockResponseRenderer) WriteError(w http.ResponseWriter, err error) {
	_m.Called(w, err)
}
