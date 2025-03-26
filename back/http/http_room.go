package http

import (
	"back/provider"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func setUpRoomEndpointV1(r *gin.RouterGroup, p provider.ProviderManager) {
	rr := r.Group("/room")
	rr.POST("/new", p.Ap().ValidateToken(), func(ctx *gin.Context) {
		var req newRoom
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		if err := p.Rp().CreateNewRoom(req.Name); err != nil {
			internalErrorRes(err, ctx)
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{"message": "successfully created"})
	})

	rr.GET("/all", p.Ap().ValidateToken(), func(ctx *gin.Context) {
		rooms, err := p.Rp().GetAllRooms()
		if err != nil {
			internalErrorRes(err, ctx)
		}
		ctx.JSON(http.StatusOK, rooms)
	})

	rr.GET("/", p.Ap().ValidateToken(), func(ctx *gin.Context) {
		strId, ok := ctx.GetQuery("id")
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is not set"})
		}
		id, err := strconv.ParseUint(strId, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		room, ce := p.Rp().GetRoomById(uint(id))
		if ce != nil {
			internalErrorRes(ce, ctx)
			return
		}
		ctx.JSON(http.StatusOK, room)
	})

	rr.POST("/name/change", p.Ap().ValidateToken(), func(ctx *gin.Context) {
		var req nameChange
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		if err := p.Rp().UpdateRoomName(req.Id, req.NewName); err != nil {
			internalErrorRes(err, ctx)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "successfully updated"})
	})
}

type newRoom struct {
	Name string `json:"name"`
}

type nameChange struct {
	Id      uint   `json:"id"`
	NewName string `json:"newName"`
}
