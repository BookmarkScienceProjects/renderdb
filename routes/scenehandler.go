package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/larsmoa/renderdb/db"
	"github.com/larsmoa/renderdb/httpext"
)

// --------------------------------------------------
// Middleware for injecting db.Scenes to the context.
// --------------------------------------------------
type scenesDBKeyType int

const scenesDBKey scenesDBKeyType = 0

type scenesMiddleware struct{}

func (h *scenesMiddleware) Handle(tx *sqlx.Tx, renderer httpext.ResponseRenderer,
	w http.ResponseWriter, r *http.Request) error {

	// Parse URL
	vars := mux.Vars(r)
	layerID, err := httpext.ReadInt64ID(vars, "layerID")
	if err != nil {
		renderer.WriteError(w, err)
		return err
	}

	scenesDB := db.NewScenesDB(tx, layerID)
	context.Set(r, scenesDBKey, scenesDB)
	return nil
}

func getScenesFromContext(r *http.Request) db.Scenes {
	scenes, ok := context.GetOk(r, scenesDBKey)
	if !ok {
		panic("Scenes not available in context, forgot scenesMiddleware?")
	}
	return scenes.(db.Scenes)
}

// ----------------------------------------------
// GET /worlds/{worldID}/layers/{layerID}/scenes
// ----------------------------------------------

type getScenesHandler struct{}

func (h *getScenesHandler) Handle(tx *sqlx.Tx, renderer httpext.ResponseRenderer,
	w http.ResponseWriter, r *http.Request) error {

	var err error

	// Read from database
	scenesDB := getScenesFromContext(r)
	layers, err := scenesDB.GetAll()
	if err != nil {
		renderer.WriteError(w, err)
		return err
	}

	renderer.WriteObject(w, http.StatusOK, layers)
	return nil
}

// ----------------------------------------------------------
// GET /worlds/{worldID}/layers/{layerID}/scenes/{layerID}
// ----------------------------------------------------------

type getSceneHandler struct{}

func (h *getSceneHandler) Handle(tx *sqlx.Tx, renderer httpext.ResponseRenderer,
	w http.ResponseWriter, r *http.Request) error {

	var err error
	// Parse URL
	vars := mux.Vars(r)
	sceneID, err := httpext.ReadInt64ID(vars, "sceneID")
	if err != nil {
		renderer.WriteError(w, err)
		return err
	}

	// Read from database
	scenesDB := getScenesFromContext(r)
	scene, err := scenesDB.Get(sceneID)
	if err != nil {
		err = httpext.NewHttpError(fmt.Errorf("Could not retrieve scene with id %d (reason: %s)", sceneID, err), http.StatusInternalServerError)
		renderer.WriteError(w, err)
		return err
	}

	// Respond
	if scene == nil {
		err = httpext.NewHttpError(fmt.Errorf("No scene with id %d", sceneID), http.StatusNotFound)
		renderer.WriteError(w, err)
		return err
	}
	renderer.WriteObject(w, http.StatusOK, scene)
	return nil
}

// -------------------------------------------------
// POST /worlds/{worldID}/layers/{layerID}/scenes
// -------------------------------------------------

type postSceneHandler struct {
}

func (h *postSceneHandler) Handle(tx *sqlx.Tx, renderer httpext.ResponseRenderer,
	w http.ResponseWriter, r *http.Request) error {

	var err error
	// Parse URL
	vars := mux.Vars(r)
	layerID, err := httpext.ReadInt64ID(vars, "layerID")
	if err != nil {
		renderer.WriteError(w, err)
		return err
	}

	// Parse body
	scene, err := parseSceneFromBody(r)
	if err != nil {
		renderer.WriteError(w, err)
		return err
	}
	scene.LayerID = layerID

	// Add to database
	scenesDB := getScenesFromContext(r)
	id, err := scenesDB.Add(scene)
	if err != nil {
		renderer.WriteError(w, err)
		return err
	}
	scene.ID = id

	// Return to client
	renderer.WriteObject(w, http.StatusOK, scene)
	return nil
}

func parseSceneFromBody(r *http.Request) (*db.Scene, error) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	if !decoder.More() {
		return nil, fmt.Errorf("Request body is empty")
	}

	var scene db.Scene
	err := decoder.Decode(&scene)
	if err != nil {
		return nil, fmt.Errorf("Could not decode body (%v)", err)
	}

	// Validate
	if scene.Name == "" {
		return nil, fmt.Errorf("Field 'name' must be set")
	}
	return &scene, nil
}
