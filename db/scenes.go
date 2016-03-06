package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/larsmoa/renderdb/db/helpers"
)

// Scene represents a set of geometric objects in a 'layer'.
// Scenes is the smallest updatable entity - to add or update
// geometric objects an entire scene must be added/updated.
type Scene struct {
	ID      int64  `db:"id"`
	LayerID int64  `db:"layer_id"`
	Name    string `db:"name"`
}

// Scenes contains functionality for adding and deleting scenes of
// a layer.
type Scenes interface {
	// GetAll returns all scenes in a layer.
	GetAll(layerid int64) ([]*Scene, error)
	// Get returns the scene with the given ID.
	Get(id int64) (*Scene, error)
	// Add creates a new scene in the database and returns the ID.
	Add(scene *Scene) (int64, error)
	// Delete deletes the scene with the given ID from the database.
	Delete(sceneid int64) error
}

const (
	getAllScenesSQL string = "SELECT id, layer_id, name FROM scenes WHERE layer_id = ?"
	getSceneSQL     string = "SELECT id, layer_id, name FROM scenes WHERE id = ?"
	addSceneSQL     string = "INSERT INTO scenes(layer_id, name) VALUES(:layer_id, :name)"
	deleteScenesSQL string = "DELETE FROM scenes WHERE id = ?"
)

type scenesDb struct {
	sceneID int64
	tx      *sqlx.Tx
}

func sceneConstructor() interface{} {
	return new(Scene)
}

func (db *scenesDb) GetAll(layerid int64) ([]*Scene, error) {
	items, err := helpers.GetAll(db.tx, sceneConstructor, getAllScenesSQL, db.sceneID)
	scenes := make([]*Scene, len(items))
	for i, s := range items {
		scenes[i] = s.(*Scene)
	}
	return scenes, err
}

func (db *scenesDb) Get(id int64) (*Scene, error) {
	item, err := helpers.Get(db.tx, sceneConstructor, getSceneSQL, id)
	return item.(*Scene), err
}

func (db *scenesDb) Add(scene *Scene) (int64, error) {
	result, _ := db.tx.NamedExec(addSceneSQL, scene)
	return result.LastInsertId()
}

func (db *scenesDb) Delete(id int64) error {
	_, err := db.tx.Exec(deleteLayerSQL, id)
	return err
}
