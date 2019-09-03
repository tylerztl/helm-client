// @APIVersion 1.0.0
// @Title Zig Helm RESTful API
// @Description Manage k8s Resources
// @Contact tailinzhang1993@gmail.com
// @TermsOfServiceUrl http://zhigui.com
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"helm-client/controllers"
	"helm-client/services/handlers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/releases",
			beego.NSInclude(
				&controllers.ReleaseController{HelmClient: handlers.NewReleaseHandler()},
			),
		),
		beego.NSNamespace("/repos",
			beego.NSInclude(
				&controllers.RepoController{HelmClient: handlers.NewRepoHandler()},
			),
		),
	)
	beego.AddNamespace(ns)
}
