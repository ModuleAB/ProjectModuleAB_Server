package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["moduleab_server/controllers:AppSetsController"] = append(beego.GlobalControllerRouter["moduleab_server/controllers:AppSetsController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["moduleab_server/controllers:AppSetsController"] = append(beego.GlobalControllerRouter["moduleab_server/controllers:AppSetsController"],
		beego.ControllerComments{
			"Get",
			`/:name`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["moduleab_server/controllers:AppSetsController"] = append(beego.GlobalControllerRouter["moduleab_server/controllers:AppSetsController"],
		beego.ControllerComments{
			"GetAll",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["moduleab_server/controllers:AppSetsController"] = append(beego.GlobalControllerRouter["moduleab_server/controllers:AppSetsController"],
		beego.ControllerComments{
			"Delete",
			`/:name`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["moduleab_server/controllers:AppSetsController"] = append(beego.GlobalControllerRouter["moduleab_server/controllers:AppSetsController"],
		beego.ControllerComments{
			"Put",
			`/:name`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["moduleab_server/controllers:ClientController"] = append(beego.GlobalControllerRouter["moduleab_server/controllers:ClientController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["moduleab_server/controllers:HostsController"] = append(beego.GlobalControllerRouter["moduleab_server/controllers:HostsController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["moduleab_server/controllers:HostsController"] = append(beego.GlobalControllerRouter["moduleab_server/controllers:HostsController"],
		beego.ControllerComments{
			"Get",
			`/:name`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["moduleab_server/controllers:HostsController"] = append(beego.GlobalControllerRouter["moduleab_server/controllers:HostsController"],
		beego.ControllerComments{
			"GetAll",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["moduleab_server/controllers:HostsController"] = append(beego.GlobalControllerRouter["moduleab_server/controllers:HostsController"],
		beego.ControllerComments{
			"Delete",
			`/:name`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["moduleab_server/controllers:HostsController"] = append(beego.GlobalControllerRouter["moduleab_server/controllers:HostsController"],
		beego.ControllerComments{
			"Put",
			`/:name`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["moduleab_server/controllers:OssController"] = append(beego.GlobalControllerRouter["moduleab_server/controllers:OssController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["moduleab_server/controllers:OssController"] = append(beego.GlobalControllerRouter["moduleab_server/controllers:OssController"],
		beego.ControllerComments{
			"Get",
			`/:name`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["moduleab_server/controllers:OssController"] = append(beego.GlobalControllerRouter["moduleab_server/controllers:OssController"],
		beego.ControllerComments{
			"GetAll",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["moduleab_server/controllers:OssController"] = append(beego.GlobalControllerRouter["moduleab_server/controllers:OssController"],
		beego.ControllerComments{
			"Delete",
			`/:name`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["moduleab_server/controllers:OssController"] = append(beego.GlobalControllerRouter["moduleab_server/controllers:OssController"],
		beego.ControllerComments{
			"Put",
			`/:name`,
			[]string{"put"},
			nil})

}
