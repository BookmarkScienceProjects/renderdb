package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

// StaticController hosts static content.
type StaticController struct{}

// Init registers handler for the /static/*-route.
func (c *StaticController) Init(router *mux.Router) {
	router.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
}
