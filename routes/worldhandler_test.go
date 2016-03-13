package routes

import (
	"bytes"
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

func (f *worldHandlerFixture) Setup(t *testing.T, r *http.Request) {
	var database *sql.DB
	var err error
	database, f.mockDB, err = sqlmock.New()
	assert.NoError(t, err)

	f.mockDB.ExpectBegin()
	f.db = sqlx.NewDb(database, "")
	f.tx, err = f.db.Beginx()
	assert.NoError(t, err)

	f.writer = httptest.NewRecorder()
	f.renderer = &httpext.MockResponseRenderer{}

	f.worlds = &db.MockWorlds{}
	context.Set(r, worldsDBKey, f.worlds)
}

func (f *worldHandlerFixture) Teardown(t *testing.T) {
	assert.NoError(t, f.db.Close())
}

func TestWorldsMiddleware_InjectsWorldsToContext(t *testing.T) {
	// Arrange
	r, _ := http.NewRequest("GET", "/worlds", nil)
	f := worldHandlerFixture{}
	f.Setup(t, r)
	defer f.Teardown(t)
	middleware := worldsMiddleware{}

	// Act
	err := middleware.Handle(f.tx, f.renderer, f.writer, f.request)

	// Assert
	assert.NoError(t, err)
	_, ok := context.GetOk(f.request, worldsDBKey)
	assert.True(t, ok)
}

func TestGetWorldsHandler_GetAllReturnsError_WritesError(t *testing.T) {
	// Arrange
	r, _ := http.NewRequest("GET", "/worlds", nil)
	f := worldHandlerFixture{}
	f.Setup(t, r)
	defer f.Teardown(t)

	f.worlds.On("GetAll").Return(nil, errors.New(""))
	f.renderer.On("WriteError", f.writer, mock.Anything)
	handler := getWorldsHandler{}

	// Act
	err := httpext.InvokeHandler(&handler, "GET", "/worlds",
		f.writer, r,
		f.tx, f.renderer)

	// Assert
	assert.Error(t, err)
	f.worlds.AssertExpectations(t)
	f.renderer.AssertExpectations(t)
}

func TestGetWorldsHandler_GetAllReturnsWorlds_WritesResponse(t *testing.T) {
	// Arrange
	r, _ := http.NewRequest("GET", "/worlds", nil)
	f := worldHandlerFixture{}
	f.Setup(t, r)
	defer f.Teardown(t)

	worlds := []*db.World{&db.World{}}
	f.worlds.On("GetAll").Return(worlds, nil)
	f.renderer.On("WriteObject", f.writer, 200, worlds)
	handler := getWorldsHandler{}

	// Act
	err := httpext.InvokeHandler(&handler, "GET", "/worlds",
		f.writer, r,
		f.tx, f.renderer)

	// Assert
	assert.NoError(t, err)
	f.worlds.AssertExpectations(t)
	f.renderer.AssertExpectations(t)
}

func TestGetWorldHandle_GetReturnsError_WritesError(t *testing.T) {
	// Arrange
	r, _ := http.NewRequest("GET", "/worlds/42", nil)
	f := worldHandlerFixture{}
	f.Setup(t, r)
	defer f.Teardown(t)

	f.worlds.On("Get", int64(42)).Return(nil, errors.New(""))
	f.renderer.On("WriteError", f.writer, mock.Anything)
	handler := getWorldHandler{}

	// Act
	err := httpext.InvokeHandler(&handler, "GET", "/worlds/{worldID}",
		f.writer, r,
		f.tx, f.renderer)

	// Assert
	assert.Error(t, err)
	f.worlds.AssertExpectations(t)
	f.renderer.AssertExpectations(t)
}

func TestGetWorldHandle_GetReturnsWorld_WritesResponse(t *testing.T) {
	// Arrange
	r, _ := http.NewRequest("GET", "/worlds/42", nil)
	f := worldHandlerFixture{}
	f.Setup(t, r)
	defer f.Teardown(t)

	world := &db.World{}
	f.worlds.On("Get", int64(42)).Return(world, nil)
	f.renderer.On("WriteObject", f.writer, 200, world)
	handler := getWorldHandler{}

	// Act
	err := httpext.InvokeHandler(&handler, "GET", "/worlds/{worldID}",
		f.writer, r,
		f.tx, f.renderer)

	// Assert
	assert.NoError(t, err)
	f.worlds.AssertExpectations(t)
	f.renderer.AssertExpectations(t)
}

func TestPostWorldHandle_InvalidBody_WritesError(t *testing.T) {
	// Arrange
	buffer := bytes.NewBuffer([]byte("{}"))
	r, _ := http.NewRequest("POST", "/worlds", buffer)
	f := worldHandlerFixture{}
	f.Setup(t, r)

	f.renderer.On("WriteError", f.writer, mock.Anything)
	handler := postWorldHandler{}

	// Act
	err := httpext.InvokeHandler(&handler, "POST", "/worlds",
		f.writer, r,
		f.tx, f.renderer)

	// Assert
	assert.Error(t, err)
	f.renderer.AssertExpectations(t)
}

func TestPostWorldHandle_ValidBody_WritesObject(t *testing.T) {
	// Arrange
	buffer := bytes.NewBuffer([]byte(`{"Name": "MyWorld"}`))
	r, _ := http.NewRequest("POST", "/worlds", buffer)
	f := worldHandlerFixture{}
	f.Setup(t, r)

	f.worlds.On("Add", mock.Anything).Return(int64(11), nil)
	f.renderer.On("WriteObject", f.writer, 200, mock.Anything)
	handler := postWorldHandler{}

	// Act
	err := httpext.InvokeHandler(&handler, "POST", "/worlds",
		f.writer, r,
		f.tx, f.renderer)

	// Assert
	assert.NoError(t, err)
	f.renderer.AssertExpectations(t)
	f.worlds.AssertExpectations(t)
}
