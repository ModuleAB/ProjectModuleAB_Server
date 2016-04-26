package controllers

import (
	"encoding/json"
	"fmt"
	"moduleab_server/models"
	"net/http"

	"github.com/astaxie/beego"
)

type OssController struct {
	beego.Controller
}

// @router / [post]
func (a *OssController) Post() {
	oss := new(models.Oss)
	err := json.Unmarshal(a.Ctx.Input.RequestBody, oss)
	if err != nil {
		beego.Warn("[C] Got error:", err)
		a.Data["json"] = map[string]string{
			"message": "Bad request",
			"error":   err.Error(),
		}
		a.Ctx.Output.SetStatus(http.StatusBadRequest)
		a.ServeJSON()
	}
	beego.Debug("[C] Got data:", oss)
	id, err := models.AddOss(oss)
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
func (a *OssController) Get() {
	name := a.GetString(":name")
	beego.Debug("[C] Got name:", name)
	if name != "" {
		oss := &models.Oss{
			BucketName: name,
		}
		osss, err := models.GetOss(oss)
		if err != nil {
			a.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get  with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			a.Ctx.Output.SetStatus(http.StatusInternalServerError)
			a.ServeJSON()
		}
		a.Data["json"] = osss
		if len(osss) == 0 {
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
func (a *OssController) GetAll() {
	oss := &models.Oss{}
	osss, err := models.GetOss(oss)
	if err != nil {
		a.Data["json"] = map[string]string{
			"message": fmt.Sprint("Failed to get"),
			"error":   err.Error(),
		}
		beego.Warn("[C] Got error:", err)
		a.Ctx.Output.SetStatus(http.StatusInternalServerError)
		a.ServeJSON()
	}
	a.Data["json"] = osss
	if len(osss) == 0 {
		beego.Debug("[C] Got nothing")
		a.Ctx.Output.SetStatus(http.StatusNotFound)
		a.ServeJSON()
	} else {
		a.Ctx.Output.SetStatus(http.StatusOK)
		a.ServeJSON()
	}
}

// @router /:name [delete]
func (a *OssController) Delete() {
	name := a.GetString(":name")
	beego.Debug("[C] Got name:", name)
	if name != "" {
		oss := &models.Oss{
			BucketName: name,
		}
		osss, err := models.GetOss(oss)
		if err != nil {
			a.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			a.Ctx.Output.SetStatus(http.StatusInternalServerError)
			a.ServeJSON()
		}
		if len(osss) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			a.Ctx.Output.SetStatus(http.StatusNotFound)
			a.ServeJSON()
		}
		err = models.DeleteOss(osss[0])
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
func (a *OssController) Put() {
	name := a.GetString(":name")
	beego.Debug("[C] Got oss name:", name)
	if name != "" {
		oss := &models.Oss{
			BucketName: name,
		}
		osss, err := models.GetOss(oss)
		if err != nil {
			a.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			a.Ctx.Output.SetStatus(http.StatusInternalServerError)
			a.ServeJSON()
		}
		if len(osss) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			a.Ctx.Output.SetStatus(http.StatusNotFound)
			a.ServeJSON()
		}

		err = json.Unmarshal(a.Ctx.Input.RequestBody, oss)
		oss.Id = osss[0].Id
		if err != nil {
			beego.Warn("[C] Got error:", err)
			a.Data["json"] = map[string]string{
				"message": "Bad request",
				"error":   err.Error(),
			}
			a.Ctx.Output.SetStatus(http.StatusBadRequest)
			a.ServeJSON()
		}
		beego.Debug("[C] Got oss data:", oss)
		err = models.UpdateOss(oss)
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
