package api

import (
	"encoding/json/v2"
	"net/http"

	"dev.kaesebrot.eu/go/ingestor/internal/middleware"
	"dev.kaesebrot.eu/go/ingestor/internal/repository"
)

type APIHandler struct {
	db     *repository.Repository
	prefix string
}

func NewAPIHandler(db *repository.Repository, prefix string) *APIHandler {
	return &APIHandler{db: db, prefix: prefix}
}

func (h *APIHandler) APIMux() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", middleware.Make(h.GetPing))
	mux.HandleFunc("POST /project", middleware.Make(h.PostCreateNewProject))

	return http.StripPrefix(h.prefix, mux)
}

func (h *APIHandler) GetPing(w http.ResponseWriter, r *http.Request) error {
	return writeJSON(w, map[string]string{"response": "pong"})
}

func (h *APIHandler) SearchAllProjectsPaged(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *APIHandler) PostCreateNewProject(w http.ResponseWriter, r *http.Request) error {
	project := new(repository.Project{})

	err := h.db.Save(project)

	if err != nil {
		return err
	}

	return writeJSON(w, ProjectResponseFromProject(*project, true))
}

func writeJSON(w http.ResponseWriter, object any) error {
	jData, err := json.Marshal(object)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
	return nil
}
