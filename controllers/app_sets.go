package controllers

import (
	"encoding/json"
	"fmt"
	"moduleab_server/models"
	"net/http"

	"github.com/astaxie/beego"
)

type AppSetsController struct {
	beego.Controller
}

// @router / [post]
func (a *AppSetsController) Post() {
	appSet := new(models.AppSets)
	err := json.Unmarshal(a.Ctx.Input.RequestBody, appSet)
	if err != nil {
		beego.Warn("[C] Got error:", err)
		a.Data["json"] = map[string]string{
			"message": "Bad request",
			"error":   err.Error(),
		}
		a.Ctx.Output.SetStatus(http.StatusBadRequest)
		a.ServeJSON()
	}
	beego.Debug("[C] Got data:", appSet)
	id, err := models.AddAppSet(appSet)
	if err != nil {
		beego.Warn("[C] Got error:", err)
		a.Data["json"] = map[string]string{
			"message": "Failed to add New host",
			"error":   err.Error(),
		}
		a.Ctx.Output.SetStatus(http.StatusInternalServerError)
		a.ServeJSON()
	}

	beego.Debug("[C] Got id:", id)
	a.Data["json"] = map[string]string{
		"id": id,
	}
	a.Ctx.Output.SetStatus(http.StatusCreated)
	a.ServeJSON()
}

// @router /:name [get]
func (a *AppSetsController) Get() {
	name := a.GetString(":name")
	beego.Debug("[C] Got name:", name)
	if name != "" {
		appSet := &models.AppSets{
			Name: name,
		}
		appSets, err := models.GetAppSets(appSet)
		if err != nil {
			a.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get  with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			a.Ctx.Output.SetStatus(http.StatusInternalServerError)
			a.ServeJSON()
		}
		a.Data["json"] = appSets
		if len(appSets) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			a.Ctx.Output.SetStatus(http.StatusNotFound)
			a.ServeJSON()
		} else {
			a.Ctx.Output.SetStatus(http.StatusOK)
			a.ServeJSON()
		}
	}
}

// @router / [get]
func (a *AppSetsController) GetAll() {
	appSet := &models.AppSets{}
	appSets, err := models.GetAppSets(appSet)
	if err != nil {
		a.Data["json"] = map[string]string{
			"message": fmt.Sprint("Failed to get"),
			"error":   err.Error(),
		}
		beego.Warn("[C] Got error:", err)
		a.Ctx.Output.SetStatus(http.StatusInternalServerError)
		a.ServeJSON()
	}
	a.Data["json"] = appSets
	if len(appSets) == 0 {
		beego.Debug("[C] Got nothing")
		a.Ctx.Output.SetStatus(http.StatusNotFound)
		a.ServeJSON()
	} else {
		a.Ctx.Output.SetStatus(http.StatusOK)
		a.ServeJSON()
	}
}

// @router /:name [delete]
func (a *AppSetsController) Delete() {
	name := a.GetString(":name")
	beego.Debug("[C] Got name:", name)
	if name != "" {
		appSet := &models.AppSets{
			Name: name,
		}
		appSets, err := models.GetAppSets(appSet)
		if err != nil {
			a.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			a.Ctx.Output.SetStatus(http.StatusInternalServerError)
			a.ServeJSON()
		}
		if len(appSets) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			a.Ctx.Output.SetStatus(http.StatusNotFound)
			a.ServeJSON()
		}
		err = models.DeleteAppSet(appSets[0])
		if err != nil {
			a.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to delete with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			a.Ctx.Output.SetStatus(http.StatusInternalServerError)
			a.ServeJSON()

		}
		a.Ctx.Output.SetStatus(http.StatusNoContent)
		a.ServeJSON()
	}
}

// @router /:name [put]
func (a *AppSetsController) Put() {
	name := a.GetString(":name")
	beego.Debug("[C] Got appSet name:", name)
	if name != "" {
		appSet := &models.AppSets{
			Name: name,
		}
		appSets, err := models.GetAppSets(appSet)
		if err != nil {
			a.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			a.Ctx.Output.SetStatus(http.StatusInternalServerError)
			a.ServeJSON()
		}
		if len(appSets) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			a.Ctx.Output.SetStatus(http.StatusNotFound)
			a.ServeJSON()
		}

		err = json.Unmarshal(a.Ctx.Input.RequestBody, appSet)
		appSet.Id = appSets[0].Id
		if err != nil {
			beego.Warn("[C] Got error:", err)
			a.Data["json"] = map[string]string{
				"message": "Bad request",
				"error":   err.Error(),
			}
			a.Ctx.Output.SetStatus(http.StatusBadRequest)
			a.ServeJSON()
		}
		beego.Debug("[C] Got appSet data:", appSet)
		err = models.UpdateAppSet(appSet)
		if err != nil {
			a.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to update with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			a.Ctx.Output.SetStatus(http.StatusInternalServerError)
			a.ServeJSON()

		}
		a.Ctx.Output.SetStatus(http.StatusAccepted)
		a.ServeJSON()
	}
}
