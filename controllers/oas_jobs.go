package controllers

import (
	"fmt"
	"net/http"

	"github.com/ModuleAB/ModuleAB/server/common"
	"github.com/ModuleAB/ModuleAB/server/models"

	"github.com/astaxie/beego"
)

func init() {
	AddPrivilege("GET", "^/api/v1/oasJobs", models.RoleFlagUser)
}

type OasJobsController struct {
	beego.Controller
}

func (h *OasJobsController) Prepare() {
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

// @Title getOAS
// @router /:job_id [get]
func (a *OasJobsController) Get() {
	jobId := a.GetString(":job_id")
	defer a.ServeJSON()
	beego.Debug("[C] Got job id:", jobId)
	if jobId != "" {
		oasJob := &models.OasJobs{
			JobId: jobId,
		}
		oasJobs, err := models.GetOasJobs(oasJob, 0, 0)
		if err != nil {
			a.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with job id:", jobId),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			a.Ctx.Output.SetStatus(http.StatusInternalServerError)
			return
		}
		a.Data["json"] = oasJobs
		if len(oasJobs) == 0 {
			beego.Debug("[C] Got nothing with job id:", jobId)
			a.Ctx.Output.SetStatus(http.StatusNotFound)
		} else {
			a.Ctx.Output.SetStatus(http.StatusOK)
		}
	}
}

// @Title listOAS
// @router / [get]
func (a *OasJobsController) GetAll() {
	limit, _ := a.GetInt("limit", 0)
	index, _ := a.GetInt("index", 0)

	defer a.ServeJSON()

	oasJob := &models.OasJobs{}
	oasJobs, err := models.GetOasJobs(oasJob, limit, index)
	if err != nil {
		a.Data["json"] = map[string]string{
			"message": fmt.Sprint("Failed to get"),
			"error":   err.Error(),
		}
		beego.Warn("[C] Got error:", err)
		a.Ctx.Output.SetStatus(http.StatusInternalServerError)
		return
	}
	a.Data["json"] = oasJobs
	if len(oasJobs) == 0 {
		beego.Debug("[C] Got nothing")
		a.Ctx.Output.SetStatus(http.StatusNotFound)
	} else {
		a.Ctx.Output.SetStatus(http.StatusOK)
	}
}
