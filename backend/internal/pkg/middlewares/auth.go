package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Claims структура для JWT токена
type Claims struct {
	jwt.RegisteredClaims

	UserID string `json:"userId"`
	Login  string `json:"login"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

const (
	userContextKey = "user"
	authCookieName = "auth_token"
)

var jwtSecret = os.Getenv("JWT_SECRET")

// Auth проверяет JWT токен из cookie
func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(authCookieName)
		if err != nil {
			http.Error(w, fmt.Sprintf("Auth.r.Cookie err: %v", err), http.StatusUnauthorized)
			return
		}

		claims, err := validateJWT(cookie.Value)
		if err != nil {
			http.Error(w, fmt.Sprintf("Auth.validateJWT err: %v", err), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// Генерация JWT токена (на этапе логина генерю куку авторизации)
func generateJWT(userID, login string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: userID,
		Login:  login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "go-auth-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// Валидация JWT токена (на запросах, которые разделются на конкретных юзеров проверяю куку авторизации и какому юзеру она принадлежит)
func validateJWT(tokenString string) (*Claims, error) {
	claims := new(Claims) // тож самое, что &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("validateJWT.jwt.ParseWithClaims mot ok: %v", token.Method.Alg())
			}

			return []byte(jwtSecret), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("validateJWT.jwt.ParseWithClaims err: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("validateJWT.jwt.Validate token is invalid")
	}

	return claims, nil
}

// Получение пользователя из контекста
func getUserFromContext(r *http.Request) *Claims {
	user, ok := r.Context().Value(userContextKey).(*Claims)
	if !ok {
		return nil
	}
	return user
}
