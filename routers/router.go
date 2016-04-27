// @APIVersion 1.0.0
// @Title ModulesAB API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"moduleab_server/controllers"

	"github.com/astaxie/beego"
)

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
		beego.NSNamespace("/oss",
			beego.NSInclude(
				&controllers.OssController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
