package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/larsmoa/renderdb/db"
	"github.com/larsmoa/renderdb/httpext"
)

// --------------------------------------------------
// Middleware for injecting db.Worlds to the context.
// --------------------------------------------------
type worldsDBKeyType int

const worldsDBKey worldsDBKeyType = 0

type worldsMiddleware struct{}

func (h *worldsMiddleware) Handle(tx *sqlx.Tx, _ httpext.ResponseRenderer,
	_ http.ResponseWriter, r *http.Request) error {
	worldsDB := db.NewWorldsDB(tx)
	context.Set(r, worldsDBKey, worldsDB)
	return nil
}

func getWorldsFromContext(r *http.Request) db.Worlds {
	worlds, ok := context.GetOk(r, worldsDBKey)
	if !ok {
		panic("Worlds not available in context, forgot worldsMiddleware?")
	}
	return worlds.(db.Worlds)
}

// -------------------------------
// GET /worlds
// -------------------------------

type getWorldsHandler struct{}

func (h *getWorldsHandler) Handle(tx *sqlx.Tx, renderer httpext.ResponseRenderer,
	w http.ResponseWriter, r *http.Request) error {

	worldsDB := getWorldsFromContext(r)
	worlds, err := worldsDB.GetAll()
	if err != nil {
		renderer.WriteError(w, err)
		return err
	}

	renderer.WriteObject(w, http.StatusOK, worlds)
	return nil
}

// -------------------------------
// GET /worlds/[id]
// -------------------------------

type getWorldHandler struct{}

func (h *getWorldHandler) Handle(tx *sqlx.Tx, renderer httpext.ResponseRenderer,
	w http.ResponseWriter, r *http.Request) error {

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		return httpext.NewHttpError(fmt.Errorf("id must be a number (got '%s')", vars["id"]), http.StatusBadRequest)
	}

	worldsDB := getWorldsFromContext(r)
	world, err := worldsDB.Get(id)
	if err != nil {
		return httpext.NewHttpError(fmt.Errorf("Could not retrieve world with id %d (reason: %s)", id, err), http.StatusInternalServerError)
	}

	if world == nil {
		return httpext.NewHttpError(fmt.Errorf("No world with id %d", id), http.StatusNotFound)
	}
	renderer.WriteObject(w, http.StatusOK, world)
	return nil
}

// -------------------------------
// POST /worlds
// -------------------------------

type postWorldHandler struct {
}

func (h *postWorldHandler) Handle(tx *sqlx.Tx, renderer httpext.ResponseRenderer,
	w http.ResponseWriter, r *http.Request) error {
	// Parse body
	world, err := parseWorldFromBody(r)
	if err != nil {
		renderer.WriteError(w, err)
		return err
	}

	// Add to database
	worldsDB := getWorldsFromContext(r)
	id, err := worldsDB.Add(world)
	if err != nil {
		renderer.WriteError(w, err)
		return err
	}

	// Return to client
	world.ID = id
	renderer.WriteObject(w, http.StatusOK, world)
	return nil
}

func parseWorldFromBody(r *http.Request) (*db.World, error) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	if !decoder.More() {
		return nil, fmt.Errorf("Request body is empty")
	}

	var world db.World
	err := decoder.Decode(&world)
	if err != nil {
		return nil, fmt.Errorf("Could not decode body (%v)", err)
	}

	// Validate
	if world.Name == "" {
		return nil, fmt.Errorf("Field 'name' must be set")
	}
	return &world, nil
}
