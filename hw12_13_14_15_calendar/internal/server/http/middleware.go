package internalhttp

import (
	"net/http"
	"time"

	"github.com/urfave/negroni"
)

func loggingMiddleware(next http.Handler, logg Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		lrw := negroni.NewResponseWriter(w)
		next.ServeHTTP(lrw, r)

		latency := time.Since(startTime).Milliseconds()
		logg.Infof(
			"%s [%s] %s %s %s %d %d \"%s\"",
			r.RemoteAddr,
			startTime.String(),
			r.Method,
			r.RequestURI,
			r.Proto,
			lrw.Status(),
			latency,
			r.Header.Get("User-Agent"),
		)
	})
}
