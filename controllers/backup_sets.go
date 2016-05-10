package controllers

import (
	"encoding/json"
	"fmt"
	"moduleab_server/common"
	"moduleab_server/models"
	"net/http"

	"github.com/astaxie/beego"
)

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

// @Title createBackupSet
// @router / [post]
func (h *BackupSetsController) Post() {
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

	backupSet := new(models.BackupSets)
	err := json.Unmarshal(h.Ctx.Input.RequestBody, backupSet)
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
	beego.Debug("[C] Got data:", backupSet)
	id, err := models.AddBackupSet(backupSet)
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

// @Title getBackupSet
// @router /:name [get]
func (h *BackupSetsController) Get() {
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
			h.ServeJSON()
			return
		}
		h.Data["json"] = backupSets
		if len(backupSets) == 0 {
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

// @Title listBackupSets
// @router / [get]
func (h *BackupSetsController) GetAll() {
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

	backupSet := &models.BackupSets{}
	backupSets, err := models.GetBackupSets(backupSet, limit, index)
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
	h.Data["json"] = backupSets
	if len(backupSets) == 0 {
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

// @Title deleteBackupSet
// @router /:name [delete]
func (h *BackupSetsController) Delete() {
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
			h.ServeJSON()
			return
		}
		if len(backupSets) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			h.Ctx.Output.SetStatus(http.StatusNotFound)
			h.ServeJSON()
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
			h.ServeJSON()
			return

		}
		h.Ctx.Output.SetStatus(http.StatusNoContent)
		h.ServeJSON()
		return
	}
}

// @Title updateBackupSet
// @router /:name [put]
func (h *BackupSetsController) Put() {
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
			h.ServeJSON()
			return
		}
		if len(backupSets) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			h.Ctx.Output.SetStatus(http.StatusNotFound)
			h.ServeJSON()
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
			h.ServeJSON()
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
			h.ServeJSON()
			return

		}
		h.Ctx.Output.SetStatus(http.StatusAccepted)
		h.ServeJSON()
		return
	}
}

/*************************************/

// @router /:name/hosts [post]
func (h *BackupSetsController) AddBackupSetsHosts() {
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
			h.ServeJSON()
			return
		}
		if len(backupSets) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			h.Ctx.Output.SetStatus(http.StatusNotFound)
			h.ServeJSON()
			return
		}

		hosts := make([]*models.Hosts, 0)
		err = json.Unmarshal(h.Ctx.Input.RequestBody, hosts)
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
		beego.Debug("[C] Got data:", hosts)
		err = models.AddBackupSetsHosts(backupSets[0], hosts)
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

		h.Ctx.Output.SetStatus(http.StatusNoContent)
		h.ServeJSON()
		return
	}
}

// @router /:name/hosts [delete]
func (h *BackupSetsController) DeleteBackupSetsHosts() {
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
			h.ServeJSON()
			return
		}
		if len(backupSets) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			h.Ctx.Output.SetStatus(http.StatusNotFound)
			h.ServeJSON()
			return
		}

		hosts := make([]*models.Hosts, 0)
		err = json.Unmarshal(h.Ctx.Input.RequestBody, hosts)
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
		beego.Debug("[C] Got data:", hosts)
		err = models.DeleteBackupSetsHosts(backupSets[0], hosts)
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
