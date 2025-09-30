package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	controller "github.com/elisha1995/magic-stream-movies/server/controller"
)

func main() {
	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, MagicStreamMovies")
	})

	router.GET("/movies", controller.GetMovies())

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server: ", err)
	}
}
