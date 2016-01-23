package geometry

import (
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type data struct {
	id           int64
	geometryText string
	metadata     map[string]string
}

type database interface {
	add(o Object) (int64, error)
	getMany(ids []int64) (<-chan *data, <-chan error)
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
	jsonTxt, err := json.Marshal(o.Metadata())
	if err != nil {
		return -1, err
	}

	result, err := database.db.Exec("INSERT INTO geometry_objects(geometry_text, metadata) VALUES ($1, $2)",
		o.GeometryText(), jsonTxt)
	return result.LastInsertId()
}

type row interface {
	Scan(dest ...interface{}) error
}

func parseDataRow(r row) (*data, error) {
	data := new(data)
	var jsonTxt string
	err := r.Scan(&data.id, &data.geometryText, &jsonTxt)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(jsonTxt), &data.metadata)
	return data, err
}

func (database *sqlDatabase) getMany(ids []int64) (<-chan *data, <-chan error) {
	bufferSize := 200
	dataChan := make(chan *data, bufferSize)
	errChan := make(chan error)
	go func() {
		defer close(dataChan)
		defer close(errChan)

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
			rows, err := database.db.Queryx("SELECT id, geometry_text, metadata FROM geometry_objects WHERE id IN $1", chunkIds)
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
