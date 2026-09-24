package api

import (
	"time"
	"uuid"

	"dev.kaesebrot.eu/go/ingestor/internal/repository"
)

type ErrorResponse struct {
	Message   string         `json:"message"`
	Code      int            `json:"code"`
	Traceback string         `json:"traceback,omitempty"`
	Details   map[string]any `json:"details,omitempty"`
}

type PagedProjectResponse struct {
	Projects []ProjectResponse `json:"projects"`
}

type ProjectResponse struct {
	ID          uuid.UUID `json:"id"`
	ShareToken  string    `json:"share_token"`
	AdminToken  string    `json:"admin_token,omitempty"`
	Description string    `json:"description"`
	UpdatedAt   string    `json:"updated_at"`
	CreatedAt   string    `json:"created_at"`
}

func ProjectResponseFromProject(project repository.Project, withAdminToken bool) *ProjectResponse {
	adminToken := ""
	if withAdminToken {
		adminToken = project.AdminToken
	}
	return new(ProjectResponse{
		ID:          project.ID,
		ShareToken:  project.ShareToken,
		AdminToken:  adminToken,
		Description: project.Description,
		UpdatedAt:   project.UpdatedAt.Format(time.RFC3339),
		CreatedAt:   project.CreatedAt.Format(time.RFC3339),
	})
}
