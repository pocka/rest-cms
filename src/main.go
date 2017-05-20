package main

import (
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
)

func main(){
	router := gin.New()

	// Specify middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Bind routes
	router.GET("/ping", ping)

	// Run server
	router.Run(":80")
}

// Health check api
func ping(c *gin.Context) {
	c.String(http.StatusOK, "pong!")
}
