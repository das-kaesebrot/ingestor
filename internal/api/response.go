package api

import "dev.kaesebrot.eu/go/ingestor/internal/repository"

func PagedProjectsToResponse(ps []repository.Project, pageNum int, perPage int, total int64, withAdminToken bool) PagedProjects {
	pagedProjects := PagedProjects{}
	pagedProjects.Page = Page{PageNum: pageNum, PerPage: perPage, Total: int(total)}
	pagedProjects.Projects = make([]Project, len(ps))

	for i, project := range ps {
		pagedProjects.Projects[i] = ProjectToResponse(&project, withAdminToken)
	}

	return pagedProjects
}

func ProjectToResponse(p *repository.Project, withAdminToken bool) Project {
	shareTokenPermissions := make([]ProjectShareTokenPermission, 0, 2)
	if p.ShareTokenDownloadAllowed {
		shareTokenPermissions = append(shareTokenPermissions, Download)
	}
	if p.ShareTokenUploadAllowed {
		shareTokenPermissions = append(shareTokenPermissions, Upload)
	}

	projectResponse := Project{
		CreatedAt:             p.CreatedAt,
		Description:           p.Description,
		ExpiresAt:             p.ExpiresAt,
		Id:                    p.ID,
		ShareToken:            p.ShareToken,
		ShareTokenPermissions: shareTokenPermissions,
		UpdatedAt:             p.UpdatedAt,
	}

	if p.Start != nil {
		projectResponse.Start = new(*p.Start)
	}

	if p.End != nil {
		projectResponse.End = new(*p.End)
	}

	if withAdminToken {
		projectResponse.AdminToken = &p.AdminToken
	}

	return projectResponse
}

func PagedUsersToResponse(u []repository.UploadUser, pageNum int, perPage int, total int64) PagedUsers {
	pagedUsers := PagedUsers{}
	pagedUsers.Page = Page{PageNum: pageNum, PerPage: perPage, Total: int(total)}
	pagedUsers.Users = make([]User, len(u))

	for i, user := range u {
		pagedUsers.Users[i] = UserToResponse(&user)
	}

	return pagedUsers
}

func UserToResponse(u *repository.UploadUser) User {
	return User{
		CreatedAt: u.CreatedAt,
		Id:        u.ID,
		Name:      u.Name,
		UpdatedAt: u.UpdatedAt,
	}
}
