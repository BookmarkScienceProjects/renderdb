package sql

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/larsmoa/renderdb/db/sql"
	"github.com/stretchr/testify/assert"
)

func TestInitialize_FromEmptyDatabase_Succeeds(t *testing.T) {
	// Arrange
	db, err := sqlx.Open("sqlite3", ":memory:")
	assert.NoError(t, err, "Could not open database")
	err = sql.Initialize(f.db)
	assert.NoError(t, err, "Could not initialize database")

	// Act
	err = Initialize(db)

	// Assert
	assert.NoError(t, err)
}
