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
	AddPrivilege("GET", "^/api/v1/clientJobs", models.RoleFlagUser)
}

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

// @Title createClientJob
// @router / [post]
func (a *ClientJobsController) Post() {
	defer a.ServeJSON()
	clientJob := new(models.ClientJobs)
	err := json.Unmarshal(a.Ctx.Input.RequestBody, clientJob)
	if err != nil {
		beego.Warn("[C] Got error:", err)
		a.Data["json"] = map[string]string{
			"message": "Bad request",
			"error":   err.Error(),
		}
		a.Ctx.Output.SetStatus(http.StatusBadRequest)
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
		return
	}

	beego.Debug("[C] Got id:", id)
	a.Data["json"] = map[string]string{
		"id": id,
	}
	a.Ctx.Output.SetStatus(http.StatusCreated)
}

// @Title getClientJob
// @router /:id [get]
func (a *ClientJobsController) Get() {
	id := a.GetString(":id")
	defer a.ServeJSON()
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
			return
		}
		a.Data["json"] = clientJobs
		if len(clientJobs) == 0 {
			beego.Debug("[C] Got nothing with id:", id)
			a.Ctx.Output.SetStatus(http.StatusNotFound)
		} else {
			a.Ctx.Output.SetStatus(http.StatusOK)
		}
	}
}

// @Title listClientJobs
// @router / [get]
func (a *ClientJobsController) GetAll() {
	limit, _ := a.GetInt("limit", 0)
	index, _ := a.GetInt("index", 0)

	defer a.ServeJSON()
	clientJob := &models.ClientJobs{}
	clientJobs, err := models.GetClientJobs(clientJob, limit, index)
	if err != nil {
		a.Data["json"] = map[string]string{
			"message": fmt.Sprint("Failed to get"),
			"error":   err.Error(),
		}
		beego.Warn("[C] Got error:", err)
		a.Ctx.Output.SetStatus(http.StatusInternalServerError)
		return
	}
	a.Data["json"] = clientJobs
	if len(clientJobs) == 0 {
		beego.Debug("[C] Got nothing")
		a.Ctx.Output.SetStatus(http.StatusNotFound)
	} else {
		a.Ctx.Output.SetStatus(http.StatusOK)
	}
}

// @Title deleteClientJob
// @router /:id [delete]
func (a *ClientJobsController) Delete() {
	id := a.GetString(":id")
	defer a.ServeJSON()
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
			return
		}
		if len(clientJobs) == 0 {
			beego.Debug("[C] Got nothing with id:", id)
			a.Ctx.Output.SetStatus(http.StatusNotFound)
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
			return

		}
		a.Ctx.Output.SetStatus(http.StatusNoContent)
	}
}

// @Title updateClientJob
// @router /:id [put]
func (a *ClientJobsController) Put() {
	id := a.GetString(":id")
	defer a.ServeJSON()
	beego.Debug("[C] Got clientJob id:", id)
	if id != "" {
		clientJob := &models.ClientJobs{
			Id:    id,
			Host:  make([]*models.Hosts, 0),
			Paths: make([]*models.Paths, 0),
		}
		clientJobs, err := models.GetClientJobs(clientJob, 0, 0)
		if err != nil {
			a.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with id:", id),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			a.Ctx.Output.SetStatus(http.StatusInternalServerError)
			return
		}
		if len(clientJobs) == 0 {
			beego.Debug("[C] Got nothing with id:", id)
			a.Ctx.Output.SetStatus(http.StatusNotFound)
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
			return
		}
		a.Ctx.Output.SetStatus(http.StatusAccepted)
	}
}
