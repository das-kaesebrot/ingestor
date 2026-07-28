package handler

import (
	"net/http"

	"dev.kaesebrot.eu/go/ingestor/internal/renderer"
)

type Handler struct {
	renderer *renderer.Renderer
}

func NewHandler(renderer *renderer.Renderer) *Handler {
	return &Handler{renderer: renderer}
}

func (h *Handler) GetRoot(w http.ResponseWriter, r *http.Request) error {
	return h.renderer.RenderWithoutData(w, "home")
}
