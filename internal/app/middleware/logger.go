package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entry := slog.With(
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
			slog.String("user_agent", r.UserAgent()),
			slog.String("request_id", GetRequestId(r.Context())),
		)
		start := time.Now()

		next.ServeHTTP(w, r)

		dur := time.Since(start)
		entry.Debug(fmt.Sprintf("request handled in %dms", dur.Milliseconds()))
	})
}
