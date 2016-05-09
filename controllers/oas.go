package controllers

import (
	"encoding/json"
	"fmt"
	"moduleab_server/common"
	"moduleab_server/models"
	"net/http"

	"github.com/astaxie/beego"
)

type OasController struct {
	beego.Controller
}

// @Title createOAS
// @router / [post]
func (a *OasController) Post() {
	if a.Ctx.Input.Header("Signature") != "" {
		err := common.AuthWithKey(a.Ctx)
		if err != nil {
			a.Data["json"] = map[string]string{
				"error": err.Error(),
			}
			a.Ctx.Output.SetStatus(http.StatusForbidden)
			a.ServeJSON()
		}
	} else {
		if a.GetSession("id") == nil {
			a.Data["json"] = map[string]string{
				"error": "You need login first.",
			}
			a.Ctx.Output.SetStatus(http.StatusForbidden)
			a.ServeJSON()
		}
	}

	oas := new(models.Oas)
	err := json.Unmarshal(a.Ctx.Input.RequestBody, oas)
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

	o, err := common.NewOasClient(oas.Endpoint)
	if err != nil {
		beego.Warn("[C] Got error:", err)
		a.Data["json"] = map[string]string{
			"message": "Bad config",
			"error":   err.Error(),
		}
		a.Ctx.Output.SetStatus(http.StatusInternalServerError)
		a.ServeJSON()
		return
	}
	oas.VaultId, err = o.GetOasVaultId(oas.VaultName)
	if err != nil {
		beego.Warn("[C] Got error:", err)
		a.Data["json"] = map[string]string{
			"message": "Failed to access OAS",
			"error":   err.Error(),
		}
		a.Ctx.Output.SetStatus(http.StatusInternalServerError)
		a.ServeJSON()
		return
	}

	beego.Debug("[C] Got data:", oas)
	id, err := models.AddOas(oas)
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
}

// @Title getOAS
// @router /:name [get]
func (a *OasController) Get() {
	if a.Ctx.Input.Header("Signature") != "" {
		err := common.AuthWithKey(a.Ctx)
		if err != nil {
			a.Data["json"] = map[string]string{
				"error": err.Error(),
			}
			a.Ctx.Output.SetStatus(http.StatusForbidden)
			a.ServeJSON()
		}
	} else {
		if a.GetSession("id") == nil {
			a.Data["json"] = map[string]string{
				"error": "You need login first.",
			}
			a.Ctx.Output.SetStatus(http.StatusForbidden)
			a.ServeJSON()
		}
	}

	name := a.GetString(":name")
	beego.Debug("[C] Got name:", name)
	if name != "" {
		oas := &models.Oas{
			VaultName: name,
		}
		oass, err := models.GetOas(oas, 0, 0)
		if err != nil {
			a.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get  with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			a.Ctx.Output.SetStatus(http.StatusInternalServerError)
			a.ServeJSON()
			return
		}
		a.Data["json"] = oass
		if len(oass) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
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

// @Title listOAS
// @router / [get]
func (a *OasController) GetAll() {
	if a.Ctx.Input.Header("Signature") != "" {
		err := common.AuthWithKey(a.Ctx)
		if err != nil {
			a.Data["json"] = map[string]string{
				"error": err.Error(),
			}
			a.Ctx.Output.SetStatus(http.StatusForbidden)
			a.ServeJSON()
		}
	} else {
		if a.GetSession("id") == nil {
			a.Data["json"] = map[string]string{
				"error": "You need login first.",
			}
			a.Ctx.Output.SetStatus(http.StatusForbidden)
			a.ServeJSON()
		}
	}

	limit, _ := a.GetInt("limit", 0)
	index, _ := a.GetInt("index", 0)

	oas := &models.Oas{}
	oass, err := models.GetOas(oas, limit, index)
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
	a.Data["json"] = oass
	if len(oass) == 0 {
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

// @Title deleteOAS
// @router /:name [delete]
func (a *OasController) Delete() {
	if a.Ctx.Input.Header("Signature") != "" {
		err := common.AuthWithKey(a.Ctx)
		if err != nil {
			a.Data["json"] = map[string]string{
				"error": err.Error(),
			}
			a.Ctx.Output.SetStatus(http.StatusForbidden)
			a.ServeJSON()
		}
	} else {
		if a.GetSession("id") == nil {
			a.Data["json"] = map[string]string{
				"error": "You need login first.",
			}
			a.Ctx.Output.SetStatus(http.StatusForbidden)
			a.ServeJSON()
		}
	}

	name := a.GetString(":name")
	beego.Debug("[C] Got name:", name)
	if name != "" {
		oas := &models.Oas{
			VaultName: name,
		}
		oass, err := models.GetOas(oas, 0, 0)
		if err != nil {
			a.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			a.Ctx.Output.SetStatus(http.StatusInternalServerError)
			a.ServeJSON()
			return
		}
		if len(oass) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			a.Ctx.Output.SetStatus(http.StatusNotFound)
			a.ServeJSON()
			return
		}
		err = models.DeleteOas(oass[0])
		if err != nil {
			a.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to delete with name:", name),
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

// @Title updateOAS
// @router /:name [put]
func (a *OasController) Put() {
	if a.Ctx.Input.Header("Signature") != "" {
		err := common.AuthWithKey(a.Ctx)
		if err != nil {
			a.Data["json"] = map[string]string{
				"error": err.Error(),
			}
			a.Ctx.Output.SetStatus(http.StatusForbidden)
			a.ServeJSON()
		}
	} else {
		if a.GetSession("id") == nil {
			a.Data["json"] = map[string]string{
				"error": "You need login first.",
			}
			a.Ctx.Output.SetStatus(http.StatusForbidden)
			a.ServeJSON()
		}
	}

	name := a.GetString(":name")
	beego.Debug("[C] Got oas name:", name)
	if name != "" {
		oas := &models.Oas{
			VaultName: name,
		}
		oass, err := models.GetOas(oas, 0, 0)
		if err != nil {
			a.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			a.Ctx.Output.SetStatus(http.StatusInternalServerError)
			a.ServeJSON()
			return
		}
		if len(oass) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			a.Ctx.Output.SetStatus(http.StatusNotFound)
			a.ServeJSON()
			return
		}

		err = json.Unmarshal(a.Ctx.Input.RequestBody, oas)
		oas.Id = oass[0].Id
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
		beego.Debug("[C] Got oas data:", oas)
		err = models.UpdateOas(oas)
		if err != nil {
			a.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to update with name:", name),
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
