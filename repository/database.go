package repository

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/ungerik/go3d/float64/vec3"
)

const (
	insertGeometrySQL string = `INSERT INTO geometry_objects(
            bounds_x_min, bounds_y_min, bounds_z_min, 
            bounds_x_max, bounds_y_max, bounds_z_max, 
            geometry_data, metadata) 
          VALUES (?, ?, ?, 
                  ?, ?, ?, 
                  ?, ?)`
	selectGeometrySQL string = `SELECT id, 
                bounds_x_min, bounds_y_min, bounds_z_min, 
                bounds_x_max, bounds_y_max, bounds_z_max,
                geometry_data, metadata 
            FROM geometry_objects`
)

type data struct {
	id           int64
	bounds       vec3.Box
	geometryData []byte
	metadata     map[string]interface{}
}

type database interface {
	add(o Object) (int64, error)
	getMany(ids []int64) (<-chan *data, <-chan error)
	getAll() (<-chan *data, <-chan error)
}

type sqlDatabase struct {
	db *sqlx.DB
}

func newSQLDatabase(db *sqlx.DB) *sqlDatabase {
	database := new(sqlDatabase)
	database.db = db
	return database
}

func (database *sqlDatabase) add(o Object) (int64, error) {
	if o == nil {
		return -1, fmt.Errorf("Cannot add nil")
	}
	jsonBuf, err := json.Marshal(o.Metadata())
	if err != nil {
		return -1, err
	}

	jsonTxt := strings.Trim(string(jsonBuf), "\"")
	bounds_min := o.Bounds().Min
	bounds_max := o.Bounds().Max
	result, err := database.db.Exec(insertGeometrySQL,
		bounds_min[0], bounds_min[1], bounds_min[2],
		bounds_max[0], bounds_max[1], bounds_max[2],
		o.GeometryData(), jsonTxt)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

type row interface {
	Scan(dest ...interface{}) error
}

func parseDataRow(r row) (*data, error) {
	data := new(data)
	var jsonTxt string
	err := r.Scan(&data.id,
		&data.bounds.Min[0], &data.bounds.Min[1], &data.bounds.Min[2],
		&data.bounds.Max[0], &data.bounds.Max[1], &data.bounds.Max[2],
		&data.geometryData, &jsonTxt)
	if err != nil {
		return nil, err
	}

	if jsonTxt != "" {
		err = json.Unmarshal([]byte(jsonTxt), &data.metadata)
	}
	return data, err
}

func (database *sqlDatabase) getMany(ids []int64) (<-chan *data, <-chan error) {
	bufferSize := 200
	dataChan := make(chan *data, bufferSize)
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
			chunkIds := ids[i:lastElement]

			// TODO: Consider if this should be optimized by creating a temporary table
			// http://explainextended.com/2009/08/18/passing-parameters-in-mysql-in-list-vs-temporary-table/
			q, args, _ := sqlx.In(fmt.Sprintf("%s WHERE id IN (?)", selectGeometrySQL), chunkIds)
			q = sqlx.Rebind(sqlx.QUESTION, q)
			rows, err := database.db.Queryx(q, args...)
			if err != nil {
				errChan <- err
				return
			}
			defer rows.Close()

			// Parse rows
			var result *data
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

func (database *sqlDatabase) getAll() (<-chan *data, <-chan error) {
	bufferSize := 200
	dataChan := make(chan *data, bufferSize)
	errChan := make(chan error)
	go func() {
		defer close(dataChan)

		rows, err := database.db.Queryx(selectGeometrySQL)
		if err != nil {
			errChan <- err
			return
		}
		defer rows.Close()

		// Parse rows
		var result *data
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
