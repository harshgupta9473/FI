package middleware

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/harshgupta9473/fi/utils"
)

type Claims struct {
	UserName string `json:"username"`
	jwt.RegisteredClaims
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Missing Authorization header"})
			return
		}
		parts := strings.Split(tokenHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid Authorization header format"})
			return
		}
		token:=parts[1]
		claim, err := validateJWTToken(token)
		if err != nil {
			utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid or expired token"})
			return
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
