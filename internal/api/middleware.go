package api

import (
	"context"
	"encoding/json/v2"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	middleware "github.com/oapi-codegen/nethttp-middleware"
)

func OapiRequestValidatorErrorHandler(_ context.Context, err error, w http.ResponseWriter, _ *http.Request, opts middleware.ErrorHandlerOpts) {
	status := opts.StatusCode
	detail := err.Error()

	switch status {
	case http.StatusBadRequest:
		writeProblemJSON(w, status, BadRequest{
			Status: status,
			Title:  "The request was malformed or invalid.",
			AdditionalProperties: map[string]interface{}{
				"errors": collectBadRequestErrors(err),
			},
		})
	case http.StatusNotFound:
		writeProblemJSON(w, status, NotFound{Status: status, Title: "The requested object was not found."})
	default:
		status = http.StatusInternalServerError
		writeProblemJSON(w, status, InternalServerError{Status: status, Title: "An internal server error has occured.", Detail: &detail})
	}
}

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

			w.Header().Set("Content-Type", "application/problem+json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.MarshalWrite(w, InternalServerError{
				Title:  "An internal server error has occured.",
				Status: http.StatusInternalServerError,
			})
		}()
		next.ServeHTTP(w, r)
	})
}
