package middleware

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
	"time"
)

type Claims struct {
	UserName string `json:"username"`
	jwt.RegisteredClaims
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {

		}
		claim, err := validateJWTToken(tokenHeader)
		if err != nil {

		}
		ctx := context.WithValue(r.Context(), "userID", claim.UserName)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CreateJWTToken(username string) (string, error) {
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	claims := Claims{
		UserName: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func validateJWTToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {

	}
	if !token.Valid {
		return claims, errors.New("invalid token")
	}
	return claims, nil
}
