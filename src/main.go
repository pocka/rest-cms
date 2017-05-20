package main

import (
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
)

func main(){
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong!")
	})

	router.Run(":80")
}
