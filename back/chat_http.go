package main

import (
	"back/interfaces"
	"net/http"

	"github.com/gin-gonic/gin"
)

func setupRouter(p providerManager) *gin.Engine {
	e := gin.Default()
	r := e.Group("/api/v1")

	r.POST("/SignUp", func(ctx *gin.Context) {
		var req signUpAndSignIn
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		if ce := p.ap.SignUp(req.email, req.hashedPassword); ce != nil {
			internalErrorRes(ce, ctx)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "new account is created!"})
	})

	r.GET("/SignIn", func(ctx *gin.Context) {
		var req signUpAndSignIn
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		tk, ce := p.ap.SignIn(req.email, req.hashedPassword)
		if ce != nil {
			internalErrorRes(ce, ctx)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "successfully signined!",
			"token":   tk,
		})
	})

	return e
}

func internalErrorRes(e interfaces.CustomError, ctx *gin.Context) {
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": e.Error(), "internal error": e.Unwrap().Error()})
}
