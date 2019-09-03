package controllers

import (
	"encoding/json"
	"fmt"
	"helm-client/commons"
	"helm-client/services"

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
// @Success 200 {object} 	commons.ReleaseResource
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

// @Title Get
// @Description get a named release content
// @@Param	repo   path 	string	true	"The release you want to get"
// @Success 200 {object}  commons.ReleaseResource
// @Failure 403
// @router /:name [get]
func (r *ReleaseController) Get() {
	releaseName := r.GetString(":name")
	if releaseName == "" {
		r.CustomAbort(403, "Release name must be provide")
	}
	releaseExtended, err := r.HelmClient.GetRelease(releaseName)
	if err == nil {
		r.Data["json"] = releaseExtended
		r.ServeJSON()
	} else {
		r.CustomAbort(403, err.Error())
	}
}

// @Title Delete
// @Description delete release
// @@Param	repo   path 	string	true	"The release you want to delete"
// @Success 200 {string}  Release has been Deleted
// @Failure 403
// @router /:name [delete]
func (r *ReleaseController) Delete() {
	releaseName := r.GetString(":name")
	if releaseName == "" {
		r.CustomAbort(403, "Release name must be provide")
	}
	_, err := r.HelmClient.DeleteRelease(releaseName)
	if err == nil {
		r.Data["json"] = fmt.Sprintf("%q has been deleted", releaseName)
		r.ServeJSON()
	} else {
		r.CustomAbort(403, err.Error())
	}
}
