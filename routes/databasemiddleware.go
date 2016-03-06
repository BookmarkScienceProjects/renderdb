package routes

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/jmoiron/sqlx"
)

type ctxDbKeyType int
type ctxTxKeyType int

// ctxDbKey is used to store database instance (type sqlx.DB) in the context for
// the request.
const ctxDbKey ctxDbKeyType = 0

// ctxTxKey is used to store a database transaction (type sqlx.Tx) in the context
// for the request. The transaction is
const ctxTxKey ctxTxKeyType = 0

// DatabaseMiddleware injects a sqlx.DB-instance in the context (gorilla/context)
// Implements negroni.Handler
type DatabaseMiddleware struct {
	db *sqlx.DB
}

// ServeHTTP injects the database in the context (gorilla/context) and invokes the
// next handler. DbKey
func (m *DatabaseMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	context.Set(r, ctxDbKey, m.db)
	next.ServeHTTP(w, r)
}

// NewDatabaseMiddleware initializes a new middleware for injecting a database
// to the http request contexts.
func NewDatabaseMiddleware(db *sqlx.DB) *DatabaseMiddleware {
	return &DatabaseMiddleware{db}
}
