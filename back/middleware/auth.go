package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func AuthMW() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			ctx.Abort()
			return
		}
		if !strings.HasPrefix(tokenString, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			ctx.Abort()
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		err := validateToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func GenerateToken(id uint) (string, error) {
	key, err := loadSecretKey()
	if err != nil {
		return "", &AuthMWError{msg: "failed to load secret key", err: err}
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  strconv.FormatUint(uint64(id), 10),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", &AuthMWError{msg: "failed to generate token", err: err}
	}
	return tokenString, nil
}

func validateToken(tokenString string) error {
	key, err := loadSecretKey()
	if err != nil {
		return &AuthMWError{msg: "failed to load secret key", err: err}
	}

	_, err = jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return &AuthMWError{msg: "failed to validate error", err: err}
	}
	return nil
}

func loadSecretKey() (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", fmt.Errorf("failed to load dotenv %v", err)
	}
	key := os.Getenv("SECRET_KEY")
	if key == "" {
		return "", errors.New("secret key not set")
	}
	return key, nil
}

type AuthMWError struct {
	msg string
	err error
}

func (e *AuthMWError) Error() string {
	return fmt.Sprintf("error occurred in auth middleware %v (%v)", e.msg, e.err)
}

func (e *AuthMWError) Unwrap() error {
	return e.err
}
