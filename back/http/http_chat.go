package http

import (
	"back/provider"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func setUpChatEndpointV1(r *gin.RouterGroup, p provider.ProviderManager) {
	cr := r.Group("/chat")

	cr.POST("/new", p.Ap().ValidateToken(), func(ctx *gin.Context) {
		var req newChat
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		if err := p.Cp().PostChat(req.Message, req.UserId, req.RoomId); err != nil {
			internalErrorRes(err, ctx)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "successfully created"})
	})

	cr.GET("/", p.Ap().ValidateToken(), func(ctx *gin.Context) {
		strId, ok := ctx.GetQuery("id")
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "id not set"})
		}
		id, err := strconv.ParseUint(strId, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		}
		chats, ce := p.Cp().GetRecentChatsFromOneRoom(uint(id))
		if ce != nil {
			internalErrorRes(ce, ctx)
		}
		ctx.JSON(http.StatusOK, chats)
	})
}

type newChat struct {
	Message string `json:"message"`
	UserId  uint   `json:"userId"`
	RoomId  uint   `json:"roomId"`
}
