package usecase

import (
	"back/db"
	"back/interfaces"
	"back/middleware"
	"fmt"
)

type IAuthService interface {
	SignUp(email string, hashedPassword string) interfaces.CustomError
	SignIn(id string, hashedPassword string) (string, error)
}

type AuthService struct {
	repo *db.AuthDriver
}

func NewAuthService(repo *db.AuthDriver) *AuthService {
	return &AuthService{repo: repo}
}

func (as *AuthService) SignUp(email string, hashedPassword string) interfaces.CustomError {
	if err := as.repo.Create(email, hashedPassword); err != nil {
		return &AuthServiceError{msg: "failed to create auth record", err: err}
	}
	return nil
}

func (as *AuthService) SignIn(email string, hashedPassword string) (string, error) {
	auth, err := as.repo.CheckIfExist(email, hashedPassword)
	if err != nil {
		return "", &AuthServiceError{msg: "failed to fetch auth record", err: err}
	}
	if auth == nil {
		return "", &AuthServiceError{msg: "fetched auth was nil", err: fmt.Errorf("email: %v", email)}
	}
	tokenString, err := middleware.GenerateToken(auth.ID)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

type AuthServiceError struct {
	msg string
	err error
}

func (e *AuthServiceError) Error() string {
	return fmt.Sprintf("error occured in auth service %v (%v)", e.msg, e.err)
}

func (e *AuthServiceError) Unwrap() error {
	return e.err
}
