package main

import (
	"gopkg.in/gin-gonic/gin.v1"
)

func main(){
	// Run server
	NewServer().Run(":80")
}

// Creates and returns server
func NewServer() *gin.Engine {
	router := gin.New()

	// Specify middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Bind routes
	router.GET("/ping", ping)

	return router
}

// Health check api
func ping(c *gin.Context) {
	c.String(200, "pong!")
}
