package handler

import (
	"net/http"

	"dev.kaesebrot.eu/go/ingestor/internal/renderer"
	"dev.kaesebrot.eu/go/ingestor/internal/repository"
)

type Handler struct {
	renderer *renderer.Renderer
	db       *repository.Repository
}

func NewHandler(renderer *renderer.Renderer, db *repository.Repository) *Handler {
	return &Handler{renderer: renderer, db: db}
}

func (h *Handler) GetRoot(w http.ResponseWriter, r *http.Request) error {
	return h.renderer.RenderWithoutData(w, "home")
}
