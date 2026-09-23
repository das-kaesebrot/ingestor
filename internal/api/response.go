package api

import (
	"uuid"

	"dev.kaesebrot.eu/go/ingestor/internal/repository"
)

type ProblemDetails struct {
}

type ProjectResponse struct {
	ID          uuid.UUID `json:"id"`
	ShareToken  string    `json:"share_token"`
	AdminToken  string    `json:"admin_token,omitempty"`
	Description string    `json:"description"`
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
	})
}
