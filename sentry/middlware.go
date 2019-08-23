package sentry

import (
	sdk "github.com/getsentry/sentry-go"
	"net/http"
)

// Record http request data into sentry scope
func LogResponseMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hub := sdk.CurrentHub()
		hub.Scope().SetRequest(sdk.Request{}.FromHTTPRequest(r))
		handler.ServeHTTP(w, r)
	}
}
