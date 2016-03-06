package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

// StaticController hosts static content.
type StaticController struct{}

// NewStaticController registers handler for the /static/-route.
func NewStaticController(router *mux.Router) *StaticController {
	controller := StaticController{}
	router.PathPrefix("/static/").Handler(http.FileServer(http.Dir("./static/")))
	return &controller
}
