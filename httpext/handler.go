package httpext

import (
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type Handler interface {
	// Handle handles a HTTP request and returns an error if the operation fails. The implementor
	// is responsible for writing the error to the response.
	Handle(tx *sqlx.Tx, renderer ResponseRenderer, w http.ResponseWriter, r *http.Request) error
}

func NewHttpHandler(db *sqlx.DB, renderer ResponseRenderer, h Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Initialize transaction
		tx, err := db.Beginx()
		if err != nil {
			err = NewHttpError(fmt.Errorf("Could not open transaction (Reason: %v)", err), http.StatusInternalServerError)
			renderer.WriteError(w, err)
			return
		}

		// Commit or rollback transation after handler is done
		defer func() {
			if err == nil {
				if err = tx.Commit(); err != nil {
					renderer.WriteError(w, err)
				}
				return
			}
			renderer.WriteError(w, err)
			tx.Rollback()
		}()

		// Run handler
		err = h.Handle(tx, renderer, w, r)
	})
}
