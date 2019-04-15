package services

import (
	"zig-helm/commons"

	rls "k8s.io/helm/pkg/proto/hapi/services"
)

// Client is an interface for managing Helm Chart releases
type Release interface {
	ListReleases() (*rls.ListReleasesResponse, error)
	InstallRelease(*commons.InstallReleaseRequest) (*rls.InstallReleaseResponse, error)
	GetRelease(*commons.GetReleaseRequest) (*rls.GetReleaseContentResponse, error)
	DeleteRelease(*commons.DeleteReleaseRequest) (*rls.UninstallReleaseResponse, error)
}

type Repo interface {
	AddRepo(*commons.AddRepoRequest) error
	RemoveRepo(*commons.RemoveRepoRequest) error
	ListRepos() (*commons.ListReposResponse, error)
}
