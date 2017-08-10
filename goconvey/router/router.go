package router

import "github.com/gin-gonic/gin"

// Define a single route. i.e route path, method type
// handler and api version number
type Route struct {
	Path    string
	Method  string
	Handler func(*gin.Context)
	Version string
}

// A slice of routes collection
type Routes []Route

// Create a new gin engine instance
func NewRouter() *gin.Engine {
	router := gin.Default()

	gin.SetMode(gin.ReleaseMode)

	return router
}
