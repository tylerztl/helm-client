package handlers

import (
	"helm-client/commons"
	helm_client "helm-client/services/helm"

	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/release"
	rls "k8s.io/helm/pkg/proto/hapi/services"
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
func (h *ReleaseHandler) ListReleases() (*commons.ListResult, error) {
	res, err := h.HelmClient.ListReleases(
		helm.ReleaseListLimit(256),
		helm.ReleaseListOffset(""),
		helm.ReleaseListFilter(""),
		helm.ReleaseListSort(int32(rls.ListSort_LAST_RELEASED)),
		helm.ReleaseListOrder(int32(rls.ListSort_DESC)),
		helm.ReleaseListStatuses([]release.Status_Code{
			release.Status_UNKNOWN,
			release.Status_DEPLOYED,
			release.Status_DELETED,
			release.Status_DELETING,
			release.Status_FAILED,
			release.Status_PENDING_INSTALL,
			release.Status_PENDING_UPGRADE,
			release.Status_PENDING_ROLLBACK,
		}),
		helm.ReleaseListNamespace(""),
	)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}

	rels := commons.FilterList(res.GetReleases())

	return commons.GetListResult(rels, res.Next), nil
}

// GetRelease gets the information of an existing release
func (h *ReleaseHandler) GetRelease(releaseName string) (*commons.ReleaseExtended, error) {
	res, err := h.HelmClient.ReleaseContent(releaseName)
	if err != nil {
		return nil, err
	}
	return commons.MakeReleaseExtendedResource(res.Release), nil
}

// InstallRelease wraps helms client installReleae method
func (h *ReleaseHandler) InstallRelease(request *commons.InstallReleaseRequest) (*commons.ReleaseResource, error) {
	chartPath, err := commons.LocateChartPath(request.Repo, request.Username, request.Password, request.ChartName,
		request.ChartVersion, request.Verify, "", "", "", "")
	if err != nil {
		return nil, err
	}

	ns := request.Namespace
	if ns == "" {
		ns = commons.DefaultNamespace()
	}

	res, err := h.HelmClient.InstallRelease(
		chartPath,
		ns,
		helm.ValueOverrides([]byte{}),
		helm.ReleaseName(request.ReleaseName),
		helm.InstallTimeout(300),
		helm.InstallWait(false),
		helm.InstallDescription(request.Description))
	if err != nil {
		return nil, err
	}
	return commons.MakeReleaseResource(res.Release), nil
}

// DeleteRelease deletes an existing helm chart
func (h *ReleaseHandler) DeleteRelease(releaseName string) (*rls.UninstallReleaseResponse, error) {
	opts := []helm.DeleteOption{
		helm.DeleteDryRun(false),
		helm.DeleteDisableHooks(false),
		helm.DeletePurge(false),
		helm.DeleteTimeout(300),
		helm.DeleteDescription(""),
	}
	return h.HelmClient.DeleteRelease(releaseName, opts...)
}
