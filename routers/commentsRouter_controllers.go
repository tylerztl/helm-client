package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["helm-client/controllers:ReleaseController"] = append(beego.GlobalControllerRouter["helm-client/controllers:ReleaseController"],
        beego.ControllerComments{
            Method: "List",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["helm-client/controllers:ReleaseController"] = append(beego.GlobalControllerRouter["helm-client/controllers:ReleaseController"],
        beego.ControllerComments{
            Method: "Install",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["helm-client/controllers:ReleaseController"] = append(beego.GlobalControllerRouter["helm-client/controllers:ReleaseController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/:name`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["helm-client/controllers:ReleaseController"] = append(beego.GlobalControllerRouter["helm-client/controllers:ReleaseController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:name`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["helm-client/controllers:RepoController"] = append(beego.GlobalControllerRouter["helm-client/controllers:RepoController"],
        beego.ControllerComments{
            Method: "Add",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["helm-client/controllers:RepoController"] = append(beego.GlobalControllerRouter["helm-client/controllers:RepoController"],
        beego.ControllerComments{
            Method: "List",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["helm-client/controllers:RepoController"] = append(beego.GlobalControllerRouter["helm-client/controllers:RepoController"],
        beego.ControllerComments{
            Method: "Remove",
            Router: `/:repo`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["helm-client/controllers:UserController"] = append(beego.GlobalControllerRouter["helm-client/controllers:UserController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["helm-client/controllers:UserController"] = append(beego.GlobalControllerRouter["helm-client/controllers:UserController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["helm-client/controllers:UserController"] = append(beego.GlobalControllerRouter["helm-client/controllers:UserController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/:uid`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["helm-client/controllers:UserController"] = append(beego.GlobalControllerRouter["helm-client/controllers:UserController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:uid`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["helm-client/controllers:UserController"] = append(beego.GlobalControllerRouter["helm-client/controllers:UserController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:uid`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["helm-client/controllers:UserController"] = append(beego.GlobalControllerRouter["helm-client/controllers:UserController"],
        beego.ControllerComments{
            Method: "Login",
            Router: `/login`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["helm-client/controllers:UserController"] = append(beego.GlobalControllerRouter["helm-client/controllers:UserController"],
        beego.ControllerComments{
            Method: "Logout",
            Router: `/logout`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
