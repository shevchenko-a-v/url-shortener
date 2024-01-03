package middleware

import (
	"log/slog"
	"net/http"
	"strings"
)

type AuthorizationChecker interface {
	AreCredentialsValid(string, string) bool
}

func BasicAuth(authorizationChecker AuthorizationChecker, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()

		if r.Method != http.MethodPost || !strings.HasPrefix(r.URL.String(), "/save") || (ok && authorizationChecker.AreCredentialsValid(username, password)) {
			slog.Debug("credentials are valid or not required")
			next.ServeHTTP(w, r)
			return
		}

		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
