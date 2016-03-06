package httpext

import (
	"net/http"

	"github.com/jmoiron/sqlx"
)

// Chain creates a handler that executes the handlers in sequence.
// If a handler returns an error or writes any bytes to the
// http.ResponseWriter, the chain is broken and no more handlers
// are processed.
func Chain(handlers ...Handler) Handler {
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

type responseWriterProxy struct {
	w          http.ResponseWriter
	hasWritten bool
}

func (p *responseWriterProxy) Header() http.Header {
	return p.Header()
}

func (p *responseWriterProxy) WriteHeader(statusCode int) {
	p.w.WriteHeader(statusCode)
}

func (p *responseWriterProxy) Write(data []byte) (int, error) {
	if len(data) > 0 {
		p.hasWritten = true
		return p.w.Write(data)
	}
}
