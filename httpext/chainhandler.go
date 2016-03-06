package httpext

import (
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
)

// ChainHandler extends Handler and is used by Chain()
// to allow for incrementally building chained handlers.
type ChainHandler interface {
	Handler
	// Then creates a new handler adds all handlers from this
	// chain and the handler provided.
	Then(h Handler) ChainHandler
}

// Chain creates a handler that executes the handlers in sequence.
// If a handler returns an error or writes any bytes to the
// http.ResponseWriter, the chain is broken and no more handlers
// are processed.
func Chain(handlers ...Handler) ChainHandler {
	for i, h := range handlers {
		if h == nil {
			panic(fmt.Errorf("Argument %d was nil", i))
		}
	}
	return &chainedHandler{handlers}
}

type chainedHandler struct {
	handlers []Handler
}

func (c *chainedHandler) Handle(
	tx *sqlx.Tx,
	renderer ResponseRenderer,
	w http.ResponseWriter,
	r *http.Request) error {

	writerProxy := &responseWriterProxy{w: w}
	for _, h := range c.handlers {
		err := h.Handle(tx, renderer, writerProxy, r)
		if err != nil {
			// Bail out because of error
			return err
		} else if writerProxy.hasWritten {
			// Bail out because middleware has written response
			return nil
		}
	}
	return nil
}

func (c *chainedHandler) Then(h Handler) ChainHandler {
	if h == nil {
		panic("Handler cannot be nil")
	}
	extended := new(chainedHandler)
	extended.handlers = append(c.handlers, h)
	return extended
}

type responseWriterProxy struct {
	w          http.ResponseWriter
	hasWritten bool
}

func (p *responseWriterProxy) Header() http.Header {
	return p.w.Header()
}

func (p *responseWriterProxy) WriteHeader(statusCode int) {
	p.w.WriteHeader(statusCode)
}

func (p *responseWriterProxy) Write(data []byte) (int, error) {
	p.hasWritten = true
	return p.w.Write(data)
}
