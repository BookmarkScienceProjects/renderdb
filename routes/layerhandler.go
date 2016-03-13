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
// Middleware for injecting db.Layers to the context.
// --------------------------------------------------
type layersDBKeyType int

const layersDBKey layersDBKeyType = 0

type layersMiddleware struct{}

func (h *layersMiddleware) Handle(tx *sqlx.Tx, _ httpext.ResponseRenderer,
	_ http.ResponseWriter, r *http.Request) error {
	layersDB := db.NewLayersDB(tx)
	context.Set(r, layersDBKey, layersDB)
	return nil
}

func getLayersFromContext(r *http.Request) db.Layers {
	layers, ok := context.GetOk(r, layersDBKey)
	if !ok {
		panic("Layers not available in context, forgot layersMiddleware?")
	}
	return layers.(db.Layers)
}

// -------------------------------
// GET /worlds/{worldID}/layers
// -------------------------------

type getLayersHandler struct{}

func (h *getLayersHandler) Handle(tx *sqlx.Tx, renderer httpext.ResponseRenderer,
	w http.ResponseWriter, r *http.Request) error {

	var err error
	// Parse URL
	vars := mux.Vars(r)
	worldID, err := httpext.ReadInt64ID(vars, "worldID")
	if err != nil {
		fmt.Println("getLayersHandler err1", err)
		renderer.WriteError(w, err)
		return err
	}

	// Read from database
	layersDB := getLayersFromContext(r)
	layers, err := layersDB.GetAll(worldID)
	if err != nil {
		renderer.WriteError(w, err)
		return err
	}

	renderer.WriteObject(w, http.StatusOK, layers)
	return nil
}

// -------------------------------
// GET /worlds/{worldID}/layers/{layerID}
// -------------------------------

type getLayerHandler struct{}

func (h *getLayerHandler) Handle(tx *sqlx.Tx, renderer httpext.ResponseRenderer,
	w http.ResponseWriter, r *http.Request) error {
	fmt.Println("getLayersHandler begin")

	var err error
	// Parse URL
	vars := mux.Vars(r)
	worldID, err := httpext.ReadInt64ID(vars, "worldID")
	if err != nil {
		fmt.Println("getLayersHandler err1")
		renderer.WriteError(w, err)
		return err
	}
	layerID, err := httpext.ReadInt64ID(vars, "layerID")
	if err != nil {
		fmt.Println("getLayersHandler err2")
		renderer.WriteError(w, err)
		return err
	}

	// Read from database
	fmt.Println("getLayersHandler", worldID, layerID)
	layersDB := getLayersFromContext(r)
	layer, err := layersDB.Get(worldID, layerID)
	if err != nil {
		err = httpext.NewHttpError(fmt.Errorf("Could not retrieve layer with id %d in world %d (reason: %s)", layerID, worldID, err), http.StatusInternalServerError)
		renderer.WriteError(w, err)
		return err
	}

	// Respond
	if layer == nil {
		err = httpext.NewHttpError(fmt.Errorf("No layer with id %d in world %d", layerID, worldID), http.StatusNotFound)
		renderer.WriteError(w, err)
		return err
	}
	renderer.WriteObject(w, http.StatusOK, layer)
	return nil
}

// -------------------------------
// POST /worlds/{worldID}/layers
// -------------------------------

type postLayerHandler struct {
}

func (h *postLayerHandler) Handle(tx *sqlx.Tx, renderer httpext.ResponseRenderer,
	w http.ResponseWriter, r *http.Request) error {

	var err error
	// Parse URL
	vars := mux.Vars(r)
	worldID, err := httpext.ReadInt64ID(vars, "worldID")
	if err != nil {
		renderer.WriteError(w, err)
		return err
	}

	// Parse body
	layer, err := parseLayerFromBody(r)
	if err != nil {
		renderer.WriteError(w, err)
		return err
	}
	layer.WorldID = worldID

	// Add to database
	layersDB := getLayersFromContext(r)
	id, err := layersDB.Add(layer)
	if err != nil {
		renderer.WriteError(w, err)
		return err
	}
	layer.ID = id

	// Return to client
	renderer.WriteObject(w, http.StatusOK, layer)
	return nil
}

func parseLayerFromBody(r *http.Request) (*db.Layer, error) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	if !decoder.More() {
		return nil, fmt.Errorf("Request body is empty")
	}

	var layer db.Layer
	err := decoder.Decode(&layer)
	if err != nil {
		return nil, fmt.Errorf("Could not decode body (%v)", err)
	}

	// Validate
	if layer.Name == "" {
		return nil, fmt.Errorf("Field 'name' must be set")
	}
	return &layer, nil
}
