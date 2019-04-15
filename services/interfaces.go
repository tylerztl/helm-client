package services

import (
	"zig-helm/commons"

	rls "k8s.io/helm/pkg/proto/hapi/services"
)

// Client is an interface for managing Helm Chart releases
type Release interface {
	ListReleases() (*commons.ListResult, error)
	InstallRelease(*commons.InstallReleaseRequest) (*commons.ReleaseResource, error)
	GetRelease(release string) (*commons.ReleaseExtended, error)
	DeleteRelease(release string) (*rls.UninstallReleaseResponse, error)
}

type Repo interface {
	AddRepo(*commons.AddRepoRequest) error
	RemoveRepo(repo string) error
	ListRepos() (*commons.ListReposResponse, error)
}
