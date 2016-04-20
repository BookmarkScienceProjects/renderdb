// Package routes implements a REST API for dealing with geometry.
//
// The API is layered into 'worlds', 'layers', 'scenes' and 'objects'.
//
// - A 'world' typically defines a separate project, e.g. a building site.
//   Each world can have many layers.
// - A 'layer' typically defines a group of data, e.g. all the plumbing in the building.
//   Each layer can have many scenes.
// - A 'scene' is the smallest distingushable feature and can e.g. represents
//   all plumbing in floor 2. Each scene can have many 'objects'.
// - 'Objects' are geometric 3D entities that can be rendered. Each object
//   can have JSON metadata which can be used for dynamic filtering.
//   There is no API to add single objects or to query for objects based
//   on IDs. To add objects a new 'scene' must be added.
//
// Data management endpoints:
// --------------------------
// POST 	/worlds
// - Adds a new world
// GET  	/worlds
// - Returns metadata for all known worlds
// GET  	/worlds/{id}
// - Returns metadata the world with the given ID
// DELETE 	/worlds/{id} 	(Not implemented yet)
// - Deletes the world with the given ID. Deletes all
//   layers in the world.
// POST 	/worlds/{id}/layers
// - Adds a new layer to the given world
// GET 		/worlds/{id}/layers
// - Returns metadata for all layers in the world
// GET 		/worlds/{id}/layers/{id}
// - Returns metadata for the layer with the given ID
// DELETE   /worlds/{id}/layers/{id} 	(Not implemented yet)
// - Deletes the layer with the given ID. Deletes all
//   scenes in the layer.
// POST 	/worlds/{id}/layers/{id}/scenes
// - Adds a new scene to the given layer. Scenes are specified
//   using the Wavefront OBJ-format and each group
//   in the file is considered to be a separate object.
// PUT 		/worlds/{id}/layers/{id}/scenes/{id}	(Not implemented yet)
// - Replaces all geometry in a scene. Supports the same
//   formats as the POST request.
// GET 		/worlds/{id}/layers/{id}/scenes
// - Returns metadata for all scenes in the layer.
// GET 		/worlds/{id}/layers/{id}/scenes/{id}
// - Returns metadata for the given scene.
// DELETE 	/worlds/{id}/layers/{id}/scenes/{id}	(Not implemented yet)
// - Deletes the scene with the given ID and all the objects
//   in the scene.
//
// Geometry query endpoints:
// -------------------------
// GET /world/{id}/geometry?{filter}&{options}	(Not implemented yet)
// - Gets all geometry in the world that matches the filter.
// GET /world/{id}/layers/{id}/geometry?{filter}&{options}	(Not implemented yet)
// - Gets all geometry in the layer that matches the filter.
//
// Filters is used to filter away unwanted data, e.g. based on location or distance
// to camera.
// Options are used to e.g. sort the results by distance to a camera, or
// restrict the number of returned triangles
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

// RegisterLayersRoutes registers handlers for the "/worlds/{worldID}/layers"-route.
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

/*
// RegisterGeometryRoutes registers handelrs for the "/worlds/{worldID}/layers/{layerID}/scenes/{sceneID}/geometry"-route.
func RegisterGeometryRoutes(router *mux.Router, db *sqlx.DB) {
	renderer := httpext.NewJSONResponseRenderer()

	middleware := httpext.Chain(&objectsMiddleware{})
	postSceneGeometry := httpext.NewHttpHandler(db, renderer, middleware.Then(&postSceneGeometryHandler{}))
	router = router.Path("/worlds/{worldID:[0-9]+}/layers/{layerID:[0-9]+}/scenes/{sceneID:[0-9]+}").Subrouter()
	router.Handle("/geometry", postSceneGeometry).Methods("POST")
}
*/
