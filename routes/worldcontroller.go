package routes

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/larsmoa/renderdb/httpext"
)

type ctxWorldKeyType int

const (
	ctxWorldKey ctxWorldKeyType = 0
)

// WorldController handles requests to "/world".
type WorldController struct {
}

// NewWorldController creates a new controller for the "/world"-route.
func NewWorldController(router *mux.Router, db *sqlx.DB) *WorldController {
	controller := WorldController{}
	controller.init(router, db)
	return &controller
}

// Init initializes the routes for "/world".
func (c *WorldController) init(router *mux.Router, db *sqlx.DB) {
	renderer := httpext.NewJSONResponseRenderer()
	getWorlds := httpext.NewHttpHandler(db, renderer, &getWorldsHandler{})
	getWorld := httpext.NewHttpHandler(db, renderer, &getWorldHandler{})
	postWorld := httpext.NewHttpHandler(db, renderer, &postWorldHandler{})

	router.Handle("/worlds", getWorlds).Methods("GET")
	router.Handle("/worlds/{id:[0-9]+}", getWorld).Methods("GET")
	router.Handle("/worlds", postWorld).Methods("POST")
}
