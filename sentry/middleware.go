package sentry

import (
	sdk "github.com/getsentry/sentry-go"
	"net/http"
)

// Record http request data into sentry scope
func LogResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hub := sdk.CurrentHub()
		hub.Scope().SetRequest(sdk.Request{}.FromHTTPRequest(r))
		next.ServeHTTP(w, r)
	})
}
