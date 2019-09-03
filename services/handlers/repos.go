package handlers

import (
	"helm-client/commons"

	"k8s.io/helm/pkg/repo"
)

// Repo describes a chart repository
type Repo struct {
	Name     string `json:"name" valid:"alpha,required"`
	Url      string `json:"url" valid:"url,required"`
	Username string `json:"username"`
	Password string `json:"password"`
	CertFile string `json:"certFile"`
	KeyFile  string `json:"keyFile"`
	CaFile   string `json:"caFile"`
	Noupdate bool   `json:"noupdate"`
}

// RepoHandlers defines handlers that serve chart data
type RepoHandler struct {
}

// NewRepoHandlers takes a datastore.Session implementation and returns a RepoHandlers struct
func NewRepoHandler() *RepoHandler {
	return &RepoHandler{}
}

// ListRepos returns all repositories
func (r *RepoHandler) ListRepos() (*commons.ListReposResponse, error) {

	f, err := repo.LoadRepositoriesFile(commons.GetConfig().Home.RepositoryFile())
	if err != nil {
		return nil, err
	}
	if len(f.Repositories) == 0 {
		return nil, err
	}

	return &commons.ListReposResponse{Repo: f.Repositories}, nil
}

// AddRepo adds a repo to the list of enabled repositories to index
func (r *RepoHandler) AddRepo(request *commons.AddRepoRequest) error {
	return commons.AddRepository(request.Name, request.Url, request.Username, request.Password,
		commons.GetConfig().Home, request.CertFile, request.KeyFile, request.CaFile, request.Noupdate)
}

// DeleteRepo deletes a repo from the list of enabled repositories to index
func (r *RepoHandler) RemoveRepo(repo string) error {
	return commons.RemoveRepoLine(repo, commons.GetConfig().Home)
}
