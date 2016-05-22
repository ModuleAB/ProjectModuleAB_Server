package controllers

import (
	"moduleab_server/version"

	"github.com/astaxie/beego"
)

type VersionController struct {
	beego.Controller
}

// @router / [get]
func (v *VersionController) Get() {
	v.Data["json"] = map[string]string{
		"code":    version.Version.GetCode(),
		"version": version.Version.GetVersion(),
	}
	v.ServeJSON()
}
