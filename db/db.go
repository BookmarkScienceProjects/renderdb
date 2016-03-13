package db

import "github.com/jmoiron/sqlx"

func NewWorldsDB(tx *sqlx.Tx) Worlds {
	return &worldsDb{tx}
}

func NewLayersDB(tx *sqlx.Tx, worldID int64) Layers {
	return &layersDb{tx, worldID}
}

func NewScenesDB(tx *sqlx.Tx, scenes *Scene) Scenes {
	return &scenesDb{scenes.ID, tx}
}

func NewObjectsDb(tx *sqlx.Tx, world *World) Objects {
	return &objectsDb{world.ID, tx}
}
