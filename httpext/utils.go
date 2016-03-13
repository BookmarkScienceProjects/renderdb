package httpext

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

// ReadInt64ID returns an int64 ID read from the mux vars (URL vars).
// vars is generated from a request by using 'mux.Vars(r)'.
// If the operation fails a HttpError with status code http.StatusBadRequest
// is returned.
func ReadInt64ID(vars map[string]string, fieldName string) (int64, HttpError) {
	strValue, ok := vars[fieldName]
	if !ok {
		return -1, NewHttpError(fmt.Errorf("'%s' is not set", fieldName), http.StatusBadRequest)
	}

	id, err := strconv.ParseInt(strValue, 10, 64)
	if err != nil {
		return -1, NewHttpError(fmt.Errorf("Expected numeric field '%s', but got '%s'", fieldName, strValue), http.StatusBadRequest)
	}

	return id, nil
}

var invokeCount int

// InvokeHandler invokes the handler given. This test should be used in unit tests only.
func InvokeHandler(h Handler, method, path string,
	w http.ResponseWriter, r *http.Request,
	tx *sqlx.Tx, renderer ResponseRenderer) error {

	// Add a new sub-path for each invocation since
	// we cannot (easily) remove old handler
	invokeCount++
	router := mux.NewRouter()
	http.Handle(fmt.Sprintf("/%d", invokeCount), router)

	var err error
	router.Methods(method).Path(path).
		HandlerFunc(
			func(w http.ResponseWriter, _ *http.Request) {
				err = h.Handle(tx, renderer, w, r)
			})

		// Modify the request to add "/%d" to the request-URL
	r.URL.RawPath = fmt.Sprintf("/%d%s", invokeCount, r.URL.RawPath)
	router.ServeHTTP(w, r)
	return err
}
