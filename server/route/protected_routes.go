package route

import (
	"github.com/elisha1995/magic-stream-movies/server/controller"
	"github.com/elisha1995/magic-stream-movies/server/middleware"
	"github.com/gin-gonic/gin"
)

func SetupProtectedRoutes(router *gin.Engine) {
	router.Use(middleware.AuthMiddleWare())

	router.GET("/movies/:imdb_id", controller.GetMovie())
	router.POST("/movies", controller.AddMovie())
}
