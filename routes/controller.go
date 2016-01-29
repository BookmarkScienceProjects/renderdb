package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type RequestContext struct {
	Vars map[string]string
	Body []byte
}

type Controller struct{}

func (c *Controller) CreateContext(r *http.Request) (*RequestContext, HttpError) {
	var err error
	defer r.Body.Close()

	ctx := new(RequestContext)
	ctx.Body, err = ioutil.ReadAll(r.Body)
	ctx.Vars = mux.Vars(r)
	return ctx, EncapulateIfError(err, http.StatusBadRequest)
}

func (c *Controller) ParseBody(ctx *RequestContext, v interface{}) HttpError {
	err := json.Unmarshal(ctx.Body, &v)
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
