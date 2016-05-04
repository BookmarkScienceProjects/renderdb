package sql

import (
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestInitialize_FromEmptyDatabase_Succeeds(t *testing.T) {
	// Arrange
	db, err := sqlx.Open("sqlite3", ":memory:")
	assert.NoError(t, err, "Could not open database")

	// Act
	err = Initialize(db)

	// Assert
	assert.NoError(t, err)
}
