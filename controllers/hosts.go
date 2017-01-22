package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ModuleAB/ModuleAB/server/common"
	"github.com/ModuleAB/ModuleAB/server/models"

	"github.com/astaxie/beego"
)

func init() {
	AddPrivilege("GET", "^/api/v1/hosts", models.RoleFlagUser)
}

type HostsController struct {
	beego.Controller
}

func (h *HostsController) Prepare() {
	if h.Ctx.Input.Header("Signature") != "" {
		err := common.AuthWithKey(h.Ctx)
		if err != nil {
			h.Data["json"] = map[string]string{
				"error": err.Error(),
			}
			h.Ctx.Output.SetStatus(http.StatusForbidden)
			h.ServeJSON()
		}
	} else {
		id := h.GetSession("id")
		if id == nil {
			h.Data["json"] = map[string]string{
				"error": "You need login first.",
			}
			h.Ctx.Output.SetStatus(http.StatusUnauthorized)
			h.ServeJSON()
		} else {
			if !CheckPrivileges(id.(string), h.Ctx) {
				h.Data["json"] = map[string]string{
					"error": "No privileges.",
				}
				h.Ctx.Output.SetStatus(http.StatusForbidden)
				h.ServeJSON()
			}
		}
	}
}

// @Title createHost
// @Description create Host
// @Param	host 	body 	object true	"host"
// @Success 201 {object} models.Hosts
// @Failure 400 Hostname or IP missing
// @Failure 500 Failure on writing database
// @router / [post]
func (h *HostsController) Post() {
	host := new(models.Hosts)
	defer h.ServeJSON()
	err := json.Unmarshal(h.Ctx.Input.RequestBody, host)
	if err != nil {
		beego.Warn("[C] Got error:", err)
		h.Data["json"] = map[string]string{
			"message": "Bad request",
			"error":   err.Error(),
		}
		h.Ctx.Output.SetStatus(http.StatusBadRequest)
		return
	}
	beego.Debug("[C] Got data:", host)
	id, err := models.AddHost(host)
	if err != nil {
		beego.Warn("[C] Got error:", err)
		h.Data["json"] = map[string]string{
			"message": "Failed to add New host",
			"error":   err.Error(),
		}
		h.Ctx.Output.SetStatus(http.StatusInternalServerError)
		return
	}

	beego.Debug("[C] Got id:", id)
	h.Data["json"] = map[string]string{
		"id": id,
	}
	h.Ctx.Output.SetStatus(http.StatusCreated)
}

// @Title getHost
// @Description get Host info
// @Param	body		body 	models.Hosts	true		"body for host content"
// @Success 200 {string} models.Hosts.Id
// @Failure 403 body is empty
// @router /:name [get]
func (h *HostsController) Get() {
	name := h.GetString(":name")
	defer h.ServeJSON()
	beego.Debug("[C] Got name:", name)
	if name != "" {
		host := &models.Hosts{
			Name: name,
		}
		hosts, err := models.GetHosts(host, 0, 0)
		if err != nil {
			h.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			h.Ctx.Output.SetStatus(http.StatusInternalServerError)
			return
		}
		h.Data["json"] = hosts
		if len(hosts) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			h.Ctx.Output.SetStatus(http.StatusNotFound)
		} else {
			h.Ctx.Output.SetStatus(http.StatusOK)
		}
	}
}

// @Title listHosts
// @Description get all Host info
// @Success 200
// @router / [get]
func (h *HostsController) GetAll() {
	limit, _ := h.GetInt("limit", 0)
	index, _ := h.GetInt("index", 0)

	defer h.ServeJSON()

	host := &models.Hosts{}
	hosts, err := models.GetHosts(host, limit, index)
	if err != nil {
		h.Data["json"] = map[string]string{
			"message": fmt.Sprint("Failed to get"),
			"error":   err.Error(),
		}
		beego.Warn("[C] Got error:", err)
		h.Ctx.Output.SetStatus(http.StatusInternalServerError)
		return
	}
	h.Data["json"] = hosts
	if len(hosts) == 0 {
		beego.Debug("[C] Got nothing")
		h.Ctx.Output.SetStatus(http.StatusNotFound)
	} else {
		h.Ctx.Output.SetStatus(http.StatusOK)
	}
}

// @Title deleteHost
// @Description delete host
// @Success 204
// @Failure 404
// @router /:name [delete]
func (h *HostsController) Delete() {
	name := h.GetString(":name")
	defer h.ServeJSON()
	beego.Debug("[C] Got name:", name)
	if name != "" {
		host := &models.Hosts{
			Name: name,
		}
		hosts, err := models.GetHosts(host, 0, 0)
		if err != nil {
			h.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			h.Ctx.Output.SetStatus(http.StatusInternalServerError)
			return
		}
		if len(hosts) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			h.Ctx.Output.SetStatus(http.StatusNotFound)
			return
		}
		err = models.DeleteHost(hosts[0])
		if err != nil {
			h.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to delete with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			h.Ctx.Output.SetStatus(http.StatusInternalServerError)
			return
		}
		h.Ctx.Output.SetStatus(http.StatusNoContent)
	}
}

// @Title updateHost
// @Description update host
// @Success 204
// @Failure 404
// @router /:name [put]
func (h *HostsController) Put() {
	name := h.GetString(":name")
	defer h.ServeJSON()
	beego.Debug("[C] Got name:", name)
	if name != "" {
		host := &models.Hosts{
			Name: name,
			// 外键关系也需要初始化，否则会出现问题，反向关系则不用
			AppSet: new(models.AppSets),
		}
		hosts, err := models.GetHosts(host, 0, 0)
		if err != nil {
			h.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			h.Ctx.Output.SetStatus(http.StatusInternalServerError)
			return
		}
		if len(hosts) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			h.Ctx.Output.SetStatus(http.StatusNotFound)
			return
		}

		err = json.Unmarshal(h.Ctx.Input.RequestBody, host)
		if err != nil {
			beego.Warn("[C] Got error:", err)
			h.Data["json"] = map[string]string{
				"message": "Bad request",
				"error":   err.Error(),
			}
			h.Ctx.Output.SetStatus(http.StatusBadRequest)
			return
		}
		host.Id = hosts[0].Id
		beego.Debug("[C] Got host data:", host)
		err = models.UpdateHost(host)
		if err != nil {
			h.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to update with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			h.Ctx.Output.SetStatus(http.StatusInternalServerError)
			return

		}
		h.Ctx.Output.SetStatus(http.StatusAccepted)
	}
}
