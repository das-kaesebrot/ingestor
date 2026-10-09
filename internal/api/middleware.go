package api

import (
	"encoding/json/v2"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

func LogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next.ServeHTTP(w, r)
		elapsedTime := time.Since(startTime)
		slog.Debug("request handled", "client", r.RemoteAddr, "method", r.Method, "path", r.URL.Path, "elapsed", elapsedTime)
	})
}

func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}
			if rec == http.ErrAbortHandler {
				panic(rec)
			}

			slog.Error("panic recovered",
				"panic", rec,
				"method", r.Method,
				"path", r.URL.Path,
				"stack", string(debug.Stack()),
			)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.MarshalWrite(w, InternalServerError{
				Title:  "An internal server error has occured.",
				Status: http.StatusInternalServerError,
			})
		}()
		next.ServeHTTP(w, r)
	})
}
