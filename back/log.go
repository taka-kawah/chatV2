package main

import (
	"back/interfaces"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func errLog(e interfaces.CustomError) {
	log.Fatalf("%v, %v", e.Error(), e.Unwrap().Error())
}

func internalErrorRes(e interfaces.CustomError, ctx *gin.Context) {
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": e.Error(), "internal error": e.Unwrap().Error()})
}

func marshalErrorRes(e error, ctx *gin.Context) {
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": "failed to marshal response", "internal error": e,
	})
}
