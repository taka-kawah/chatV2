package http

import (
	"back/provider"
	"net/http"

	"github.com/gin-gonic/gin"
)

func setUpAuthEndPointV1(r *gin.RouterGroup, p provider.ProviderManager) {
	ar := r.Group("/auth")

	ar.POST("/new", func(ctx *gin.Context) {
		var req signUpAndSignIn
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if ce := p.Ap().SignUp(req.Email, req.HashedPassword); ce != nil {
			internalErrorRes(ce, ctx)
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{"message": "new account is created!"})
	})

	ar.POST("/", func(ctx *gin.Context) {
		var req signUpAndSignIn
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		tk, ce := p.Ap().SignIn(req.Email, req.HashedPassword)
		if ce != nil {
			internalErrorRes(ce, ctx)
			return
		}
		ctx.SetCookie("token", tk, 100, "/", "localhost", false, true)
		ctx.JSON(http.StatusOK, gin.H{
			"message": "successfully signined!",
		})
	})
}

type signUpAndSignIn struct {
	Email          string `json:"email"`
	HashedPassword string `json:"hashedPassword"`
}
