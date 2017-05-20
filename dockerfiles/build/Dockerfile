FROM golang:1.8.1-alpine

MAINTAINER pocka <pockawoooh@gmail.com>

RUN apk --update add \
        git gcc musl-dev

RUN go get gopkg.in/gin-gonic/gin.v1 && \
	go get github.com/appleboy/gofight && \
	go get github.com/stretchr/testify/assert
