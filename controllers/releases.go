package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"zig-helm/commons"
	"zig-helm/services"
)

type ReleaseController struct {
	HelmClient services.Client
	beego.Controller
}

// @Title Install
// @Description install release
// @Param	body	body 	commons.InstallReleaseRequest	true 	"body content"
// @Success 200 {object} commons.InstallReleaseResponse
// @Failure 403
// @router /vpc [post]
func (r *ReleaseController) Install() {
	installReleaseRequest := new(commons.InstallReleaseRequest)
	err := json.Unmarshal(r.Ctx.Input.RequestBody, installReleaseRequest)
	if nil != err {
		r.CustomAbort(403, err.Error())
	}
	installReleaseResponse, err := r.HelmClient.InstallRelease(installReleaseRequest)
	if err == nil {
		r.Data["json"] = installReleaseResponse
		r.ServeJSON()
	} else {
		r.CustomAbort(403, err.Error())
	}
}
