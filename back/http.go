package main

import (
	"github.com/gin-gonic/gin"
)

func setupRouter(p providerManager) *gin.Engine {
	e := gin.Default()
	r := e.Group("/api/v1")

	setUpSignEndPointV1(r, p)
	setUpUserEndpointV1(r, p)
	setUpRoomEndpointV1(r, p)
	setUpChatEndpointV1(r, p)

	return e
}
