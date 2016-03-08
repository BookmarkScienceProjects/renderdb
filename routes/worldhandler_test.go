package routes

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/context"
	"github.com/jmoiron/sqlx"
	"github.com/larsmoa/renderdb/db"
	"github.com/larsmoa/renderdb/httpext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type worldHandlerFixture struct {
	mockDB sqlmock.Sqlmock
	db     *sqlx.DB
	tx     *sqlx.Tx

	worlds *db.MockWorlds

	request  *http.Request
	writer   *httptest.ResponseRecorder
	renderer *httpext.MockResponseRenderer
}

func (f *worldHandlerFixture) Setup(t *testing.T) {
	var database *sql.DB
	var err error
	database, f.mockDB, err = sqlmock.New()
	assert.NoError(t, err)

	f.mockDB.ExpectBegin()
	f.db = sqlx.NewDb(database, "")
	f.tx, err = f.db.Beginx()
	assert.NoError(t, err)

	f.request, _ = http.NewRequest("GET", "", nil)
	f.writer = httptest.NewRecorder()
	f.renderer = &httpext.MockResponseRenderer{}

	f.worlds = &db.MockWorlds{}
	context.Set(f.request, worldsDBKey, f.worlds)
}

func (f *worldHandlerFixture) Teardown(t *testing.T) {
	assert.NoError(t, f.db.Close())
}

func TestWorldsMiddleware_InjectsWorldsToContext(t *testing.T) {
	// Arrange
	f := worldHandlerFixture{}
	f.Setup(t)
	defer f.Teardown(t)
	middleware := worldsMiddleware{}

	// Act
	err := middleware.Handle(f.tx, f.renderer, f.writer, f.request)

	// Assert
	assert.NoError(t, err)
	_, ok := context.GetOk(f.request, worldsDBKey)
	assert.True(t, ok)
}

func TestGetWorldsHandler_Handle_GetAllReturnsError_WritesError(t *testing.T) {
	// Arrange
	f := worldHandlerFixture{}
	f.Setup(t)
	defer f.Teardown(t)

	f.worlds.On("GetAll").Return(nil, errors.New(""))
	f.renderer.On("WriteError", f.writer, mock.Anything)
	handler := getWorldsHandler{}

	// Act
	err := handler.Handle(f.tx, f.renderer, f.writer, f.request)

	// Assert
	assert.Error(t, err)
	f.worlds.AssertExpectations(t)
	f.renderer.AssertExpectations(t)
}

func TestGetWorldsHandler_Handle_GetAllReturnsWorlds_WritesResponse(t *testing.T) {
	// Arrange
	f := worldHandlerFixture{}
	f.Setup(t)
	defer f.Teardown(t)

	worlds := []*db.World{&db.World{}}
	f.worlds.On("GetAll").Return(worlds, nil)
	f.renderer.On("WriteObject", f.writer, 200, worlds)
	handler := getWorldsHandler{}

	// Act
	err := handler.Handle(f.tx, f.renderer, f.writer, f.request)

	// Assert
	assert.NoError(t, err)
	f.worlds.AssertExpectations(t)
	f.renderer.AssertExpectations(t)
}
