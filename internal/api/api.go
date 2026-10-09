package api

import (
	"encoding/json/v2"
	"errors"
	"net/http"

	"dev.kaesebrot.eu/go/ingestor/internal/repository"
	"gorm.io/gorm"
)

type APIHandler struct {
	repo *repository.Repository
}

// CreateDownloadProjectFilteredMediaArchive implements [ServerInterface].
func (a *APIHandler) CreateDownloadProjectFilteredMediaArchive(w http.ResponseWriter, r *http.Request, shareToken ShareToken) {
	panic("unimplemented")
}

// CreateProject implements [ServerInterface].
func (a *APIHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// CreateProjectUser implements [ServerInterface].
func (a *APIHandler) CreateProjectUser(w http.ResponseWriter, r *http.Request, shareToken ShareToken) {
	panic("unimplemented")
}

// DeleteProject implements [ServerInterface].
func (a *APIHandler) DeleteProject(w http.ResponseWriter, r *http.Request, adminToken AdminToken) {
	panic("unimplemented")
}

// DeleteProjectMedia implements [ServerInterface].
func (a *APIHandler) DeleteProjectMedia(w http.ResponseWriter, r *http.Request, adminToken AdminToken, mediaId MediaId) {
	panic("unimplemented")
}

// DeleteProjectMediaBatched implements [ServerInterface].
func (a *APIHandler) DeleteProjectMediaBatched(w http.ResponseWriter, r *http.Request, adminToken AdminToken) {
	panic("unimplemented")
}

// DeleteProjectUser implements [ServerInterface].
func (a *APIHandler) DeleteProjectUser(w http.ResponseWriter, r *http.Request, adminToken AdminToken, userId UserId) {
	panic("unimplemented")
}

// DownloadProjectFilteredMediaArchive implements [ServerInterface].
func (a *APIHandler) DownloadProjectFilteredMediaArchive(w http.ResponseWriter, r *http.Request, shareToken ShareToken, downloadToken FilteredDownloadToken, params DownloadProjectFilteredMediaArchiveParams) {
	panic("unimplemented")
}

// DownloadProjectMediaArchive implements [ServerInterface].
func (a *APIHandler) DownloadProjectMediaArchive(w http.ResponseWriter, r *http.Request, shareToken ShareToken, params DownloadProjectMediaArchiveParams) {
	panic("unimplemented")
}

// DownloadProjectMediaFile implements [ServerInterface].
func (a *APIHandler) DownloadProjectMediaFile(w http.ResponseWriter, r *http.Request, shareToken ShareToken, mediaId MediaId, params DownloadProjectMediaFileParams) {
	panic("unimplemented")
}

// DownloadProjectMediaPreview implements [ServerInterface].
func (a *APIHandler) DownloadProjectMediaPreview(w http.ResponseWriter, r *http.Request, shareToken ShareToken, mediaId MediaId, params DownloadProjectMediaPreviewParams) {
	panic("unimplemented")
}

// GetPing implements [ServerInterface].
func (a *APIHandler) GetPing(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, PingResponse{Message: PongMessagePong})
}

// GetProjectBackgroundJobs implements [ServerInterface].
func (a *APIHandler) GetProjectBackgroundJobs(w http.ResponseWriter, r *http.Request, shareToken ShareToken, params GetProjectBackgroundJobsParams) {
	panic("unimplemented")
}

// GetProjectByShareToken implements [ServerInterface].
func (a *APIHandler) GetProjectByShareToken(w http.ResponseWriter, r *http.Request, shareToken ShareToken) {
	_, err := a.repo.ProjectByShareToken(r.Context(), shareToken)

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			writeJSON(w, http.StatusNotFound, NotFound{Status: http.StatusNotFound, Title: "The requested project was not found.", AdditionalProperties: map[string]interface{}{"share_token": shareToken}})
		default:
			writeJSON(w, http.StatusInternalServerError, InternalServerError{Status: http.StatusInternalServerError, Title: err.Error()})
		}
		return
	}
}

// GetProjectMedia implements [ServerInterface].
func (a *APIHandler) GetProjectMedia(w http.ResponseWriter, r *http.Request, shareToken ShareToken, mediaId MediaId) {
	panic("unimplemented")
}

// GetProjectUser implements [ServerInterface].
func (a *APIHandler) GetProjectUser(w http.ResponseWriter, r *http.Request, shareToken ShareToken, userId UserId) {
	panic("unimplemented")
}

// GetSystemCapabilities implements [ServerInterface].
func (a *APIHandler) GetSystemCapabilities(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// ListProjectMediaMetadata implements [ServerInterface].
func (a *APIHandler) ListProjectMediaMetadata(w http.ResponseWriter, r *http.Request, shareToken ShareToken, params ListProjectMediaMetadataParams) {
	panic("unimplemented")
}

// ListProjectUsers implements [ServerInterface].
func (a *APIHandler) ListProjectUsers(w http.ResponseWriter, r *http.Request, shareToken ShareToken, params ListProjectUsersParams) {
	panic("unimplemented")
}

// ListProjects implements [ServerInterface].
func (a *APIHandler) ListProjects(w http.ResponseWriter, r *http.Request, params ListProjectsParams) {
	panic("unimplemented")
}

// PatchProject implements [ServerInterface].
func (a *APIHandler) PatchProject(w http.ResponseWriter, r *http.Request, adminToken AdminToken) {
	panic("unimplemented")
}

// PatchProjectUser implements [ServerInterface].
func (a *APIHandler) PatchProjectUser(w http.ResponseWriter, r *http.Request, shareToken ShareToken, userId UserId) {
	panic("unimplemented")
}

// SyncProjectMediaFile implements [ServerInterface].
func (a *APIHandler) SyncProjectMediaFile(w http.ResponseWriter, r *http.Request, shareToken ShareToken, mediaId MediaId) {
	panic("unimplemented")
}

// UpdateProjectMedia implements [ServerInterface].
func (a *APIHandler) UpdateProjectMedia(w http.ResponseWriter, r *http.Request, shareToken ShareToken, mediaId MediaId) {
	panic("unimplemented")
}

// UploadProjectMedia implements [ServerInterface].
func (a *APIHandler) UploadProjectMedia(w http.ResponseWriter, r *http.Request, shareToken ShareToken) {
	panic("unimplemented")
}

var _ ServerInterface = (*APIHandler)(nil)

func NewAPIHandler(db *repository.Repository) *APIHandler {
	return &APIHandler{repo: db}
}

func writeJSON(w http.ResponseWriter, status int, object any) {
	jData, err := json.Marshal(object)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(jData)
	if err != nil {
		panic(err)
	}
}
