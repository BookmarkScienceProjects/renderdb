-- +migrate Up
CREATE TABLE worlds(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name REAL NOT NULL);
CREATE TABLE layers(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    world_id INTEGER,
    name REAL NOT NULL,
    FOREIGN KEY(world_id) REFERENCES worlds(id)
);
CREATE TABLE scenes(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    layer_id INTEGER,
    name REAL NOT NULL,
    FOREIGN KEY (layer_id) REFERENCES layers(id)
);
CREATE TABLE geometry_objects(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    world_id INTEGER NOT NULL,
    layer_id INTEGER NOT NULL,
    scene_id INTEGER NOT NULL,
    bounds_x_min REAL NOT NULL,
    bounds_y_min REAL NOT NULL,
    bounds_z_min REAL NOT NULL,
    bounds_x_max REAL NOT NULL,
    bounds_y_max REAL NOT NULL,
    bounds_z_max REAL NOT NULL,
    geometry_data BLOB NOT NULL,
    metadata STRING NOT NULL,
	FOREIGN KEY(world_id) REFERENCES worlds(id),
	FOREIGN KEY(layer_id) REFERENCES layers(id),
	FOREIGN KEY(scene_id) REFERENCES scenes(id));

-- +migrate Down
DROP TABLE geometry_objects;
DROP TABLE scenes;
DROP TABLE layers;
DROP TABLE worlds;