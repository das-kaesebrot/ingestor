package api

import (
	"encoding/json/v2"
	"net/http"

	"dev.kaesebrot.eu/go/ingestor/internal/repository"
)

type APIHandler struct {
	repo   *repository.Repository
	prefix string
}

func NewAPIHandler(db *repository.Repository, prefix string) *APIHandler {
	return &APIHandler{repo: db, prefix: prefix}
}

func (h *APIHandler) APIMux() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", MakeMiddleware(h.GetPing))
	mux.HandleFunc("POST /projects", MakeMiddleware(h.PostCreateNewProject))
	mux.HandleFunc("GET /projects/{token}", MakeMiddleware(h.GetProject))
	mux.HandleFunc("GET /projects/admin/{token}", MakeMiddleware(h.GetProjectByAdminToken))

	return http.StripPrefix(h.prefix, mux)
}

func (h *APIHandler) GetPing(w http.ResponseWriter, r *http.Request) error {
	return writeJSON(w, map[string]string{"response": "pong"})
}

func (h *APIHandler) SearchAllProjectsPaged(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *APIHandler) GetProject(w http.ResponseWriter, r *http.Request) error {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return err
	}

	var project repository.Project
	if result := h.repo.DB.First(&project, id); result.Error != nil {
		return err
	}

	return writeJSON(w, ProjectResponseFromProject(project, false))
}

func (h *APIHandler) PostCreateNewProject(w http.ResponseWriter, r *http.Request) error {
	var parsedReq CreateProjectRequest
	if err := json.UnmarshalRead(r.Body, parsedReq); err != nil {
		return err
	}

	project := new(repository.Project{})
	project.Description = parsedReq.Description

	err := h.repo.Save(project)

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
	_, err = w.Write(jData)
	return err
}
