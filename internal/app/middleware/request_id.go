package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type ContextKey string

const RequestIdKey ContextKey = "request_id"

func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := uuid.New()
		ctx = context.WithValue(ctx, RequestIdKey, id.String())
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func GetRequestId(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	ctxId := ctx.Value(RequestIdKey)
	id, ok := ctxId.(string)
	if !ok {
		slog.Error("wrong request id format")
		return ""
	}
	return id
}
