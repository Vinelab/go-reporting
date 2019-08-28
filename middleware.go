package reporting

import "net/http"

func LogPanicMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer LogPanic()
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
