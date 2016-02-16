package routes

import "github.com/go-martini/martini"

// StaticController hosts static content.
type StaticController struct{}

// Init registers handler for the /static/*-route.
func (c *StaticController) Init(router *martini.ClassicMartini) {
	router.Use(martini.Static("static/", martini.StaticOptions{Prefix: "/static/"}))
}
