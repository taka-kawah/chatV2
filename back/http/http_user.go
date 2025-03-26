package http

import (
	"back/provider"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func setUpUserEndpointV1(r *gin.RouterGroup, p provider.ProviderManager) {
	ur := r.Group("/user")

	ur.GET("/all", p.Ap().ValidateToken(), func(ctx *gin.Context) {
		users, ce := p.Up().GetAllUsers()
		if ce != nil {
			internalErrorRes(ce, ctx)
		}
		ctx.JSON(http.StatusOK, users)
	})

	ur.GET("/", p.Ap().ValidateToken(), func(ctx *gin.Context) {
		email, ok := ctx.GetQuery("email")
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "email is not set"})
			return
		}
		user, ce := p.Up().GetFromEmail(email)
		if ce != nil {
			internalErrorRes(ce, ctx)
			return
		}
		if user == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "user not created"})
			return
		}
		ctx.JSON(http.StatusOK, user)
	})

	ur.POST("/new", p.Ap().ValidateToken(), func(ctx *gin.Context) {
		var req userInfoCreation
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		if err := p.Up().RegisterAccount(req.Name, req.Email); err != nil {
			internalErrorRes(err, ctx)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "successfully created!"})
	})

	ur.POST("/name/change", p.Ap().ValidateToken(), func(ctx *gin.Context) {
		var req userInfoChange
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		if err := p.Up().UpdateName(req.Id, req.NewName); err != nil {
			internalErrorRes(err, ctx)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "successfully updated!"})
	})

	ur.GET("/delete", p.Ap().ValidateToken(), func(ctx *gin.Context) {
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
		if ce := p.Up().Delete(uint(id)); ce != nil {
			internalErrorRes(ce, ctx)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "successfully deleted"})
	})
}

type userInfoCreation struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type userInfoChange struct {
	Id      uint   `json:"id"`
	NewName string `json:"newName"`
}
