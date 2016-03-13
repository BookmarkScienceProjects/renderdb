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

type sceneHandlerFixture struct {
	mockDB sqlmock.Sqlmock
	db     *sqlx.DB
	tx     *sqlx.Tx

	scenes *db.MockScenes

	writer   *httptest.ResponseRecorder
	renderer *httpext.MockResponseRenderer
}

func (f *sceneHandlerFixture) Setup(t *testing.T, r *http.Request) {
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

	f.scenes = &db.MockScenes{}
	context.Set(r, scenesDBKey, f.scenes)
}

func (f *sceneHandlerFixture) Teardown(t *testing.T) {
	assert.NoError(t, f.db.Close())
}

func TestScenesMiddleware_Success(t *testing.T) {
	// Arrange
	r, _ := http.NewRequest("GET", "/worlds/13/layers/42/scenes", nil)
	f := sceneHandlerFixture{}
	f.Setup(t, r)
	defer f.Teardown(t)
	middleware := scenesMiddleware{}

	// Act
	err := httpext.InvokeHandler(&middleware, "GET", "/worlds/{worldID}/layers/{layerID}/scenes",
		f.writer, r, f.tx, f.renderer)

	// Assert
	assert.NoError(t, err)
}

func TestGetScenesHandler_GetAllReturnsError_WritesError(t *testing.T) {
	// Arrange
	r, _ := http.NewRequest("GET", "/worlds/13/layers/42/scenes", nil)
	f := sceneHandlerFixture{}
	f.Setup(t, r)
	defer f.Teardown(t)

	f.scenes.On("GetAll").Return(nil, errors.New(""))
	f.renderer.On("WriteError", f.writer, mock.Anything)
	handler := getScenesHandler{}

	// Act
	err := httpext.InvokeHandler(&handler, "GET", "/worlds/{worldID}/layers/{layerID}/scenes",
		f.writer, r, f.tx, f.renderer)

	// Assert
	assert.Error(t, err)
	f.scenes.AssertExpectations(t)
	f.renderer.AssertExpectations(t)
}

func TestGetScenesHandler_GetAllReturnsLayers_WritesResponse(t *testing.T) {
	// Arrange
	r, _ := http.NewRequest("GET", "/worlds/13/layers/42/scenes", nil)
	f := sceneHandlerFixture{}
	f.Setup(t, r)
	defer f.Teardown(t)

	scenes := []*db.Scene{&db.Scene{}, &db.Scene{}}
	f.scenes.On("GetAll").Return(scenes, nil)
	f.renderer.On("WriteObject", f.writer, 200, scenes)
	handler := getScenesHandler{}

	// Act
	err := httpext.InvokeHandler(&handler, "GET", "/worlds/{worldID}/layers/{layerID}/scenes",
		f.writer, r, f.tx, f.renderer)

	// Assert
	assert.NoError(t, err)
	f.scenes.AssertExpectations(t)
	f.renderer.AssertExpectations(t)
}

func TestGetSceneHander_GetAllReturnsError_WritesError(t *testing.T) {
	// Arrange
	r, _ := http.NewRequest("GET", "/worlds/13/layers/42/scenes/88", nil)
	f := sceneHandlerFixture{}
	f.Setup(t, r)
	defer f.Teardown(t)

	f.scenes.On("Get", int64(88)).Return(nil, errors.New(""))
	f.renderer.On("WriteError", f.writer, mock.Anything)
	handler := getSceneHandler{}

	// Act
	err := httpext.InvokeHandler(&handler, "GET", "/worlds/{worldID}/layers/{layerID}/scenes/{sceneID}",
		f.writer, r, f.tx, f.renderer)

	// Assert
	assert.Error(t, err)
	f.scenes.AssertExpectations(t)
	f.renderer.AssertExpectations(t)
}

func TestGetSceneHandler_GetAllReturnsScenes_WritesResponse(t *testing.T) {
	// Arrange
	r, _ := http.NewRequest("GET", "/worlds/13/layers/42/scenes/88", nil)
	f := sceneHandlerFixture{}
	f.Setup(t, r)
	defer f.Teardown(t)

	f.scenes.On("Get", int64(88)).Return(&db.Scene{}, nil)
	f.renderer.On("WriteObject", f.writer, 200, mock.Anything)
	handler := getSceneHandler{}

	// Act
	err := httpext.InvokeHandler(&handler, "GET", "/worlds/{worldID}/layers/{layerID}/scenes/{sceneID}",
		f.writer, r, f.tx, f.renderer)

	// Assert
	assert.NoError(t, err)
	f.scenes.AssertExpectations(t)
	f.renderer.AssertExpectations(t)
}

func TestPostSceneHandler_InvalidBody_WritesError(t *testing.T) {
	// Arrange
	buffer := bytes.NewBuffer([]byte("{}"))
	r, _ := http.NewRequest("POST", "/worlds/42/layers/13/scenes", buffer)
	f := sceneHandlerFixture{}
	f.Setup(t, r)
	defer f.Teardown(t)

	f.renderer.On("WriteError", f.writer, mock.Anything)
	handler := postSceneHandler{}

	// Act
	err := httpext.InvokeHandler(&handler, "POST", "/worlds/{worldID}/layers/{layerID}/scenes",
		f.writer, r, f.tx, f.renderer)

	// Assert
	assert.Error(t, err)
	f.renderer.AssertExpectations(t)
}

func TestPostSceneHandler_ValidBody_AddsToDatabaseAndWritesObject(t *testing.T) {
	// Arrange
	buffer := bytes.NewBuffer([]byte(`{"name":"MyScene"}`))
	r, _ := http.NewRequest("POST", "/worlds/42/layers/13/scenes", buffer)
	f := sceneHandlerFixture{}
	f.Setup(t, r)
	defer f.Teardown(t)

	f.scenes.On("Add", mock.Anything).Return(int64(1), nil)
	f.renderer.On("WriteObject", f.writer, 200, mock.Anything)
	handler := postSceneHandler{}

	// Act
	err := httpext.InvokeHandler(&handler, "POST", "/worlds/{worldID}/layers/{layerID}/scenes",
		f.writer, r, f.tx, f.renderer)

	// Assert
	assert.NoError(t, err)
	f.scenes.AssertExpectations(t)
	f.renderer.AssertExpectations(t)
}
