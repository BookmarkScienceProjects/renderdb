package routes

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/larsmoa/renderdb/httpext"
)

// RegisterWorldsRoutes registers handlers for the "/worlds"-route.
func RegisterWorldsRoutes(router *mux.Router, db *sqlx.DB) {
	renderer := httpext.NewJSONResponseRenderer()
	middleware := httpext.Chain(&worldsMiddleware{})
	getWorlds := httpext.NewHttpHandler(db, renderer, middleware.Then(&getWorldsHandler{}))
	getWorld := httpext.NewHttpHandler(db, renderer, middleware.Then(&getWorldHandler{}))
	postWorld := httpext.NewHttpHandler(db, renderer, middleware.Then(&postWorldHandler{}))

	router.Handle("/worlds", getWorlds).Methods("GET")
	router.Handle("/worlds/{worldID:[0-9]+}", getWorld).Methods("GET")
	router.Handle("/worlds", postWorld).Methods("POST")
}
