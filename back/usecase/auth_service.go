package usecase

import (
	"back/db"
	"back/middleware/authentication"
	"back/provider"
	"fmt"

	"github.com/gin-gonic/gin"
)

type AuthService struct {
	d *db.AuthDriver
	m *authentication.AuthMiddleware
}

func NewAuthService(d *db.AuthDriver, m *authentication.AuthMiddleware) provider.AuthProvider {
	return &AuthService{d: d, m: m}
}

func (as *AuthService) SignUp(email string, hashedPassword string) provider.CustomError {
	if err := as.d.Create(email, hashedPassword); err != nil {
		return &authServiceError{msg: "failed to create auth record", err: err}
	}
	return nil
}

func (as *AuthService) SignIn(email string, hashedPassword string) (string, provider.CustomError) {
	auth, err := as.d.CheckIfExist(email, hashedPassword)
	if err != nil {
		return "", &authServiceError{msg: "failed to fetch auth record", err: err}
	}
	if auth == nil {
		return "", &authServiceError{msg: "fetched auth was nil", err: fmt.Errorf("email: %v", email)}
	}
	tokenString, err := as.m.GenerateToken(auth.ID)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (as *AuthService) ValidateToken() gin.HandlerFunc {
	return as.m.AuthMW()
}

type authServiceError struct {
	msg string
	err error
}

func (e *authServiceError) Error() string {
	return fmt.Sprintf("error occured in auth service %v (%v)", e.msg, e.err)
}

func (e *authServiceError) Unwrap() error {
	return e.err
}
