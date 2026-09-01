package api

import (
	"encoding/json/v2"
	"net/http"

	"dev.kaesebrot.eu/go/ingestor/internal/repository"
)

type APIHandler struct {
	db *repository.Repository
}

func NewAPIHandler(db *repository.Repository) *APIHandler {
	return &APIHandler{db: db}
}

func (h *APIHandler) PostCreateNewProject(w http.ResponseWriter, r *http.Request) error {
	// project := new(repository.Project)
	return nil
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
