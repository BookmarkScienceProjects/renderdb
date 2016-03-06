-- +migrate Up
CREATE TABLE scenes(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    layer_id INTEGER,
    name REAL NOT NULL,
    FOREIGN KEY (layer_id) REFERENCES layers(id)
);

-- +migrate Down
DROP TABLE scenes;