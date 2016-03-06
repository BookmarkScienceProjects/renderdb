package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/larsmoa/renderdb/db/helpers"
)

// World represents collection of data that relates to each other, but not necessarily
// to data in other worlds. A world consists of many 'layers' - e.g. buildings, roads,
// landscape etc.
type World struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

// Worlds is an interface for retrieving worlds which contains
// 'layers' which in turn contains 'scenes' of geometric objects.
type Worlds interface {
	// GetAll returns all known 'worlds'.
	GetAll() ([]*World, error)
	// Get returns the world with the given ID.
	Get(id int64) (*World, error)
	// Add adds a new world and returns the ID
	// of the world. The ID field provided world is also
	// updated.
	Add(world *World) (int64, error)
	// Delete deletes the world with the given ID.
	Delete(worldid int64) error
}

const (
	getAllWorldsSQL string = "SELECT id, name FROM worlds ORDER BY name"
	getWorldSQL     string = "SELECT id, name FROM worlds WHERE id = ?"
	addWorldSQL     string = "INSERT INTO worlds(name) VALUES (:name)"
	deleteWorldSQL  string = "DELETE FROM worlds WHERE id = ?"
)

type worldsDb struct {
	tx *sqlx.Tx
}

func worldConstructor() interface{} {
	return new(World)
}

func (db *worldsDb) GetAll() ([]*World, error) {
	items, err := helpers.GetAll(db.tx, worldConstructor, getAllWorldsSQL)
	worlds := make([]*World, len(items))
	for i, s := range items {
		worlds[i] = s.(*World)
	}
	return worlds, err
}

func (db *worldsDb) Get(id int64) (*World, error) {
	item, err := helpers.Get(db.tx, worldConstructor, getWorldSQL, id)
	return item.(*World), err
}

func (db *worldsDb) Add(world *World) (int64, error) {
	result, err := db.tx.NamedExec(addWorldSQL, world)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func (db *worldsDb) Delete(id int64) error {
	_, err := db.tx.Exec(deleteWorldSQL, id)
	return err
}
