package handlers

import (
	"fmt"
	"k8s.io/helm/pkg/helm"
	rls "k8s.io/helm/pkg/proto/hapi/services"
	"strings"
	"zig-helm/commons"
)

// ListReleases returns the list of helm releases
//func (h *HelmHandler) ListReleases(request *commons.ListReleasesRequest) (*rls.ListReleasesResponse, error) {
//	stats := []release.Status_Code{
//		release.Status_DEPLOYED,
//	}
//	resp, err := h.HelmClient.ListReleases(
//		helm.ReleaseListFilter(""),
//		helm.ReleaseListSort(int32(services.ListSort_LAST_RELEASED)),
//		helm.ReleaseListOrder(int32(services.ListSort_DESC)),
//		helm.ReleaseListStatuses(stats),
//	)
//
//	if err != nil {
//		log.WithError(err).Error("Can't retrieve the list of releases")
//		return nil, err
//	}
//
//	return resp, err
//}
//
//// GetRelease gets the information of an existing release
//func (h *HelmHandler) GetRelease(request *commons.GetReleaseRequest) (*rls.GetReleaseContentResponse, error) {
//	// TODO, find a way to retrieve all the information in a single call
//	// We get the information about the release
//	release, err := client.ReleaseContent(releaseName)
//	if err != nil {
//		return nil, err
//	}
//
//	// Now we populate the resources string
//	status, err := client.ReleaseStatus(releaseName)
//	if err != nil {
//		return nil, err
//	}
//	release.Release.Info = status.Info
//	return release, err
//}

// InstallRelease wraps helms client installReleae method
func (h *HelmHandler) InstallRelease(request *commons.InstallReleaseRequest) (*rls.InstallReleaseResponse, error) {

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
//func (h *HelmHandler) DeleteRelease(request *commons.DeleteReleaseRequest) (*rls.UninstallReleaseResponse, error) {
//	opts := []helm.DeleteOption{
//		helm.DeleteDryRun(false),
//		helm.DeletePurge(false),
//		helm.DeleteTimeout(300),
//	}
//	return client.DeleteRelease(releaseName, opts...)
//}
