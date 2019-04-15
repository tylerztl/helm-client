package handlers

import (
	"fmt"
	"k8s.io/helm/pkg/helm"
	rls "k8s.io/helm/pkg/proto/hapi/services"
	"strings"
	"zig-helm/commons"
	helm_client "zig-helm/services/helm"
)

type ReleaseHandler struct {
	HelmClient helm.Interface
}

func NewReleaseHandler() *ReleaseHandler {
	return &ReleaseHandler{
		HelmClient: helm_client.GetClient(),
	}
}

// ListReleases returns the list of helm releases
func (h *ReleaseHandler) ListReleases() (*rls.ListReleasesResponse, error) {
	return h.HelmClient.ListReleases(
		helm.ReleaseListFilter(""),
		helm.ReleaseListSort(int32(rls.ListSort_LAST_RELEASED)),
		helm.ReleaseListOrder(int32(rls.ListSort_DESC)),
	)
}

// GetRelease gets the information of an existing release
func (h *ReleaseHandler) GetRelease(request *commons.GetReleaseRequest) (*rls.GetReleaseContentResponse, error) {
	return h.HelmClient.ReleaseContent(request.ReleaseName)
}

// InstallRelease wraps helms client installReleae method
func (h *ReleaseHandler) InstallRelease(request *commons.InstallReleaseRequest) (*rls.InstallReleaseResponse, error) {

	idSplit := strings.Split(request.ChartID, "/")
	if len(idSplit) != 2 || idSplit[0] == "" || idSplit[1] == "" {
		return nil, fmt.Errorf("chartId must include the repository name. i.e: stable/wordpress")
	}

	// Search chart package and get local path
	repo, chartName := idSplit[0], idSplit[1]

	chartPath, err := commons.LocateChartPath(repo, "", "", chartName, request.ChartVersion, true, "", "", "", "")
	if err != nil {
		return nil, err
	}

	ns := request.Namespace
	if ns == "" {
		ns = "default"
	}

	return h.HelmClient.InstallRelease(
		chartPath,
		ns,
		helm.ValueOverrides([]byte{}),
		helm.ReleaseName(request.ReleaseName),
		helm.InstallDryRun(request.DryRun))
}

// DeleteRelease deletes an existing helm chart
func (h *ReleaseHandler) DeleteRelease(request *commons.DeleteReleaseRequest) (*rls.UninstallReleaseResponse, error) {
	opts := []helm.DeleteOption{
		helm.DeleteDryRun(false),
		helm.DeletePurge(false),
		helm.DeleteTimeout(300),
	}
	return h.HelmClient.DeleteRelease(request.ReleaseName, opts...)
}
