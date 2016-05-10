package controllers

import (
	"encoding/json"
	"fmt"
	"moduleab_server/common"
	"moduleab_server/models"
	"net/http"

	"github.com/astaxie/beego"
)

type ClientJobsController struct {
	beego.Controller
}

func (h *ClientJobsController) Prepare() {
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

// @Title createClientJob
// @router / [post]
func (a *ClientJobsController) Post() {
	if models.CheckPrivileges(
		a.GetSession("id").(string),
		models.RoleFlagOperator,
	) {
		a.Data["json"] = map[string]string{
			"error": "No privilege",
		}
		a.Ctx.Output.SetStatus(http.StatusForbidden)
		a.ServeJSON()
	}

	clientJob := new(models.ClientJobs)
	err := json.Unmarshal(a.Ctx.Input.RequestBody, clientJob)
	if err != nil {
		beego.Warn("[C] Got error:", err)
		a.Data["json"] = map[string]string{
			"message": "Bad request",
			"error":   err.Error(),
		}
		a.Ctx.Output.SetStatus(http.StatusBadRequest)
		a.ServeJSON()
		return
	}
	beego.Debug("[C] Got data:", clientJob)
	id, err := models.AddClientJob(clientJob)
	if err != nil {
		beego.Warn("[C] Got error:", err)
		a.Data["json"] = map[string]string{
			"message": "Failed to add New host",
			"error":   err.Error(),
		}
		a.Ctx.Output.SetStatus(http.StatusInternalServerError)
		a.ServeJSON()
		return
	}

	beego.Debug("[C] Got id:", id)
	a.Data["json"] = map[string]string{
		"id": id,
	}
	a.Ctx.Output.SetStatus(http.StatusCreated)
	a.ServeJSON()
	return
}

// @Title getClientJob
// @router /:id [get]
func (a *ClientJobsController) Get() {
	if models.CheckPrivileges(
		a.GetSession("id").(string),
		models.RoleFlagUser,
	) {
		a.Data["json"] = map[string]string{
			"error": "No privilege",
		}
		a.Ctx.Output.SetStatus(http.StatusForbidden)
		a.ServeJSON()
	}

	id := a.GetString(":id")
	beego.Debug("[C] Got id:", id)
	if id != "" {
		clientJob := &models.ClientJobs{
			Id: id,
		}
		clientJobs, err := models.GetClientJobs(clientJob, 0, 0)
		if err != nil {
			a.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get  with id:", id),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			a.Ctx.Output.SetStatus(http.StatusInternalServerError)
			a.ServeJSON()
			return
		}
		a.Data["json"] = clientJobs
		if len(clientJobs) == 0 {
			beego.Debug("[C] Got nothing with id:", id)
			a.Ctx.Output.SetStatus(http.StatusNotFound)
			a.ServeJSON()
			return
		} else {
			a.Ctx.Output.SetStatus(http.StatusOK)
			a.ServeJSON()
			return
		}
	}
}

// @Title listClientJobs
// @router / [get]
func (a *ClientJobsController) GetAll() {
	if models.CheckPrivileges(
		a.GetSession("id").(string),
		models.RoleFlagUser,
	) {
		a.Data["json"] = map[string]string{
			"error": "No privilege",
		}
		a.Ctx.Output.SetStatus(http.StatusForbidden)
		a.ServeJSON()
	}

	limit, _ := a.GetInt("limit", 0)
	index, _ := a.GetInt("index", 0)

	clientJob := &models.ClientJobs{}
	clientJobs, err := models.GetClientJobs(clientJob, limit, index)
	if err != nil {
		a.Data["json"] = map[string]string{
			"message": fmt.Sprint("Failed to get"),
			"error":   err.Error(),
		}
		beego.Warn("[C] Got error:", err)
		a.Ctx.Output.SetStatus(http.StatusInternalServerError)
		a.ServeJSON()
		return
	}
	a.Data["json"] = clientJobs
	if len(clientJobs) == 0 {
		beego.Debug("[C] Got nothing")
		a.Ctx.Output.SetStatus(http.StatusNotFound)
		a.ServeJSON()
		return
	} else {
		a.Ctx.Output.SetStatus(http.StatusOK)
		a.ServeJSON()
		return
	}
}

// @Title deleteClientJob
// @router /:id [delete]
func (a *ClientJobsController) Delete() {
	if models.CheckPrivileges(
		a.GetSession("id").(string),
		models.RoleFlagOperator,
	) {
		a.Data["json"] = map[string]string{
			"error": "No privilege",
		}
		a.Ctx.Output.SetStatus(http.StatusForbidden)
		a.ServeJSON()
	}

	id := a.GetString(":id")
	beego.Debug("[C] Got id:", id)
	if id != "" {
		clientJob := &models.ClientJobs{
			Id: id,
		}
		clientJobs, err := models.GetClientJobs(clientJob, 0, 0)
		if err != nil {
			a.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with id:", id),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			a.Ctx.Output.SetStatus(http.StatusInternalServerError)
			a.ServeJSON()
			return
		}
		if len(clientJobs) == 0 {
			beego.Debug("[C] Got nothing with id:", id)
			a.Ctx.Output.SetStatus(http.StatusNotFound)
			a.ServeJSON()
			return
		}
		err = models.DeleteClientJob(clientJobs[0])
		if err != nil {
			a.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to delete with id:", id),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			a.Ctx.Output.SetStatus(http.StatusInternalServerError)
			a.ServeJSON()
			return

		}
		a.Ctx.Output.SetStatus(http.StatusNoContent)
		a.ServeJSON()
		return
	}
}

// @Title updateClientJob
// @router /:id [put]
func (a *ClientJobsController) Put() {
	if models.CheckPrivileges(
		a.GetSession("id").(string),
		models.RoleFlagOperator,
	) {
		a.Data["json"] = map[string]string{
			"error": "No privilege",
		}
		a.Ctx.Output.SetStatus(http.StatusForbidden)
		a.ServeJSON()
	}

	id := a.GetString(":id")
	beego.Debug("[C] Got clientJob id:", id)
	if id != "" {
		clientJob := &models.ClientJobs{
			Id: id,
		}
		clientJobs, err := models.GetClientJobs(clientJob, 0, 0)
		if err != nil {
			a.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with id:", id),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			a.Ctx.Output.SetStatus(http.StatusInternalServerError)
			a.ServeJSON()
			return
		}
		if len(clientJobs) == 0 {
			beego.Debug("[C] Got nothing with id:", id)
			a.Ctx.Output.SetStatus(http.StatusNotFound)
			a.ServeJSON()
			return
		}

		err = json.Unmarshal(a.Ctx.Input.RequestBody, clientJob)
		clientJob.Id = clientJobs[0].Id
		if err != nil {
			beego.Warn("[C] Got error:", err)
			a.Data["json"] = map[string]string{
				"message": "Bad request",
				"error":   err.Error(),
			}
			a.Ctx.Output.SetStatus(http.StatusBadRequest)
			a.ServeJSON()
			return
		}
		beego.Debug("[C] Got clientJob data:", clientJob)
		err = models.UpdateClientJob(clientJob)
		if err != nil {
			a.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to update with id:", id),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			a.Ctx.Output.SetStatus(http.StatusInternalServerError)
			a.ServeJSON()
			return

		}
		a.Ctx.Output.SetStatus(http.StatusAccepted)
		a.ServeJSON()
		return
	}
}
