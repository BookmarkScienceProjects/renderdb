package db

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/ungerik/go3d/float64/vec3"
)

const (
	insertGeometrySQL string = `INSERT INTO geometry_objects(
            world_id, layer_id, scene_id,
            bounds_x_min, bounds_y_min, bounds_z_min, 
            bounds_x_max, bounds_y_max, bounds_z_max, 
            geometry_data, metadata) 
          VALUES (?, ?, ?, 
                  ?, ?, ?, 
                  ?, ?)`
	selectGeometrySQL string = `SELECT id,
                world_id, layer_id, scene_id,
                bounds_x_min, bounds_y_min, bounds_z_min, 
                bounds_x_max, bounds_y_max, bounds_z_max,
                geometry_data, metadata 
            FROM geometry_objects WHERE world_id = ?`
)

type ObjectSelector interface {
	CreateWhereClause() (string, []interface{})
}

// Objects represents a collection of geometric entities assosciated with
// a 'world'. Note that this collection holds all objects for all
// layers and scenes in a world. To restrict the result from the
// queries to only include certain layers and/or scenes 'ObjectSelector's
// can be used.
type Objects interface {
	Add(o Object) (int64, error)
	GetMany(ids []int64) (<-chan Object, <-chan error)
	GetAll() (<-chan Object, <-chan error)
}

type objectsDb struct {
	worldID int64
	tx      *sqlx.Tx
}

func (db *objectsDb) Add(o Object) (int64, error) {
	if o == nil {
		return -1, fmt.Errorf("Cannot add nil")
	}
	jsonBuf, err := json.Marshal(o.Metadata())
	if err != nil {
		return -1, err
	}

	jsonTxt := strings.Trim(string(jsonBuf), "\"")
	boundsMin := o.Bounds().Min
	boundsMax := o.Bounds().Max
	result, err := db.tx.Exec(insertGeometrySQL,
		boundsMin[0], boundsMin[1], boundsMin[2],
		boundsMax[0], boundsMax[1], boundsMax[2],
		o.GeometryData(), jsonTxt)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func (db *objectsDb) GetMany(ids []int64) (<-chan Object, <-chan error) {
	bufferSize := 200
	dataChan := make(chan Object, bufferSize)
	errChan := make(chan error)
	go func() {
		defer close(dataChan)

		// Split into several fetch operations
		retrievedCount := 0
		for i := 0; i < len(ids); i = i + bufferSize {
			lastElement := i + bufferSize
			if lastElement > len(ids) {
				lastElement = len(ids)
			}
			chunkIds := append([]int64{db.worldID}, ids[i:lastElement]...)

			// TODO: Consider if this should be optimized by creating a temporary table
			// http://explainextended.com/2009/08/18/passing-parameters-in-mysql-in-list-vs-temporary-table/
			q, args, _ := sqlx.In(fmt.Sprintf("%s AND id IN (?)", selectGeometrySQL), chunkIds)
			q = sqlx.Rebind(sqlx.QUESTION, q)
			rows, err := db.tx.Queryx(q, args...)
			if err != nil {
				errChan <- err
				return
			}
			defer rows.Close()

			// Parse rows
			var result Object
			for rows.Next() {
				result, err = parseDataRow(rows)
				if err != nil {
					errChan <- err
					return
				}
				dataChan <- result
				retrievedCount++
			}
		}
		if retrievedCount < len(ids) {
			errChan <- fmt.Errorf("Expected %d rows, but got %d", len(ids), retrievedCount)
			return
		}
	}()
	return dataChan, errChan
}

func (db *objectsDb) GetAll() (<-chan Object, <-chan error) {
	bufferSize := 200
	dataChan := make(chan Object, bufferSize)
	errChan := make(chan error)
	go func() {
		defer close(dataChan)

		rows, err := db.tx.Queryx(selectGeometrySQL, db.worldID)
		if err != nil {
			errChan <- err
			return
		}
		defer rows.Close()

		// Parse rows
		var result Object
		for rows.Next() {
			result, err = parseDataRow(rows)
			if err != nil {
				errChan <- err
				return
			}
			dataChan <- result
		}
	}()
	return dataChan, errChan
}

// Internals below:

type row interface {
	Scan(dest ...interface{}) error
}

type objectData struct {
	id           int64
	worldID      int64
	layerID      int64
	sceneID      int64
	bounds       vec3.Box
	geometryData []byte
	metadata     map[string]interface{}
}

func parseDataRow(r row) (Object, error) {
	data := new(objectData)
	var jsonTxt string

	err := r.Scan(&data.id,
		&data.worldID, &data.layerID, &data.sceneID,
		&data.bounds.Min[0], &data.bounds.Min[1], &data.bounds.Min[2],
		&data.bounds.Max[0], &data.bounds.Max[1], &data.bounds.Max[2],
		&data.geometryData, &jsonTxt)
	if err != nil {
		return nil, err
	}

	if jsonTxt != "" {
		err = json.Unmarshal([]byte(jsonTxt), &data.metadata)
	}
	if err != nil {
		return nil, err
	}
	return NewSimpleObject(data.bounds, data.geometryData, data.metadata), nil
}
