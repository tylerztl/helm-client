package commons

type InstallReleaseRequest struct {
	ChartID      string
	ChartVersion string
	Namespace    string
	ReleaseName  string
	DryRun       bool
}

type ListReleasesRequest struct {
}

type GetReleaseRequest struct {
	ReleaseName string
}

type DeleteReleaseRequest struct {
	ReleaseName string
}
