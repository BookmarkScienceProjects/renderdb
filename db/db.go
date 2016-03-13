package db

import "github.com/jmoiron/sqlx"

func NewWorldsDB(tx *sqlx.Tx) Worlds {
	return &worldsDb{tx}
}

func NewLayersDB(tx *sqlx.Tx, worldID int64) Layers {
	return &layersDb{tx, worldID}
}

func NewScenesDB(tx *sqlx.Tx, layerID int64) Scenes {
	return &scenesDb{tx, layerID}
}

func NewObjectsDb(tx *sqlx.Tx, world *World) Objects {
	return &objectsDb{world.ID, tx}
}
