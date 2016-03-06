package httpext

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockHandler struct {
	mock.Mock
}

const (
	any = mock.Anything
)

func (m *mockHandler) Handle(tx *sqlx.Tx, renderer ResponseRenderer, w http.ResponseWriter, r *http.Request) error {
	args := m.Called(tx, renderer, w, r)
	return args.Error(0)
}

func TestChainedHandler_Handle_FirstHandlerFails_BailsOutAndReturnsError(t *testing.T) {
	// Arrange
	h1 := new(mockHandler)
	any := mock.Anything
	h1.On("Handle", any, any, any, any).Return(errors.New("fail"))
	h2 := new(mockHandler)

	// Act
	ch := Chain(h1, h2)
	err := ch.Handle(nil, nil, nil, nil)

	// Assert
	assert.Error(t, err)
	h1.AssertCalled(t, "Handle", any, any, any, any)
	h2.AssertNumberOfCalls(t, "Handle", 0)
}

func TestChainedHandler_Handle_FirstHandlerWritesResponse_StopsProcessing(t *testing.T) {
	// Arrange
	h1 := new(mockHandler)
	h1.On("Handle", any, any, any, any).Run(
		func(args mock.Arguments) {
			w := args.Get(2).(http.ResponseWriter)
			w.Write([]byte{1, 2, 3})
		}).Return(nil)
	h2 := new(mockHandler)

	// Act
	w := httptest.NewRecorder()
	ch := Chain(h1, h2)
	err := ch.Handle(nil, nil, w, nil)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []byte{1, 2, 3}, w.Body.Bytes())
	h1.AssertCalled(t, "Handle", any, any, any, any)
	h2.AssertNotCalled(t, "Handle", any, any, any, any)
}

func TestChainedHandler_Handle_HandlersSucceed_ProcessesAllHandlers(t *testing.T) {
	// Arrange
	h1 := new(mockHandler)
	h1.On("Handle", any, any, any, any).Return(nil)
	h2 := new(mockHandler)
	h2.On("Handle", any, any, any, any).Return(nil)

	// Act
	ch := Chain(h1, h2)
	err := ch.Handle(nil, nil, nil, nil)

	// Assert
	assert.NoError(t, err)
	h1.AssertCalled(t, "Handle", any, any, any, any)
	h2.AssertCalled(t, "Handle", any, any, any, any)
}

func TestChain_NilHandler_Panics(t *testing.T) {
	assert.Panics(t, func() {
		Chain(nil)
	})
}

func TestChainedHandler_Then_NilHandler_Panics(t *testing.T) {
	// Arrange
	h1 := new(mockHandler)
	chain := Chain(h1)

	// Act & Assert
	assert.Panics(t, func() {
		chain.Then(nil)
	})
}

func TestChainedHandler_Then_ValidHandler_ReturnsExtendedHandler(t *testing.T) {
	// Arrange
	chain := Chain()
	h := new(mockHandler)
	h.On("Handle", any, any, any, any).Return(nil)

	// Act
	extended := chain.Then(h)
	extended.Handle(nil, nil, nil, nil)

	// Assert
	h.AssertCalled(t, "Handle", any, any, any, any)
}

func TestChainedHandler_Then_ValidHandler_DoesNotModifyOriginal(t *testing.T) {
	// Arrange
	chain := Chain()
	h := new(mockHandler)
	h.On("Handle", any, any, any, any).Return(nil)

	// Act
	chain.Then(h)
	chain.Handle(nil, nil, nil, nil)

	// Assert
	h.AssertNotCalled(t, "Handle", any, any, any, any)
}
