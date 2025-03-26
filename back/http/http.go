package http

import (
	"back/provider"

	"github.com/gin-gonic/gin"
)

func SetupRouter(p provider.ProviderManager) *gin.Engine {
	e := gin.Default()
	r := e.Group("/api/v1")

	setUpAuthEndPointV1(r, p)
	setUpUserEndpointV1(r, p)
	setUpRoomEndpointV1(r, p)
	setUpChatEndpointV1(r, p)

	return e
}
