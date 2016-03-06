package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/larsmoa/renderdb/db"
	"github.com/larsmoa/renderdb/httpext"
)

// -------------------------------
// GET /worlds
// -------------------------------

type getWorldsHandler struct{}

func (h *getWorldsHandler) Handle(tx *sqlx.Tx, renderer httpext.ResponseRenderer,
	w http.ResponseWriter, r *http.Request) error {

	worldsDB := db.NewWorldsDB(tx)
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
		err = httpext.NewHttpError(fmt.Errorf("id must be a number (got '%s')", vars["id"]), http.StatusBadRequest)
		renderer.WriteError(w, err)
		return err
	}

	worldsDB := db.NewWorldsDB(tx)
	world, err := worldsDB.Get(id)
	if err != nil {
		err = httpext.NewHttpError(fmt.Errorf("id must be a number (got '%s')", vars["id"]), http.StatusBadRequest)
		renderer.WriteError(w, err)
		return err
	}

	if world != nil {
		renderer.WriteObject(w, http.StatusNotFound, nil)
		return sql.ErrNoRows
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
	worldsDB := db.NewWorldsDB(tx)
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
