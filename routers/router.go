// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"zig-helm/controllers"
	"zig-helm/services/handlers"

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
