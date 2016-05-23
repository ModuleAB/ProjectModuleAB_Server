package controllers

import (
	"encoding/json"
	"fmt"
	"moduleab_server/common"
	"moduleab_server/models"
	"net/http"

	"github.com/astaxie/beego"
)

func init() {
	AddPrivilege("GET", "^/api/v1/oss", models.RoleFlagUser)
}

type OssController struct {
	beego.Controller
}

func (h *OssController) Prepare() {
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

// @Title createOSS
// @router / [post]
func (a *OssController) Post() {
	oss := new(models.Oss)
	defer a.ServeJSON()
	err := json.Unmarshal(a.Ctx.Input.RequestBody, oss)
	if err != nil {
		beego.Warn("[C] Got error:", err)
		a.Data["json"] = map[string]string{
			"message": "Bad request",
			"error":   err.Error(),
		}
		a.Ctx.Output.SetStatus(http.StatusBadRequest)
		return
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
		return
	}

	beego.Debug("[C] Got id:", id)
	a.Data["json"] = map[string]string{
		"id": id,
	}
	a.Ctx.Output.SetStatus(http.StatusCreated)
}

// @Title getOSS
// @router /:name [get]
func (a *OssController) Get() {
	name := a.GetString(":name")
	defer a.ServeJSON()
	beego.Debug("[C] Got name:", name)
	if name != "" {
		oss := &models.Oss{
			BucketName: name,
		}
		osss, err := models.GetOss(oss, 0, 0)
		if err != nil {
			a.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get  with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			a.Ctx.Output.SetStatus(http.StatusInternalServerError)
			return
		}
		a.Data["json"] = osss
		if len(osss) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			a.Ctx.Output.SetStatus(http.StatusNotFound)
		} else {
			a.Ctx.Output.SetStatus(http.StatusOK)
		}
	}
}

// @Title listOSS
// @router / [get]
func (a *OssController) GetAll() {
	limit, _ := a.GetInt("limit", 0)
	index, _ := a.GetInt("index", 0)

	defer a.ServeJSON()

	oss := &models.Oss{}
	osss, err := models.GetOss(oss, limit, index)
	if err != nil {
		a.Data["json"] = map[string]string{
			"message": fmt.Sprint("Failed to get"),
			"error":   err.Error(),
		}
		beego.Warn("[C] Got error:", err)
		a.Ctx.Output.SetStatus(http.StatusInternalServerError)
		return
	}
	a.Data["json"] = osss
	if len(osss) == 0 {
		beego.Debug("[C] Got nothing")
		a.Ctx.Output.SetStatus(http.StatusNotFound)
	} else {
		a.Ctx.Output.SetStatus(http.StatusOK)
	}
}

// @Title deleteOSS
// @router /:name [delete]
func (a *OssController) Delete() {
	name := a.GetString(":name")
	defer a.ServeJSON()
	beego.Debug("[C] Got name:", name)
	if name != "" {
		oss := &models.Oss{
			BucketName: name,
		}
		osss, err := models.GetOss(oss, 0, 0)
		if err != nil {
			a.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			a.Ctx.Output.SetStatus(http.StatusInternalServerError)
			return
		}
		if len(osss) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			a.Ctx.Output.SetStatus(http.StatusNotFound)
			return
		}
		err = models.DeleteOss(osss[0])
		if err != nil {
			a.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to delete with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			a.Ctx.Output.SetStatus(http.StatusInternalServerError)
			return
		}
		a.Ctx.Output.SetStatus(http.StatusNoContent)
	}
}

// @Title updateOSS
// @router /:name [put]
func (a *OssController) Put() {
	name := a.GetString(":name")
	defer a.ServeJSON()
	beego.Debug("[C] Got oss name:", name)
	if name != "" {
		oss := &models.Oss{
			BucketName: name,
		}
		osss, err := models.GetOss(oss, 0, 0)
		if err != nil {
			a.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			a.Ctx.Output.SetStatus(http.StatusInternalServerError)
			return
		}
		if len(osss) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			a.Ctx.Output.SetStatus(http.StatusNotFound)
			return
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
			return
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
			return
		}
		a.Ctx.Output.SetStatus(http.StatusAccepted)
	}
}
