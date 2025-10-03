package route

import (
	"github.com/elisha1995/magic-stream-movies/server/controller"
	"github.com/elisha1995/magic-stream-movies/server/middleware"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetupProtectedRoutes(router *gin.Engine, client *mongo.Client) {
	router.Use(middleware.AuthMiddleWare())

	router.GET("/movies/:imdb_id", controller.GetMovie(client))
	router.POST("/movies", controller.AddMovie(client))
	router.PATCH("/update-review/:imdb_id", controller.AdminReviewUpdate(client))
	router.GET("/recommended-movies", controller.GetRecommendedMovies(client))
}
