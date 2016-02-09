package geometry

import (
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/ungerik/go3d/float64/vec3"
)

func makeTimeoutChan(ms time.Duration) <-chan bool {
	timeoutCh := make(chan bool, 1)
	go func() {
		defer close(timeoutCh)
		time.Sleep(time.Second)
	}()
	return timeoutCh
}

type databaseFixture struct {
	db *sqlx.DB
}

func (f *databaseFixture) Setup(t *testing.T) {

	var err error
	f.db, err = sqlx.Open("sqlite3", ":memory:")
	assert.NoError(t, err, "Could not open database")

	f.db.MustExec(`
        CREATE TABLE IF NOT EXISTS geometry_objects(
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                bounds_x_min REAL NOT NULL,
                bounds_y_min REAL NOT NULL,
                bounds_z_min REAL NOT NULL,
                bounds_x_max REAL NOT NULL,
                bounds_y_max REAL NOT NULL,
                bounds_z_max REAL NOT NULL,
                geometry_data BLOB NOT NULL,
                metadata STRING NOT NULL)`)
}

func (f *databaseFixture) Teardown(t *testing.T) {
	assert.NoError(t, f.db.Close())
}

func TestDatabase_Add_NilElement_ReturnsError(t *testing.T) {
	// Arrange
	f := databaseFixture{}
	f.Setup(t)
	defer f.Teardown(t)
	database := sqlDatabase{f.db}

	// Act
	_, err := database.add(nil)

	// Assert
	assert.Error(t, err)
}

func TestDatabase_Add_ValidElement_InsertsElement(t *testing.T) {
	// Arrange
	f := databaseFixture{}
	f.Setup(t)
	defer f.Teardown(t)
	database := sqlDatabase{f.db}
	obj := new(SimpleObject)
	obj.bounds = &vec3.Box{}
	obj.geometryData = []byte{1}

	// Act
	id, err := database.add(obj)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
}

func TestDatabase_GetMany_NonExistantId_ReturnsError(t *testing.T) {
	// Arrange
	f := databaseFixture{}
	f.Setup(t)
	defer f.Teardown(t)
	database := sqlDatabase{f.db}

	// Act
	_, errCh := database.getMany([]int64{1337})

	// Assert
	err, _ := <-errCh
	assert.Error(t, err)
}

func TestDatabase_GetMany_NoIdsRequested_ReturnsEmpty(t *testing.T) {
	// Arrange
	f := databaseFixture{}
	f.Setup(t)
	defer f.Teardown(t)
	database := sqlDatabase{f.db}

	// Act
	dataCh, errCh := database.getMany([]int64{})

	// Assert
	_, dataChOpen := <-dataCh
	assert.False(t, dataChOpen)
	select {
	case err := <-errCh:
		assert.NoError(t, err)
	default:
	}
}

func TestDatabase_GetMany_ValidId_ReturnsData(t *testing.T) {
	// Arrange
	f := databaseFixture{}
	f.Setup(t)
	defer f.Teardown(t)
	database := sqlDatabase{f.db}
	r, _ := f.db.Exec(insertGeometrySQL, 0, 0, 0, 1, 1, 1, "ABC", "{}")
	id, _ := r.LastInsertId()
	rows, _ := f.db.Query("SELECT id FROM geometry_objects WHERE ID IN (1)")
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&id)
	}
	// Act
	dataCh, errCh := database.getMany([]int64{id})

	// Assert
	select {
	case data := <-dataCh:
		assert.NotNil(t, data)
	case err := <-errCh:
		assert.Fail(t, "Did not expect to receive error", "%v", err)
	case <-makeTimeoutChan(time.Second):
		assert.Fail(t, "Timeout while waiting for data")
	}
}

func TestDatabase_GetAll_PopulatedDb_ReturnsData(t *testing.T) {
	// Arrange
	f := databaseFixture{}
	f.Setup(t)
	defer f.Teardown(t)
	database := sqlDatabase{f.db}
	_, err := f.db.Exec(insertGeometrySQL, 0, 0, 0, 1, 1, 1, "", "{}")
	assert.NoError(t, err)

	// Act
	dataCh, errCh := database.getAll()

	// Assert
	select {
	case data := <-dataCh:
		assert.NotNil(t, data)
	case err := <-errCh:
		assert.Fail(t, "Did not expect to receive error", "%v", err)
	case <-makeTimeoutChan(time.Second):
		assert.Fail(t, "Timeout while waiting for data")
	}
}
