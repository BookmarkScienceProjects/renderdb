package httpext

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ResponseRenderer interface {
	WriteEmpty(w http.ResponseWriter, statusCode int)
	WriteObject(w http.ResponseWriter, statusCode int, val interface{})
	WriteError(w http.ResponseWriter, err error)
}

type jsonResponseRenderer struct {
}

// NewJSONResponseRenderer creates a renderer that outputs JSON.
func NewJSONResponseRenderer() ResponseRenderer {
	return &jsonResponseRenderer{}
}

func (r *jsonResponseRenderer) WriteEmpty(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}

func (r *jsonResponseRenderer) WriteObject(w http.ResponseWriter, statusCode int, val interface{}) {
	w.WriteHeader(statusCode)
	buffer, err := json.Marshal(val)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Could not marshal value '%+v'", val)))
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.Write(buffer)

	}
}

func (r *jsonResponseRenderer) WriteError(w http.ResponseWriter, err error) {
	var httpError HttpError
	var ok bool
	if httpError, ok = err.(HttpError); !ok {
		httpError = NewHttpError(err, http.StatusInternalServerError)
	}
	r.WriteObject(w, httpError.StatusCode(), httpError)
}
