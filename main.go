//go:generate go run tools/setup-web.go

package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"math"
	"net/http"
	"os"
	"path"

	"dev.kaesebrot.eu/go/ingestor/internal/api"
	"dev.kaesebrot.eu/go/ingestor/internal/repository"
	"dev.kaesebrot.eu/go/ingestor/internal/utility"
)

var (
	Version = "v0.0.1-dev"
	GitHash = "0000000000000000000000000000000000000000"
)

//go:embed web/app/build
var webFS embed.FS

var frontendFilesRoot = "web/app/build"

func main() {
	var logLevelStr, host string
	var portUnparsed int
	var versionFlag bool

	flag.StringVar(&logLevelStr, "l", "", "Log level")
	flag.StringVar(&logLevelStr, "loglevel", "", "Log level")
	flag.StringVar(&host, "host", "[::]", "HTTP server host")
	flag.IntVar(&portUnparsed, "p", 8000, "HTTP server port")
	flag.IntVar(&portUnparsed, "port", 8000, "HTTP server port")
	flag.BoolVar(&versionFlag, "v", false, "print version information")
	flag.BoolVar(&versionFlag, "version", false, "print version information")
	flag.Parse()

	if versionFlag {
		fmt.Printf("%v\n", Version)
		return
	}

	log.Printf("Version: %v", Version)

	dbFile := path.Clean(os.Getenv("INGESTOR_DB_FILE"))
	if dbFile == "." {
		dbFile = "ingestor.db"
	}

	db, err := repository.New(dbFile)
	if err != nil {
		utility.HandleErr("Error initializing database", err)
	}

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

	apiPrefix := "/api/v1"
	a := api.NewAPIHandler(db, apiPrefix)
	mux := http.NewServeMux()

	staticFS, err := fs.Sub(webFS, frontendFilesRoot)
	if err != nil {
		utility.HandleErr("Error while creating sub FS for static file server", err)
	}
	mux.Handle(apiPrefix+"/", a.APIMux())
	mux.Handle("/", http.FileServerFS(staticFS))

	slog.Info("Server ready", "host", host, "port", port)
	http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), mux)
}
