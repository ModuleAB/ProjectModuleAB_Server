package controllers

import (
	"encoding/json"
	"fmt"
	"moduleab_server/common"
	"moduleab_server/models"
	"net/http"

	"github.com/astaxie/beego"
)

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
		if h.GetSession("id") == nil {
			h.Data["json"] = map[string]string{
				"error": "You need login first.",
			}
			h.Ctx.Output.SetStatus(http.StatusUnauthorized)
			h.ServeJSON()
		}
		if models.CheckPrivileges(
			h.GetSession("id").(string),
			models.RoleFlagOperator,
		) {
			h.Data["json"] = map[string]string{
				"error": "No privilege",
			}
			h.Ctx.Output.SetStatus(http.StatusForbidden)
			h.ServeJSON()
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
	if models.CheckPrivileges(
		h.GetSession("id").(string),
		models.RoleFlagOperator,
	) {
		h.Data["json"] = map[string]string{
			"error": "No privilege",
		}
		h.Ctx.Output.SetStatus(http.StatusForbidden)
		h.ServeJSON()
	}

	host := new(models.Hosts)
	err := json.Unmarshal(h.Ctx.Input.RequestBody, host)
	if err != nil {
		beego.Warn("[C] Got error:", err)
		h.Data["json"] = map[string]string{
			"message": "Bad request",
			"error":   err.Error(),
		}
		h.Ctx.Output.SetStatus(http.StatusBadRequest)
		h.ServeJSON()
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
		h.ServeJSON()
		return
	}

	beego.Debug("[C] Got id:", id)
	h.Data["json"] = map[string]string{
		"id": id,
	}
	h.Ctx.Output.SetStatus(http.StatusCreated)
	h.ServeJSON()
	return
}

// @Title getHost
// @Description get Host info
// @Param	body		body 	models.Hosts	true		"body for host content"
// @Success 200 {string} models.Hosts.Id
// @Failure 403 body is empty
// @router /:name [get]
func (h *HostsController) Get() {
	if models.CheckPrivileges(
		h.GetSession("id").(string),
		models.RoleFlagUser,
	) {
		h.Data["json"] = map[string]string{
			"error": "No privilege",
		}
		h.Ctx.Output.SetStatus(http.StatusForbidden)
		h.ServeJSON()
	}

	name := h.GetString(":name")
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
			h.ServeJSON()
			return
		}
		h.Data["json"] = hosts
		if len(hosts) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			h.Ctx.Output.SetStatus(http.StatusNotFound)
			h.ServeJSON()
			return
		} else {
			h.Ctx.Output.SetStatus(http.StatusOK)
			h.ServeJSON()
			return
		}
	}
}

// @Title listHosts
// @Description get all Host info
// @Success 200
// @router / [get]
func (h *HostsController) GetAll() {
	if models.CheckPrivileges(
		h.GetSession("id").(string),
		models.RoleFlagUser,
	) {
		h.Data["json"] = map[string]string{
			"error": "No privilege",
		}
		h.Ctx.Output.SetStatus(http.StatusForbidden)
		h.ServeJSON()
	}

	limit, _ := h.GetInt("limit", 0)
	index, _ := h.GetInt("index", 0)

	host := &models.Hosts{}
	hosts, err := models.GetHosts(host, limit, index)
	if err != nil {
		h.Data["json"] = map[string]string{
			"message": fmt.Sprint("Failed to get"),
			"error":   err.Error(),
		}
		beego.Warn("[C] Got error:", err)
		h.Ctx.Output.SetStatus(http.StatusInternalServerError)
		h.ServeJSON()
		return
	}
	h.Data["json"] = hosts
	if len(hosts) == 0 {
		beego.Debug("[C] Got nothing")
		h.Ctx.Output.SetStatus(http.StatusNotFound)
		h.ServeJSON()
		return
	} else {
		h.Ctx.Output.SetStatus(http.StatusOK)
		h.ServeJSON()
		return
	}
}

