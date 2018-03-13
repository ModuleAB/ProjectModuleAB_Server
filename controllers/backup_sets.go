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
	AddPrivilege("GET", "^/api/v1/backupSets", models.RoleFlagUser)
}

type BackupSetsController struct {
	beego.Controller
}

func (h *BackupSetsController) Prepare() {
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

// @Title createBackupSet
// @router / [post]
func (h *BackupSetsController) Post() {
	defer h.ServeJSON()
	backupSet := new(models.BackupSets)
	err := json.Unmarshal(h.Ctx.Input.RequestBody, backupSet)
	if err != nil {
		beego.Warn("[C] Got error:", err)
		h.Data["json"] = map[string]string{
			"message": "Bad request",
			"error":   err.Error(),
		}
		h.Ctx.Output.SetStatus(http.StatusBadRequest)
		return
	}
	beego.Debug("[C] Got data:", backupSet)
	id, err := models.AddBackupSet(backupSet)
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

// @Title getBackupSet
// @router /:name [get]
func (h *BackupSetsController) Get() {
	name := h.GetString(":name")
	defer h.ServeJSON()
	beego.Debug("[C] Got name:", name)
	if name != "" {
		backupSet := &models.BackupSets{
			Name: name,
		}
		backupSets, err := models.GetBackupSets(backupSet, 0, 0)
		if err != nil {
			h.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get  with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			h.Ctx.Output.SetStatus(http.StatusInternalServerError)
			return
		}
		h.Data["json"] = backupSets
		if len(backupSets) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			h.Ctx.Output.SetStatus(http.StatusNotFound)
		} else {
			h.Ctx.Output.SetStatus(http.StatusOK)
		}
	}
}

// @Title listBackupSets
// @router / [get]
func (h *BackupSetsController) GetAll() {
	limit, _ := h.GetInt("limit", 0)
	index, _ := h.GetInt("index", 0)

	defer h.ServeJSON()

	backupSet := &models.BackupSets{}
	backupSets, err := models.GetBackupSets(backupSet, limit, index)
	if err != nil {
		h.Data["json"] = map[string]string{
			"message": fmt.Sprint("Failed to get"),
			"error":   err.Error(),
		}
		beego.Warn("[C] Got error:", err)
		h.Ctx.Output.SetStatus(http.StatusInternalServerError)
		return
	}
	h.Data["json"] = backupSets
	if len(backupSets) == 0 {
		beego.Debug("[C] Got nothing")
		h.Ctx.Output.SetStatus(http.StatusNotFound)
	} else {
		h.Ctx.Output.SetStatus(http.StatusOK)
	}
}

// @Title deleteBackupSet
// @router /:name [delete]
func (h *BackupSetsController) Delete() {
	name := h.GetString(":name")
	defer h.ServeJSON()
	beego.Debug("[C] Got name:", name)
	if name != "" {
		backupSet := &models.BackupSets{
			Name: name,
		}
		backupSets, err := models.GetBackupSets(backupSet, 0, 0)
		if err != nil {
			h.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			h.Ctx.Output.SetStatus(http.StatusInternalServerError)
			return
		}
		if len(backupSets) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			h.Ctx.Output.SetStatus(http.StatusNotFound)
			return
		}
		err = models.DeleteBackupSet(backupSets[0])
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

// @Title updateBackupSet
// @router /:name [put]
func (h *BackupSetsController) Put() {
	name := h.GetString(":name")
	defer h.ServeJSON()
	beego.Debug("[C] Got backupSet name:", name)
	if name != "" {
		backupSet := &models.BackupSets{
			Name: name,
			Oas:  new(models.Oas),
			Oss:  new(models.Oss),
		}
		backupSets, err := models.GetBackupSets(backupSet, 0, 0)
		if err != nil {
			h.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			h.Ctx.Output.SetStatus(http.StatusInternalServerError)
			return
		}
		if len(backupSets) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			h.Ctx.Output.SetStatus(http.StatusNotFound)
			return
		}

		err = json.Unmarshal(h.Ctx.Input.RequestBody, backupSet)
		backupSet.Id = backupSets[0].Id
		if err != nil {
			beego.Warn("[C] Got error:", err)
			h.Data["json"] = map[string]string{
				"message": "Bad request",
				"error":   err.Error(),
			}
			h.Ctx.Output.SetStatus(http.StatusBadRequest)
			return
		}
		beego.Debug("[C] Got backupSet data:", backupSet)
		err = models.UpdateBackupSet(backupSet)
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
