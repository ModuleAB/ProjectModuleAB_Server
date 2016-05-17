// @APIVersion 1.0.0
// @Title ModuleAB API
// @Description ModuleAB server API
// @Contact tonychyi1989@gmail.com
// @License GPLv3
// @LicenseUrl http://www.gnu.org/licenses/gpl-3.0.html
package routers

import (
	"moduleab_server/controllers"
	"net/http"
	"path"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

// Show web
func StaticFileServer(ctx *context.Context) {
	if strings.HasPrefix(ctx.Input.URL(), "/api") {
		return
	}
	filename := path.Join("web/app", ctx.Input.URL())
	http.ServeFile(ctx.ResponseWriter, ctx.Request, filename)
}

func init() {
	beego.ErrorController(&controllers.ErrorController{})
	ns := beego.NewNamespace("/api/v1",
		beego.NSNamespace("/hosts",
			beego.NSInclude(
				&controllers.HostsController{},
			),
		),
		beego.NSNamespace("/client",
			beego.NSInclude(
				&controllers.ClientController{},
			),
		),
		beego.NSNamespace("/appSets",
			beego.NSInclude(
				&controllers.AppSetsController{},
			),
		),
		beego.NSNamespace("/backupSets",
			beego.NSInclude(
				&controllers.BackupSetsController{},
			),
		),
		beego.NSNamespace("/clientJobs",
			beego.NSInclude(
				&controllers.ClientJobsController{},
			),
		),
		beego.NSNamespace("/oss",
			beego.NSInclude(
				&controllers.OssController{},
			),
		),
		beego.NSNamespace("/oas",
			beego.NSInclude(
				&controllers.OasController{},
			),
		),
		beego.NSNamespace("/oasJobs",
			beego.NSInclude(
				&controllers.OasJobsController{},
			),
		),
		beego.NSNamespace("/users",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/paths",
			beego.NSInclude(
				&controllers.PathsController{},
			),
		),
		beego.NSNamespace("/policies",
			beego.NSInclude(
				&controllers.PolicyController{},
			),
		),
		beego.NSNamespace("/records",
			beego.NSInclude(
				&controllers.RecordsController{},
			),
		),
		beego.NSNamespace("/auth",
			beego.NSInclude(
				&controllers.LoginController{},
			),
		),
		beego.NSNamespace("/roles",
			beego.NSInclude(
				&controllers.RolesController{},
			),
		),
	)
	beego.AddNamespace(ns)
	beego.InsertFilter("/", beego.BeforeRouter, StaticFileServer)
	beego.InsertFilter("/*", beego.BeforeRouter, StaticFileServer)
}
