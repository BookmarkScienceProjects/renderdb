-- +migrate Up
CREATE TABLE layers(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    world_id INTEGER,
    name REAL NOT NULL,
    FOREIGN KEY(world_id) REFERENCES worlds(id)
);

-- +migrate Down
DROP TABLE layers;