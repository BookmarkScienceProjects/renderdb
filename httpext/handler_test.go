package httpext

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

type newHTTPHandlerFixture struct {
	mockDB sqlmock.Sqlmock
	db     *sqlx.DB

	renderer ResponseRenderer
	writer   *httptest.ResponseRecorder
	request  *http.Request
}

func (f *newHTTPHandlerFixture) Setup(t *testing.T) {
	var err error
	var db *sql.DB
	db, f.mockDB, err = sqlmock.New()
	assert.NoError(t, err)

	f.db = sqlx.NewDb(db, "")
	assert.NoError(t, err)

	f.renderer = NewJSONResponseRenderer()
	f.writer = httptest.NewRecorder()
}

func (f *newHTTPHandlerFixture) Teardown(t *testing.T) {
	f.mockDB.ExpectClose()
	assert.NoError(t, f.db.Close())
}

func TestNewHttpHandler_InnerHandlerFails_RollsbackTransactionAndWritesError(t *testing.T) {
	// Arrange
	f := newHTTPHandlerFixture{}
	f.Setup(t)
	defer f.Teardown(t)
	h := mockHandler{}
	h.On("Handle", any, any, any, any).
		Return(NewHttpError(errors.New(""), http.StatusConflict))

	f.mockDB.ExpectBegin()
	f.mockDB.ExpectRollback()

	// Act
	handler := NewHttpHandler(f.db, f.renderer, &h)
	handler.ServeHTTP(f.writer, f.request)

	// Assert
	assert.NoError(t, f.mockDB.ExpectationsWereMet())
	assert.Equal(t, http.StatusConflict, f.writer.Code)
	assert.NotZero(t, f.writer.Body.Len())
}

func TestNewHttpHandler_InnerHandlerSucceeds_CommitsTransaction(t *testing.T) {
	// Arrange
	f := newHTTPHandlerFixture{}
	f.Setup(t)
	defer f.Teardown(t)
	h := mockHandler{}
	h.On("Handle", any, any, any, any).Return(nil)

	f.mockDB.ExpectBegin()
	f.mockDB.ExpectCommit()

	// Act
	handler := NewHttpHandler(f.db, f.renderer, &h)
	handler.ServeHTTP(f.writer, f.request)

	// Assert
	assert.NoError(t, f.mockDB.ExpectationsWereMet())
}

func TestNewHttpHandler_OpenTransactionFails_WritesErrorAndAborts(t *testing.T) {
	// Arrange
	f := newHTTPHandlerFixture{}
	f.Setup(t)
	defer f.Teardown(t)
	h := mockHandler{}
	h.On("Handle", any, any, any, any).Return(nil)

	f.mockDB.ExpectBegin().WillReturnError(errors.New(""))

	// Act
	handler := NewHttpHandler(f.db, f.renderer, &h)
	handler.ServeHTTP(f.writer, f.request)

	// Assert
	assert.NoError(t, f.mockDB.ExpectationsWereMet())
}
