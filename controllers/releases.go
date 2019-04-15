package controllers

import (
	"encoding/json"
	"zig-helm/commons"
	"zig-helm/services"

	"github.com/astaxie/beego"
)

type ReleaseController struct {
	HelmClient services.Release
	beego.Controller
}

// @Title List
// @Description list all releases
// @Success 200 {object} commons.ListResult
// @Failure 403
// @router / [get]
func (r *ReleaseController) List() {
	listResult, err := r.HelmClient.ListReleases()

	if err == nil {
		if listResult == nil {
			r.CustomAbort(403, "Not found releases")
		} else {
			r.Data["json"] = listResult
			r.ServeJSON()
		}
	} else {
		r.CustomAbort(403, err.Error())
	}
}

// @Title Install
// @Description install release
// @Param	body	body 	commons.InstallReleaseRequest	true 	"body content"
// @Success 200 {object} 	rls.InstallReleaseResponse
// @Failure 403
// @router / [post]
func (r *ReleaseController) Install() {
	installReleaseRequest := new(commons.InstallReleaseRequest)
	err := json.Unmarshal(r.Ctx.Input.RequestBody, installReleaseRequest)
	if nil != err {
		r.CustomAbort(403, err.Error())
	}
	releaseResource, err := r.HelmClient.InstallRelease(installReleaseRequest)
	if err == nil {
		r.Data["json"] = releaseResource
		r.ServeJSON()
	} else {
		r.CustomAbort(403, err.Error())
	}
}
