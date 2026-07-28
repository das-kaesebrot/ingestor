//go:generate go run tools/setup-web.go

package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log/slog"
	"math"
	"net/http"
	"os"

	"dev.kaesebrot.eu/go/ingestor/internal/handler"
	"dev.kaesebrot.eu/go/ingestor/internal/middleware"
	"dev.kaesebrot.eu/go/ingestor/internal/renderer"
	"dev.kaesebrot.eu/go/ingestor/internal/utility"
)

var (
	Version = "v0.0.1-dev"
	GitHash = "0000000000000000000000000000000000000000"
)

//go:embed web
var webFS embed.FS

var templateFilesRoot = "web/template"
var staticFilesRoot = "web/static"
var webStaticFilesRoot = "static"

func main() {
	var logLevelStr, host string
	var portUnparsed int

	flag.StringVar(&logLevelStr, "l", "", "Log level")
	flag.StringVar(&logLevelStr, "loglevel", "", "Log level")
	flag.StringVar(&host, "host", "[::]", "HTTP server host")
	flag.IntVar(&portUnparsed, "p", 8000, "HTTP server port")
	flag.IntVar(&portUnparsed, "port", 8000, "HTTP server port")

	flag.Parse()

	if l := os.Getenv("LOG_LEVEL"); l != "" {
		logLevelStr = l
	}

	if portUnparsed < 0 || portUnparsed > math.MaxUint16 {
		utility.HandleErr("Error while parsing port value!", fmt.Errorf("Port value has to be 0 < value <= %d", math.MaxUint16))
	}

	port := uint16(portUnparsed)

	logLevel, err := utility.ParseLogLevelFromString(logLevelStr)
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)
	if err != nil {
		slog.Warn(err.Error())
	}

	slog.Info("Starting up", "version", Version, "gitHash", GitHash)
	slog.Debug("Using slog with specified level", "loglevel", logLevel)

	renderer, err := renderer.New(webFS, staticFilesRoot, webStaticFilesRoot, templateFilesRoot, ".tmpl", map[string]any{
		"StaticLibsSubDir": "/" + webStaticFilesRoot + "/libs",
	})
	if err != nil {
		utility.HandleErr("Error while creating renderer", err)
	}

	h := handler.NewHandler(renderer)
	mux := http.NewServeMux()

	staticFS, err := fs.Sub(webFS, staticFilesRoot)
	if err != nil {
		utility.HandleErr("Error while creating sub FS for static file server", err)
	}
	mux.Handle(fmt.Sprintf("GET /%s/", webStaticFilesRoot), http.StripPrefix("/"+webStaticFilesRoot+"/", http.FileServerFS(staticFS)))
	mux.HandleFunc("GET /", middleware.Make(h.GetRoot))

	slog.Info("Server ready", "host", host, "port", port)
	http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), mux)
}
