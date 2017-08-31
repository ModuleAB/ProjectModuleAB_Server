package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:AppSetsController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:AppSetsController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:AppSetsController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:AppSetsController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:name`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:AppSetsController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:AppSetsController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:AppSetsController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:AppSetsController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:name`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:AppSetsController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:AppSetsController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:name`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:BackupSetsController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:BackupSetsController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:BackupSetsController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:BackupSetsController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:name`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:BackupSetsController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:BackupSetsController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:BackupSetsController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:BackupSetsController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:name`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:BackupSetsController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:BackupSetsController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:name`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:ClientController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:ClientController"],
		beego.ControllerComments{
			Method: "GetAliConfig",
			Router: `/config`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:ClientController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:ClientController"],
		beego.ControllerComments{
			Method: "WebSocket",
			Router: `/signal/:name/ws`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:ClientController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:ClientController"],
		beego.ControllerComments{
			Method: "GetSignals",
			Router: `/signal/:name`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:ClientController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:ClientController"],
		beego.ControllerComments{
			Method: "GetSignal",
			Router: `/signal/:name/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:ClientController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:ClientController"],
		beego.ControllerComments{
			Method: "PostSignal",
			Router: `/signal/:name`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:ClientController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:ClientController"],
		beego.ControllerComments{
			Method: "DeleteSignal",
			Router: `/signal/:name/:id`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:ClientController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:ClientController"],
		beego.ControllerComments{
			Method: "NotifySignal",
			Router: `/signal/:name/:id`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:ClientController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:ClientController"],
		beego.ControllerComments{
			Method: "GetStatus",
			Router: `/config/status`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:ClientJobsController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:ClientJobsController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:ClientJobsController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:ClientJobsController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:ClientJobsController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:ClientJobsController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:ClientJobsController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:ClientJobsController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:ClientJobsController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:ClientJobsController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:HostsController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:HostsController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:HostsController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:HostsController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:name`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:HostsController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:HostsController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:HostsController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:HostsController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:name`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:HostsController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:HostsController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:name`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:LoginController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:LoginController"],
		beego.ControllerComments{
			Method: "Login",
			Router: `/login`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:LoginController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:LoginController"],
		beego.ControllerComments{
			Method: "Logout",
			Router: `/logout`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:LoginController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:LoginController"],
		beego.ControllerComments{
			Method: "Check",
			Router: `/check`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:OasController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:OasController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:OasController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:OasController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:name`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:OasController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:OasController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:OasController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:OasController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:name`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:OasController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:OasController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:name`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:OasJobsController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:OasJobsController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:job_id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:OasJobsController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:OasJobsController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:OssController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:OssController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:OssController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:OssController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:name`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:OssController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:OssController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:OssController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:OssController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:name`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:OssController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:OssController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:name`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:PathsController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:PathsController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:PathsController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:PathsController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:PathsController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:PathsController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:PathsController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:PathsController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:PathsController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:PathsController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:PolicyController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:PolicyController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:PolicyController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:PolicyController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:name`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:PolicyController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:PolicyController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:PolicyController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:PolicyController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:name`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:PolicyController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:PolicyController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:name`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:RecordsController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:RecordsController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:RecordsController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:RecordsController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:RecordsController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:RecordsController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:RecordsController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:RecordsController"],
		beego.ControllerComments{
			Method: "Recover",
			Router: `/:id/recover`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:RolesController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:RolesController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:UserController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:UserController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:name`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:UserController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:UserController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:name`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:UserController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:name`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:VersionController"] = append(beego.GlobalControllerRouter["github.com/ModuleAB/ModuleAB/server/controllers:VersionController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

}
