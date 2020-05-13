package utils

import (
	"context"
	"net/http"
)

func LoggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("x-request-id")

		ctx := context.WithValue(r.Context(), "x-request-id", id)
		r = r.WithContext(ctx)

		InfoContext(r.Context(), "Request [%s] %q", r.Method, r.URL.String())
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
