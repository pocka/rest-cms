package main

import (
	"testing"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	gofight.New().
		GET("/ping").
		Run(NewServer(), func(res gofight.HTTPResponse, req gofight.HTTPRequest) {
			assert.Equal(t, "pong!", res.Body.String())
			assert.Equal(t, 200, res.Code)
		})
}
