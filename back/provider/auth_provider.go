package provider

import "github.com/gin-gonic/gin"

type AuthProvider interface {
	SignUp(email string, hashedPassword string, token string) CustomError
	SignIn(email string, hashedPassword string) (string, CustomError)
	SetToken(email string, hashedPassword string, token string) CustomError
	ValidateToken() gin.HandlerFunc
}
