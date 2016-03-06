package httpext

import (
	"net/http"

	"github.com/jmoiron/sqlx"
)

type Handler interface {
	Handle(tx *sqlx.Tx, renderer ResponseRenderer, w http.ResponseWriter, r *http.Request) error
}

func NewHttpHandler(db *sqlx.DB, renderer ResponseRenderer, h Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Initialize transaction
		tx, err := db.Beginx()
		if err != nil {
			renderer.WriteError(w, err)
			return
		}

		// Commit or rollback transation after handler is done
		defer func() {
			if err == nil {
				err = tx.Commit()
			} else {
				err = tx.Rollback()
			}
			if err != nil {
				renderer.WriteError(w, err)
			}
		}()

		// Run handler
		err = h.Handle(tx, renderer, w, r)
	})
}
