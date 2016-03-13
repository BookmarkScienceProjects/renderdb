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

type layerHandlerFixture struct {
	mockDB sqlmock.Sqlmock
	db     *sqlx.DB
	tx     *sqlx.Tx

	layers *db.MockLayers

	writer   *httptest.ResponseRecorder
	renderer *httpext.MockResponseRenderer
}

func (f *layerHandlerFixture) Setup(t *testing.T, r *http.Request) {
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

	f.layers = &db.MockLayers{}
	context.Set(r, layersDBKey, f.layers)
}

func (f *layerHandlerFixture) Teardown(t *testing.T) {
	assert.NoError(t, f.db.Close())
}

func TestLayersMiddleware_Success(t *testing.T) {
	// Arrange
	r, _ := http.NewRequest("GET", "/worlds/13/layers", nil)
	f := layerHandlerFixture{}
	f.Setup(t, r)
	defer f.Teardown(t)
	middleware := layersMiddleware{}

	// Act
	err := httpext.InvokeHandler(&middleware, "GET", "/worlds/{worldID}/layers",
		f.writer, r, f.tx, f.renderer)

	// Assert
	assert.NoError(t, err)
}

func TestGetLayersHandler_GetAllReturnsError_WritesError(t *testing.T) {
	// Arrange
	r, _ := http.NewRequest("GET", "/worlds/13/layers", nil)
	f := layerHandlerFixture{}
	f.Setup(t, r)
	defer f.Teardown(t)

	f.layers.On("GetAll").Return(nil, errors.New(""))
	f.renderer.On("WriteError", f.writer, mock.Anything)
	handler := getLayersHandler{}

	// Act
	err := httpext.InvokeHandler(&handler, "GET", "/worlds/{worldID}/layers",
		f.writer, r, f.tx, f.renderer)

	// Assert
	assert.Error(t, err)
	f.layers.AssertExpectations(t)
	f.renderer.AssertExpectations(t)
}

func TestGetLayersHandler_GetAllReturnsLayers_WritesResponse(t *testing.T) {
	// Arrange
	r, _ := http.NewRequest("GET", "/worlds/13/layers", nil)
	f := layerHandlerFixture{}
	f.Setup(t, r)
	defer f.Teardown(t)

	layers := []*db.Layer{&db.Layer{}, &db.Layer{}}
	f.layers.On("GetAll").Return(layers, nil)
	f.renderer.On("WriteObject", f.writer, 200, layers)
	handler := getLayersHandler{}

	// Act
	err := httpext.InvokeHandler(&handler, "GET", "/worlds/{worldID}/layers",
		f.writer, r, f.tx, f.renderer)

	// Assert
	assert.NoError(t, err)
	f.layers.AssertExpectations(t)
	f.renderer.AssertExpectations(t)
}

func TestGetLayerHander_GetAllReturnsError_WritesError(t *testing.T) {
	// Arrange
	r, _ := http.NewRequest("GET", "/worlds/13/layers/42", nil)
	f := layerHandlerFixture{}
	f.Setup(t, r)
	defer f.Teardown(t)

	f.layers.On("Get", int64(42)).Return(nil, errors.New(""))
	f.renderer.On("WriteError", f.writer, mock.Anything)
	handler := getLayerHandler{}

	// Act
	err := httpext.InvokeHandler(&handler, "GET", "/worlds/{worldID}/layers/{layerID}",
		f.writer, r, f.tx, f.renderer)

	// Assert
	assert.Error(t, err)
	f.layers.AssertExpectations(t)
	f.renderer.AssertExpectations(t)
}

func TestGetLayerHandler_GetAllReturnsLayers_WritesResponse(t *testing.T) {
	// Arrange
	r, _ := http.NewRequest("GET", "/worlds/13/layers/42", nil)
	f := layerHandlerFixture{}
	f.Setup(t, r)
	defer f.Teardown(t)

	f.layers.On("Get", int64(42)).Return(&db.Layer{}, nil)
	f.renderer.On("WriteObject", f.writer, 200, mock.Anything)
	handler := getLayerHandler{}

	// Act
	err := httpext.InvokeHandler(&handler, "GET", "/worlds/{worldID}/layers/{layerID}",
		f.writer, r, f.tx, f.renderer)

	// Assert
	assert.NoError(t, err)
	f.layers.AssertExpectations(t)
	f.renderer.AssertExpectations(t)
}

func TestPostLayerHandler_InvalidBody_WritesError(t *testing.T) {
	// Arrange
	buffer := bytes.NewBuffer([]byte("{}"))
	r, _ := http.NewRequest("POST", "/worlds/42/layers", buffer)
	f := layerHandlerFixture{}
	f.Setup(t, r)
	defer f.Teardown(t)

	f.renderer.On("WriteError", f.writer, mock.Anything)
	handler := postLayerHandler{}

	// Act
	err := httpext.InvokeHandler(&handler, "POST", "/worlds/{worldID}/layers",
		f.writer, r, f.tx, f.renderer)

	// Assert
	assert.Error(t, err)
	f.renderer.AssertExpectations(t)
}

func TestPostLayerHandler_ValidBody_AddsToDatabaseAndWritesObject(t *testing.T) {
	// Arrange
	buffer := bytes.NewBuffer([]byte(`{"name":"MyLayer"}`))
	r, _ := http.NewRequest("POST", "/worlds/42/layers", buffer)
	f := layerHandlerFixture{}
	f.Setup(t, r)
	defer f.Teardown(t)

	f.layers.On("Add", mock.Anything).Return(int64(1), nil)
	f.renderer.On("WriteObject", f.writer, 200, mock.Anything)
	handler := postLayerHandler{}

	// Act
	err := httpext.InvokeHandler(&handler, "POST", "/worlds/{worldID}/layers",
		f.writer, r, f.tx, f.renderer)

	// Assert
	assert.NoError(t, err)
	f.layers.AssertExpectations(t)
	f.renderer.AssertExpectations(t)
}