// @Title deleteHost
// @Description delete host
// @Success 204
// @Failure 404
// @router /:name [delete]
func (h *HostsController) Delete() {
	if models.CheckPrivileges(
		h.GetSession("id").(string),
		models.RoleFlagOperator,
	) {
		h.Data["json"] = map[string]string{
			"error": "No privilege",
		}
		h.Ctx.Output.SetStatus(http.StatusForbidden)
		h.ServeJSON()
	}

	name := h.GetString(":name")
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
			h.ServeJSON()
			return
		}
		if len(hosts) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			h.Ctx.Output.SetStatus(http.StatusNotFound)
			h.ServeJSON()
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
			h.ServeJSON()
			return

		}
		h.Ctx.Output.SetStatus(http.StatusNoContent)
		h.ServeJSON()
		return
	}
}

// @Title updateHost
// @Description update host
// @Success 204
// @Failure 404
// @router /:name [put]
func (h *HostsController) Put() {
	if models.CheckPrivileges(
		h.GetSession("id").(string),
		models.RoleFlagOperator,
	) {
		h.Data["json"] = map[string]string{
			"error": "No privilege",
		}
		h.Ctx.Output.SetStatus(http.StatusForbidden)
		h.ServeJSON()
	}

	name := h.GetString(":name")
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
			h.ServeJSON()
			return
		}
		if len(hosts) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			h.Ctx.Output.SetStatus(http.StatusNotFound)
			h.ServeJSON()
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
			h.ServeJSON()
			return
		}
		host.Id = hosts[0].Id
		if host.AppSet != nil {
			host.AppSet = hosts[0].AppSet
		}
		beego.Debug("[C] Got host data:", host)
		err = models.UpdateHost(host)
		if err != nil {
			h.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to update with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			h.Ctx.Output.SetStatus(http.StatusInternalServerError)
			h.ServeJSON()
			return

		}
		h.Ctx.Output.SetStatus(http.StatusAccepted)
		h.ServeJSON()
		return
	}
}

/*****************************************/

// @router /:name/paths [post]
func (h *HostsController) AddHostPaths() {
	if models.CheckPrivileges(
		h.GetSession("id").(string),
		models.RoleFlagOperator,
	) {
		h.Data["json"] = map[string]string{
			"error": "No privilege",
		}
		h.Ctx.Output.SetStatus(http.StatusForbidden)
		h.ServeJSON()
	}

	name := h.GetString(":name")
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
			h.ServeJSON()
			return
		}
		if len(hosts) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			h.Ctx.Output.SetStatus(http.StatusNotFound)
			h.ServeJSON()
			return
		}

		paths := make([]*models.Paths, 0)
		err = json.Unmarshal(h.Ctx.Input.RequestBody, paths)
		if err != nil {
			beego.Warn("[C] Got error:", err)
			h.Data["json"] = map[string]string{
				"message": "Bad request",
				"error":   err.Error(),
			}
			h.Ctx.Output.SetStatus(http.StatusBadRequest)
			h.ServeJSON()
			return
		}
		beego.Debug("[C] Got data:", paths)
		err = models.AddHostPaths(hosts[0], paths)
		if err != nil {
			beego.Warn("[C] Got error:", err)
			h.Data["json"] = map[string]string{
				"message": "Failed to add new path",
				"error":   err.Error(),
			}
			h.Ctx.Output.SetStatus(http.StatusInternalServerError)
			h.ServeJSON()
			return
		}

		h.Ctx.Output.SetStatus(http.StatusCreated)
		h.ServeJSON()
		return
	}
}

