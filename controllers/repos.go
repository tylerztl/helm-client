package controllers

import (
	"encoding/json"
	"fmt"
	"zig-helm/commons"
	"zig-helm/services"

	"github.com/astaxie/beego"
)

type RepoController struct {
	HelmClient services.Repo
	beego.Controller
}

// @Title Add
// @Description add repo
// @Param	body	body 	commons.AddRepoRequest	true 	"body content"
// @Success 200 {string} Repo has been added to your repositories
// @Failure 403
// @router / [post]
func (r *RepoController) Add() {
	addRepoRequest := new(commons.AddRepoRequest)
	err := json.Unmarshal(r.Ctx.Input.RequestBody, addRepoRequest)
	if nil != err {
		r.CustomAbort(403, err.Error())
	}
	err = r.HelmClient.AddRepo(addRepoRequest)
	if err == nil {
		r.Data["json"] = fmt.Sprintf("%q has been added to your repositories", addRepoRequest.Name)
		r.ServeJSON()
	} else {
		r.CustomAbort(403, err.Error())
	}
}

// @Title Remove
// @Description remove repo
// @Param	body	body 	commons.RemoveRepoRequest	true 	"body content"
// @Success 200 {string}  Repo has been removed from your repositories
// @Failure 403
// @router / [delete]
func (r *RepoController) Remove() {
	removeRepoRequest := new(commons.RemoveRepoRequest)
	err := json.Unmarshal(r.Ctx.Input.RequestBody, removeRepoRequest)
	if nil != err {
		r.CustomAbort(403, err.Error())
	}
	err = r.HelmClient.RemoveRepo(removeRepoRequest)
	if err == nil {
		r.Data["json"] = fmt.Sprintf("%q has been removed from your repositories", removeRepoRequest.Name)
		r.ServeJSON()
	} else {
		r.CustomAbort(403, err.Error())
	}
}

// @Title List
// @Description list all repo
// @Param	body	body 	commons.ListReposRequest	true 	"body content"
// @Success 200 {object} commons.InstallReleaseResponse
// @Failure 403
// @router / [get]
func (r *RepoController) List() {
	listReposResponse, err := r.HelmClient.ListRepos()

	repoList := make(map[string]string)
	for _, re := range listReposResponse.Repo {
		repoList[re.Name] = re.URL
	}
	if err == nil {
		r.Data["json"] = repoList
		r.ServeJSON()
	} else {
		r.CustomAbort(403, err.Error())
	}
}
