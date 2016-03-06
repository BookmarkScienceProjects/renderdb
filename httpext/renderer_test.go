package httpext

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONResponseRenderer_WriteEmpty_WritesStatusCode(t *testing.T) {
	// Arrange
	w := httptest.NewRecorder()
	renderer := NewJSONResponseRenderer()

	// Act
	renderer.WriteEmpty(w, http.StatusBadRequest)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJSONResponseRenderer_WriteError_HttpError_WritesBodyAndStatusCode(t *testing.T) {
	// Arrange
	w := httptest.NewRecorder()
	renderer := NewJSONResponseRenderer()

	// Act
	err := NewHttpError(errors.New("message"), http.StatusConflict)
	renderer.WriteError(w, err)

	// Assert
	assert.Equal(t, http.StatusConflict, w.Code)
	m := make(map[string]interface{})
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &m))
}

func TestJSONResponseRenderer_WriteError_StandardError_WritesBodyAndInternalError(t *testing.T) {
	// Arrange
	w := httptest.NewRecorder()
	renderer := NewJSONResponseRenderer()

	// Act
	renderer.WriteError(w, errors.New("message"))

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	m := make(map[string]interface{})
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &m))
}

func TestJSONResponseRenderer_WriteObject_CannotMarshal_WritesError(t *testing.T) {
	// Arrange
	w := httptest.NewRecorder()
	renderer := NewJSONResponseRenderer()

	// Act
	renderer.WriteObject(w, http.StatusOK, func() {})

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestJSONResponseRenderer_WriteObject_ValidObject_WritesBodyAndStatusCode(t *testing.T) {
	// Arrange
	w := httptest.NewRecorder()
	renderer := NewJSONResponseRenderer()

	// Act
	m := make(map[string]interface{})
	m["string"] = "value"
	m["int"] = 42
	renderer.WriteObject(w, http.StatusOK, m)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.HeaderMap.Get("Content-Type"))
	r := make(map[string]interface{})
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &r))
	assert.Equal(t, "value", m["string"])
	assert.Equal(t, 42, m["int"])
}
