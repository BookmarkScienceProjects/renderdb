package routes

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type RequestContext struct {
	Vars map[string]string
	Body io.ReadCloser
}

type Controller struct{}

func (c *Controller) CreateContext(r *http.Request) (*RequestContext, HttpError) {
	ctx := new(RequestContext)
	ctx.Body = r.Body
	ctx.Vars = mux.Vars(r)
	return ctx, nil
}

func (c *Controller) ParseBody(ctx *RequestContext, v interface{}) HttpError {
	defer ctx.Body.Close()
	decoder := json.NewDecoder(ctx.Body)
	err := decoder.Decode(&v)
	return EncapulateIfError(err, http.StatusBadRequest)
}

func (c *Controller) HandleError(w http.ResponseWriter, err HttpError) {
	w.Header().Add("Content-Type", "application/json")
	type errWrapper struct {
		Error      string `json:"error"`
		StatusCode int    `json:"statusCode"`
	}
	wrapper := errWrapper{err.Error().Error(), err.StatusCode()}
	buffer, _ := json.Marshal(wrapper)
	http.Error(w, string(buffer), err.StatusCode())
}

func (c *Controller) WriteResponse(w http.ResponseWriter, response interface{}) {
	w.Header().Add("Content-Type", "application/json")
	buffer, _ := json.Marshal(response)
	w.Write(buffer)
}
