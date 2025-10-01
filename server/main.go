package main

import (
	"fmt"
	"net/http"

	"github.com/elisha1995/magic-stream-movies/server/route"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, MagicStreamMovies")
	})

	route.SetupUnProtectedRoutes(router)
	route.SetupProtectedRoutes(router)

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server: ", err)
	}
}
