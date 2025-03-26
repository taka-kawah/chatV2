package http

import (
	"back/provider"
	"net/http"

	"github.com/gin-gonic/gin"
)

func internalErrorRes(e provider.CustomError, ctx *gin.Context) {
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": e.Error(), "internal error": e.Unwrap().Error()})
}
