package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/larsmoa/renderdb/db/helpers"
)

// Layer represents a collection of scenes that combined form 'a layer'.
//
// All scenes in a layer is visualized together, but grouping scenes into
// layers allow updating one scene without having to update all scenes in
// the layer.
//
// A layer lives in a 'world'
type Layer struct {
	ID      int64  `db:"id"`
	WorldID int64  `db:"world_id"`
	Name    string `db:"name"`
}

// Layers provide functionality for accessing layers from a database.
type Layers interface {
	// GetAll returns all layers of the given world.
	GetAll() ([]*Layer, error)
	// Get returns the layer with the specified ID.
	Get(layerid int64) (*Layer, error)
	// Add creates a new layer and returns the ID, or an error.
	Add(layer *Layer) (int64, error)
	// Delete deletes the layer with the given ID.
	Delete(layerid int64) error
}

const (
	getAllLayersSQL string = "SELECT id, world_id, name FROM layers WHERE world_id = ? ORDER BY name"
	getLayerSQL     string = "SELECT id, world_id, name FROM layers WHERE id = ? AND world_id = ?"
	addLayerSQL     string = "INSERT INTO layers(world_id, name) VALUES (:world_id, :name)"
	deleteLayerSQL  string = "DELETE FROM layers WHERE id = ? AND world_id = ?"
)

type layersDb struct {
	tx      *sqlx.Tx
	worldID int64
}

func layerConstructor() interface{} {
	return new(Layer)
}

func (db *layersDb) GetAll() ([]*Layer, error) {
	items, err := helpers.GetAll(db.tx, sceneConstructor, getAllLayersSQL, db.worldID)
	layers := make([]*Layer, len(items))
	for i, s := range items {
		layers[i] = s.(*Layer)
	}
	return layers, err
}

func (db *layersDb) Get(layerid int64) (*Layer, error) {
	item, err := helpers.Get(db.tx, sceneConstructor, getLayerSQL, layerid, db.worldID)
	return item.(*Layer), err
}

func (db *layersDb) Add(layer *Layer) (int64, error) {
	layer.WorldID = db.worldID
	result, _ := db.tx.NamedExec(addLayerSQL, layer)
	return result.LastInsertId()
}

func (db *layersDb) Delete(id int64) error {
	_, err := db.tx.Exec(deleteLayerSQL, id)
	return err
}
