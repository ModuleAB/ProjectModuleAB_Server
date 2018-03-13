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
	AddPrivilege("GET", "^/api/v1/oas", models.RoleFlagUser)
}

type OasController struct {
	beego.Controller
}

func (h *OasController) Prepare() {
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

// @Title createOAS
// @router / [post]
func (a *OasController) Post() {
	defer a.ServeJSON()
	oas := new(models.Oas)
	err := json.Unmarshal(a.Ctx.Input.RequestBody, oas)
	if err != nil {
		beego.Warn("[C] Got error:", err)
		a.Data["json"] = map[string]string{
			"message": "Bad request",
			"error":   err.Error(),
		}
		a.Ctx.Output.SetStatus(http.StatusBadRequest)
		return
	}
	beego.Debug("Got data:", oas)

	o, err := common.NewOasClient(oas.Endpoint)
	if err != nil {
		beego.Warn("[C] Got error:", err)
		a.Data["json"] = map[string]string{
			"message": "Bad config",
			"error":   err.Error(),
		}
		a.Ctx.Output.SetStatus(http.StatusInternalServerError)
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
		return
	}

	beego.Debug("[C] Got id:", id)
	a.Data["json"] = map[string]string{
		"id": id,
	}
	a.Ctx.Output.SetStatus(http.StatusCreated)
}

// @Title getOAS
// @router /:name [get]
func (a *OasController) Get() {
	name := a.GetString(":name")
	defer a.ServeJSON()
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
			return
		}
		a.Data["json"] = oass
		if len(oass) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			a.Ctx.Output.SetStatus(http.StatusNotFound)
		} else {
			a.Ctx.Output.SetStatus(http.StatusOK)
		}
	}
}

// @Title listOAS
// @router / [get]
func (a *OasController) GetAll() {
	limit, _ := a.GetInt("limit", 0)
	index, _ := a.GetInt("index", 0)

	defer a.ServeJSON()

	oas := &models.Oas{}
	oass, err := models.GetOas(oas, limit, index)
	if err != nil {
		a.Data["json"] = map[string]string{
			"message": fmt.Sprint("Failed to get"),
			"error":   err.Error(),
		}
		beego.Warn("[C] Got error:", err)
		a.Ctx.Output.SetStatus(http.StatusInternalServerError)
		return
	}
	a.Data["json"] = oass
	if len(oass) == 0 {
		beego.Debug("[C] Got nothing")
		a.Ctx.Output.SetStatus(http.StatusNotFound)
	} else {
		a.Ctx.Output.SetStatus(http.StatusOK)
	}
}

// @Title deleteOAS
// @router /:name [delete]
func (a *OasController) Delete() {
	name := a.GetString(":name")
	defer a.ServeJSON()
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
			return
		}
		if len(oass) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			a.Ctx.Output.SetStatus(http.StatusNotFound)
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
			return

		}
		a.Ctx.Output.SetStatus(http.StatusNoContent)
	}
}

// @Title updateOAS
// @router /:name [put]
func (a *OasController) Put() {
	name := a.GetString(":name")
	defer a.ServeJSON()
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
			return
		}
		if len(oass) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			a.Ctx.Output.SetStatus(http.StatusNotFound)
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
			return
		}
		a.Ctx.Output.SetStatus(http.StatusAccepted)
	}
}
