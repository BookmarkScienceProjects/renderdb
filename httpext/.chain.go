package httpext

import (
	"fmt"
	"net/http"

	"github.com/justinas/alice"
)

// Chain creates a chain http.Handler. It supports the following argument types:
// - http.Handler: generic handler
// - func (r *http.Request) error: 'passive' observing handler
// Returns an alice.Chain.
func Chain(handlers ...interface{}) alice.Chain {
	if len(handlers) == 0 {
		panic("Must provide at least one handler")
	}

	middleware := make([]alice.Constructor, len(handlers))
	for i, m := range handlers {
		if c, ok := m.(alice.Constructor); ok {
			middleware[i] = c
		} else if mw, ok := m.(func(r *http.Request) error); ok {
			cmw := chainMiddleware(mw)
			middleware[i] = cmw.handle
		} else {
			panic(fmt.Sprintf("Middleware %d must be alice.Constructor or 'func(*http.Request) error' (was %T)", i, m))
		}
	}

	return alice.New(middleware...)
}

type chainMiddleware func(r *http.Request) error

func (m chainMiddleware) handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := m(r)
		if err != nil {
			WriteJSONError(w, NewHttpError(err, http.StatusInternalServerError))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
