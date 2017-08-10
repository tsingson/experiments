package http

import (
	"github.com/gin-gonic/gin"
	"github.com/tsingson/test/goconvey/router"
	"github.com/tsingson/test/goconvey/routes"
)

// Starting http web server
// by the given port number
func StartWebServer(port string) {

	router := router.NewRouter()

	router = RegisterRoutes(router)

	router.Run(":" + port)
}

// Register all api routes
func RegisterRoutes(router *gin.Engine) *gin.Engine {
	for _, r := range routes.ApiRoutes {
		switch r.Method {
		case "GET":
			router.GET(r.Version+r.Path, r.Handler)
			break
		case "POST":
			router.POST(r.Version+r.Path, r.Handler)
			break
		default:
			router.GET(r.Path, r.Handler)
		}
	}
	return router
}