// @router /:name/paths [delete]
func (h *HostsController) DeleteHostPaths() {
	if models.CheckPrivileges(
		h.GetSession("id").(string),
		models.RoleFlagOperator,
	) {
		h.Data["json"] = map[string]string{
			"error": "No privilege",
		}
		h.Ctx.Output.SetStatus(http.StatusForbidden)
		h.ServeJSON()
	}

	name := h.GetString(":name")
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
			h.ServeJSON()
			return
		}
		if len(hosts) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			h.Ctx.Output.SetStatus(http.StatusNotFound)
			h.ServeJSON()
			return
		}

		paths := make([]*models.Paths, 0)
		err = json.Unmarshal(h.Ctx.Input.RequestBody, paths)
		if err != nil {
			beego.Warn("[C] Got error:", err)
			h.Data["json"] = map[string]string{
				"message": "Bad request",
				"error":   err.Error(),
			}
			h.Ctx.Output.SetStatus(http.StatusBadRequest)
			h.ServeJSON()
			return
		}
		beego.Debug("[C] Got data:", paths)
		err = models.DeleteHostPaths(hosts[0], paths)
		if err != nil {
			beego.Warn("[C] Got error:", err)
			h.Data["json"] = map[string]string{
				"message": "Failed to delete path",
				"error":   err.Error(),
			}
			h.Ctx.Output.SetStatus(http.StatusInternalServerError)
			h.ServeJSON()
			return
		}

		h.Ctx.Output.SetStatus(http.StatusNoContent)
		h.ServeJSON()
		return
	}
}

/*************************************/

// @router /:name/jobs [post]
func (h *HostsController) AddHostClientJobs() {
	if models.CheckPrivileges(
		h.GetSession("id").(string),
		models.RoleFlagOperator,
	) {
		h.Data["json"] = map[string]string{
			"error": "No privilege",
		}
		h.Ctx.Output.SetStatus(http.StatusForbidden)
		h.ServeJSON()
	}

	name := h.GetString(":name")
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
			h.ServeJSON()
			return
		}
		if len(hosts) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			h.Ctx.Output.SetStatus(http.StatusNotFound)
			h.ServeJSON()
			return
		}

		jobs := make([]*models.ClientJobs, 0)
		err = json.Unmarshal(h.Ctx.Input.RequestBody, jobs)
		if err != nil {
			beego.Warn("[C] Got error:", err)
			h.Data["json"] = map[string]string{
				"message": "Bad request",
				"error":   err.Error(),
			}
			h.Ctx.Output.SetStatus(http.StatusBadRequest)
			h.ServeJSON()
			return
		}
		beego.Debug("[C] Got data:", jobs)
		err = models.AddHostClientJobs(hosts[0], jobs)
		if err != nil {
			beego.Warn("[C] Got error:", err)
			h.Data["json"] = map[string]string{
				"message": "Failed to add new path",
				"error":   err.Error(),
			}
			h.Ctx.Output.SetStatus(http.StatusInternalServerError)
			h.ServeJSON()
			return
		}

		h.Ctx.Output.SetStatus(http.StatusCreated)
		h.ServeJSON()
		return
	}
}

// @router /:name/jobs [delete]
func (h *HostsController) DeleteHostClientJobs() {
	if models.CheckPrivileges(
		h.GetSession("id").(string),
		models.RoleFlagOperator,
	) {
		h.Data["json"] = map[string]string{
			"error": "No privilege",
		}
		h.Ctx.Output.SetStatus(http.StatusForbidden)
		h.ServeJSON()
	}

	name := h.GetString(":name")
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
			h.ServeJSON()
			return
		}
		if len(hosts) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			h.Ctx.Output.SetStatus(http.StatusNotFound)
			h.ServeJSON()
			return
		}

		jobs := make([]*models.ClientJobs, 0)
		err = json.Unmarshal(h.Ctx.Input.RequestBody, jobs)
		if err != nil {
			beego.Warn("[C] Got error:", err)
			h.Data["json"] = map[string]string{
				"message": "Bad request",
				"error":   err.Error(),
			}
			h.Ctx.Output.SetStatus(http.StatusBadRequest)
			h.ServeJSON()
			return
		}
		beego.Debug("[C] Got data:", jobs)
		err = models.DeleteHostClientJobs(hosts[0], jobs)
		if err != nil {
			beego.Warn("[C] Got error:", err)
			h.Data["json"] = map[string]string{
				"message": "Failed to delete path",
				"error":   err.Error(),
			}
			h.Ctx.Output.SetStatus(http.StatusInternalServerError)
			h.ServeJSON()
			return
		}

		h.Ctx.Output.SetStatus(http.StatusNoContent)
		h.ServeJSON()
		return
	}
}
