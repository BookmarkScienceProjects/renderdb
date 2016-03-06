-- +migrate Up
CREATE TABLE worlds(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name REAL NOT NULL);

-- +migrate Down
DROP TABLE worlds;