package provider

import "github.com/gin-gonic/gin"

type AuthProvider interface {
	SignUp(email string, hashedPassword string) CustomError
	SignIn(email string, hashedPassword string) (string, CustomError)
	ValidateToken() gin.HandlerFunc
}
