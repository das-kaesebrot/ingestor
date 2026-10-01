package api

import (
	"dev.kaesebrot.eu/go/ingestor/internal/repository"
	"dev.kaesebrot.eu/go/ingestor/internal/utility"
)

func ProjectResponseFromProject(project repository.Project, withAdminToken bool) *ProjectResponse {
	var adminToken = new(string)
	if withAdminToken {
		*adminToken = (project.AdminToken)
	}
	return new(ProjectResponse{
		Id:          (ProjectId)(utility.ToGoogleUUID(project.ID)),
		ShareToken:  project.ShareToken,
		AdminToken:  adminToken,
		Description: project.Description,
		UpdatedAt:   project.UpdatedAt,
		CreatedAt:   project.CreatedAt,
	})
}
