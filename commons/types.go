package commons

import "k8s.io/helm/pkg/repo"

type InstallReleaseRequest struct {
	ReleaseName  string `json:"releaseName"`
	ChartName    string `json:"chartName"`
	Repo         string `json:"repo"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	ChartVersion string `json:"chartVersion"`
	Namespace    string `json:"namespace"`
	Description  string `json:"description"`
	Verify       bool   `json:"verify"`
}

type GetReleaseRequest struct {
	ReleaseName string
}

type DeleteReleaseRequest struct {
	ReleaseName string
}

type AddRepoRequest struct {
	Name     string `json:"name"`
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
	CertFile string `json:"certFile"`
	KeyFile  string `json:"keyFile"`
	CaFile   string `json:"caFile"`
	Noupdate bool   `json:"noupdate"`
}

type ListResult struct {
	Next     string
	Releases []ListRelease
}

type ListRelease struct {
	Name       string
	Revision   int32
	Updated    string
	Status     string
	Chart      string
	AppVersion string
	Namespace  string
}

type ListReposResponse struct {
	Repo []*repo.Entry
}

type ReleaseResource struct {
	ChartIcon    string `json:"chartIcon"`
	ChartName    string `json:"chartName"`
	ChartVersion string `json:"chartVersion"`
	Name         string `json:"name"`
	Namespace    string `json:"namespace"`
	Status       string `json:"status"`
	Updated      string `json:"updated"`
}
