package middlewares

import (
	"net/http"
	"strings"
)

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем конкретные домены (или * для всех в разработке)
		origin := r.Header.Get("Origin")
		allowedOrigins := []string{"http://localhost:3000", "http://localhost:3000"}

		allowOrigin := ""
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				allowOrigin = origin
				break
			}
		}

		if allowOrigin == "" && strings.HasPrefix(origin, "http://localhost:") {
			// Разрешаем все локальные хосты для разработки
			allowOrigin = origin
		}

		if allowOrigin != "" {
			w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

			// Обрабатываем preflight OPTIONS запрос
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
