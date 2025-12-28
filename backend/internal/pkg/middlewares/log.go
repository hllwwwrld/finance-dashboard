package middlewares

import (
	"bytes"
	"log/slog"
	"net/http"
	"time"
)

// responseWriter - обертка над http.ResponseWriter для перехвата статус кода, размера и тела ответа
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.body != nil {
		rw.body.Write(b)
	}
	size, err := rw.ResponseWriter.Write(b)
	return size, err
}

// Log - middleware для логирования HTTP запросов и ответов
func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			slog.Info("HTTP request", r.Method, r.URL)

			start := time.Now()

			wrapped := &responseWriter{
				ResponseWriter: w,

				// сделал обертку и переопределил базовые методы, чтобы логировать себе код и бади
				statusCode: http.StatusOK,
				body:       &bytes.Buffer{},
			}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)

			logAttrs := []any{
				"method", r.Method,
				"path", r.URL.Path,
				"query", r.URL.RawQuery,
				"status", wrapped.statusCode,
				"duration_ms", duration.Milliseconds(),
				"response_body", wrapped.body.String(),
			}

			slog.Info("HTTP response", logAttrs...)
		},
	)
}
