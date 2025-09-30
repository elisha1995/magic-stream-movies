package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, MagicStreamMovies")
	})

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server: ", err)
	}
}
