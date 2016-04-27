package controllers

import (
	"encoding/json"
	"fmt"
	"moduleab_server/models"
	"net/http"

	"github.com/astaxie/beego"
)

type ClientController struct {
	beego.Controller
}

// @Title getClientConf
// @Description getClientConf
// @Param	body		body 	models.Hosts	true		"body for host content"
// @Success 200
// @Failure 403 body is empty
// @router /config [get]
func (c *ClientController) GetAliConfig() {
	c.Data["json"] = map[string]string{
		"ali_key":    beego.AppConfig.String("aliapi::apikey"),
		"ali_secret": beego.AppConfig.String("aliapi::secret"),
	}
	c.ServeJSON()
}

// @router /signal/:name [get]
func (c *ClientController) GetSignals() {
	name := c.GetString(":name")
	beego.Debug("[C] Got name:", name)
	if name != "" {
		host := &models.Hosts{
			Name: name,
		}
		hosts, err := models.GetHosts(host, 0, 0)
		if err != nil {
			c.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.ServeJSON()
			return
		}
		if len(hosts) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			c.Ctx.Output.SetStatus(http.StatusNotFound)
			c.ServeJSON()
			return
		}
		c.Data["json"] = models.GetSignals(hosts[0].Id)
		c.Ctx.Output.SetStatus(http.StatusOK)
		c.ServeJSON()
		return

	}
}

// @router /signal/:name/:id [get]
func (c *ClientController) GetSignal() {
	name := c.GetString(":name")
	id := c.GetString(":id")
	beego.Debug("[C] Got name:", name)
	if name != "" {
		host := &models.Hosts{
			Name: name,
		}
		hosts, err := models.GetHosts(host, 0, 0)
		if err != nil {
			c.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.ServeJSON()
			return
		}
		if len(hosts) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			c.Ctx.Output.SetStatus(http.StatusNotFound)
			c.ServeJSON()
			return
		}
		signals := models.GetSignals(hosts[0].Id)
		for _, v := range signals {
			if v["id"] == id {
				c.Data["json"] = v
				c.Ctx.Output.SetStatus(http.StatusOK)
				c.ServeJSON()
				return
			}
		}

		c.Data["json"] = map[string]string{
			"message": fmt.Sprintf("Got nothing with id:", id),
		}
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.ServeJSON()
		return

	}
}

// @router /signal/:name [post]
func (c *ClientController) PostSignal() {
	name := c.GetString(":name")
	beego.Debug("[C] Got name:", name)
	if name != "" {
		host := &models.Hosts{
			Name: name,
		}
		hosts, err := models.GetHosts(host, 0, 0)
		if err != nil {
			c.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.ServeJSON()
			return
		}
		if len(hosts) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			c.Ctx.Output.SetStatus(http.StatusNotFound)
			c.ServeJSON()
			return
		}
		var signal models.Signal
		err = json.Unmarshal(c.Ctx.Input.RequestBody, &signal)
		if err != nil {
			beego.Warn("[C] Got error:", err)
			c.Data["json"] = map[string]string{
				"message": "Bad request",
				"error":   err.Error(),
			}
			c.Ctx.Output.SetStatus(http.StatusBadRequest)
			c.ServeJSON()
			return
		}
		beego.Debug("[C] Got data:", signal)
		id, err := models.AddSignal(hosts[0].Id, signal)
		if err != nil {
			beego.Warn("[C] Got error:", err)
			c.Data["json"] = map[string]string{
				"message": "Failed to add new signal",
				"error":   err.Error(),
			}
			c.Ctx.Output.SetStatus(http.StatusBadRequest)
			c.ServeJSON()
			return
		}
		c.Data["json"] = map[string]string{
			"id": id,
		}
		c.Ctx.Output.SetStatus(http.StatusCreated)
		c.ServeJSON()
		return
	}
}

// @router /signal/:name/:id [delete]
func (c *ClientController) DeleteSignal() {
	name := c.GetString(":name")
	id := c.GetString(":id")
	beego.Debug("[C] Got name:", name)
	if name != "" {
		host := &models.Hosts{
			Name: name,
		}
		hosts, err := models.GetHosts(host, 0, 0)
		if err != nil {
			c.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.ServeJSON()
			return
		}
		if len(hosts) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			c.Ctx.Output.SetStatus(http.StatusNotFound)
			c.ServeJSON()
			return
		}
		err = models.DeleteSignal(hosts[0].Id, id)
		if err != nil {
			beego.Warn("[C] Got error:", err)
			c.Data["json"] = map[string]string{
				"message": "Delete failed",
				"error":   err.Error(),
			}
			if err == models.ErrorSignalNotFound {
				c.Ctx.Output.SetStatus(http.StatusNotFound)
			} else {
				c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			}
			c.ServeJSON()
			return
		}

		c.Data["json"] = map[string]string{
			"message": fmt.Sprintf("Got nothing with id:", id),
		}
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.ServeJSON()
		return

	}
}
