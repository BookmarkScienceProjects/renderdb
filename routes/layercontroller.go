package routes

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/larsmoa/renderdb/httpext"
)

// RegisterLayerssRoutes registers handlers for the "/worlds/{worldID}/layers"-route.
func RegisterLayersRoutes(router *mux.Router, db *sqlx.DB) {
	renderer := httpext.NewJSONResponseRenderer()
	middleware := httpext.Chain(&layersMiddleware{})
	getLayers := httpext.NewHttpHandler(db, renderer, middleware.Then(&getLayersHandler{}))
	getLayer := httpext.NewHttpHandler(db, renderer, middleware.Then(&getLayerHandler{}))
	postLayer := httpext.NewHttpHandler(db, renderer, middleware.Then(&postLayerHandler{}))

	router = router.Path("/worlds/{worldID:[0-9]+}").Subrouter()
	router.Handle("/layers", getLayers).Methods("GET")
	router.Handle("/layers/{layerID:[0-9]+}", getLayer).Methods("GET")
	router.Handle("/layers", postLayer).Methods("POST")
}
