package geometry

import (
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type geometryData struct {
	id           int64
	geometryText string
	metadata     map[string]string
}

type geometryDatabase interface {
	add(o GeometryObject) (int64, error)
	getMany(ids []int64) (<-chan *geometryData, <-chan error)
}

type sqlGeometryDatabase struct {
	db *sqlx.DB
}

func newSQLGeometryDatabase(db *sqlx.DB) *sqlGeometryDatabase {
	database := new(sqlGeometryDatabase)
	database.db = db
	return database
}

func (database *sqlGeometryDatabase) add(o GeometryObject) (int64, error) {
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

func parseGeometryDataRow(r row) (*geometryData, error) {
	data := new(geometryData)
	var jsonTxt string
	err := r.Scan(&data.id, &data.geometryText, &jsonTxt)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(jsonTxt), &data.metadata)
	return data, err
}

func (database *sqlGeometryDatabase) getMany(ids []int64) (<-chan *geometryData, <-chan error) {
	bufferSize := 200
	dataChan := make(chan *geometryData, bufferSize)
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
			var result *geometryData
			for rows.Next() {
				result, err = parseGeometryDataRow(rows)
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
