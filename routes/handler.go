package routes

import (
	"net/http"

	"github.com/jmoiron/sqlx"
)

type Handler interface {
	Handle(tx *sqlx.Tx, r http.Request, renderer ResponseRenderer) error
}

type ResponseRenderer interface {
	WriteObject(statusCode int, val interface{}) error
	WriteError(err error)
}
