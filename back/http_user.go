package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func setUpUserEndpointV1(r *gin.RouterGroup, p providerManager) {

	r.GET("/AllUserInfo", p.ap.ValidateToken(), func(ctx *gin.Context) {
		users, ce := p.up.GetAllUsers()
		if ce != nil {
			internalErrorRes(ce, ctx)
		}
		ctx.JSON(http.StatusOK, gin.H{"UsersInfo": users})
	})

	r.GET("/UserInfo", p.ap.ValidateToken(), func(ctx *gin.Context) {
		email, ok := ctx.GetQuery("email")

		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "email is not set"})
			return
		}
		user, ce := p.up.GetFromEmail(email)
		if ce != nil {
			internalErrorRes(ce, ctx)
			return
		}
		if user == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "user not created"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"MyInfo": user})
	})

	r.POST("/UserInfoCreation", p.ap.ValidateToken(), func(ctx *gin.Context) {
		var req userInfoCreation
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		if ce := p.up.RegisterAccount(req.name, req.email); ce != nil {
			internalErrorRes(ce, ctx)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "successfully created!"})
	})

	r.POST("/UserNameChange", p.ap.ValidateToken(), func(ctx *gin.Context) {
		var req userInfoChange
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		if ce := p.up.UpdateName(req.id, req.newName); ce != nil {
			internalErrorRes(ce, ctx)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "successfully updated!"})
	})

	r.POST("/UserDeletion", p.ap.ValidateToken(), func(ctx *gin.Context) {
		strId, ok := ctx.GetQuery("id")
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is not set"})
			return
		}
		id, err := strconv.ParseUint(strId, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		if ce := p.up.Delete(uint(id)); ce != nil {
			internalErrorRes(ce, ctx)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "successfully deleted"})
	})
}

type userInfoCreation struct {
	name  string
	email string
}

type userInfoChange struct {
	id      uint
	newName string
}
