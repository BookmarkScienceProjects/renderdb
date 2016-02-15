-- +migrate Up
CREATE TABLE IF NOT EXISTS geometry_objects(
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                bounds_x_min REAL NOT NULL,
                bounds_y_min REAL NOT NULL,
                bounds_z_min REAL NOT NULL,
                bounds_x_max REAL NOT NULL,
                bounds_y_max REAL NOT NULL,
                bounds_z_max REAL NOT NULL,
                geometry_data BLOB NOT NULL,
                metadata STRING NOT NULL);
                
-- +migrate Down
DROP TABLE geometry_objects;