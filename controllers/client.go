package controllers

import "github.com/astaxie/beego"

type ClientController struct {
	beego.Controller
}

// @Title getClientConf
// @Description getClientConf
// @Param	body		body 	models.Hosts	true		"body for host content"
// @Success 200
// @Failure 403 body is empty
// @router / [get]
func (c *ClientController) Get() {
	c.Data["json"] = map[string]string{
		"ali_key":    beego.AppConfig.String("aliapi::apikey"),
		"ali_secret": beego.AppConfig.String("aliapi::secret"),
	}
	c.ServeJSON()
}
