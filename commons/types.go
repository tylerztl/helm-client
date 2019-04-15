package commons

import "k8s.io/helm/pkg/repo"

type InstallReleaseRequest struct {
	ChartID      string
	ChartVersion string
	Namespace    string
	ReleaseName  string
	DryRun       bool
}

type GetReleaseRequest struct {
	ReleaseName string
}

type DeleteReleaseRequest struct {
	ReleaseName string
}

type AddRepoRequest struct {
	Name     string `json:"name" valid:"alpha,required"`
	Url      string `json:"url" valid:"url,required"`
	Username string `json:"username"`
	Password string `json:"password"`
	CertFile string `json:"certFile"`
	KeyFile  string `json:"keyFile"`
	CaFile   string `json:"caFile"`
	Noupdate bool   `json:"noupdate"`
}

type RemoveRepoRequest struct {
	Name string `json:"name" valid:"alpha,required"`
}

type ListReposResponse struct {
	Repo []*repo.Entry
}
