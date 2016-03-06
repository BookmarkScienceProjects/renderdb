package db

import "github.com/jmoiron/sqlx"

func NewWorldsDB(tx *sqlx.Tx) Worlds {
	return &worldsDb{tx}
}

func NewLayersDB(tx *sqlx.Tx, layers *Layer) Layers {
	return &layersDb{layers.ID, tx}
}

func NewScenesDB(tx *sqlx.Tx, scenes *Scene) Scenes {
	return &scenesDb{scenes.ID, tx}
}

func NewObjectsDb(tx *sqlx.Tx, world *World) Objects {
	return &objectsDb{world.ID, tx}
}
