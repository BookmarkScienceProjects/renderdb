package routes

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/larsmoa/renderdb/httpext"
)

// RegisterScenesRoutes registers handlers for the "/worlds/{worldID}/layers/{layerID}/scenes"-route.
func RegisterScenesRoutes(router *mux.Router, db *sqlx.DB) {
	renderer := httpext.NewJSONResponseRenderer()
	middleware := httpext.Chain(&layersMiddleware{})
	getScenes := httpext.NewHttpHandler(db, renderer, middleware.Then(&getScenesHandler{}))
	getScene := httpext.NewHttpHandler(db, renderer, middleware.Then(&getSceneHandler{}))
	postScene := httpext.NewHttpHandler(db, renderer, middleware.Then(&postSceneHandler{}))

	router = router.Path("/worlds/{worldID:[0-9]+}/layers/{layerID:[0-9]+}").Subrouter()
	router.Handle("/scenes", getScenes).Methods("GET")
	router.Handle("/scenes/{sceneID:[0-9]+}", getScene).Methods("GET")
	router.Handle("/scenes", postScene).Methods("POST")
}
