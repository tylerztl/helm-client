package handlers

import (
	"k8s.io/helm/pkg/helm"
	helm_client "zig-helm/services/helm"
)

type HelmHandler struct {
	HelmClient helm.Interface
}

func NewHelmHandler() *HelmHandler {
	return &HelmHandler{
		HelmClient: helm_client.GetClient(),
	}
}
