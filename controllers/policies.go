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
	AddPrivilege("GET", "^/api/v1/policies", models.RoleFlagUser)
}

type PolicyController struct {
	beego.Controller
}

func (h *PolicyController) Prepare() {
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
func (a *PolicyController) Post() {
	defer a.ServeJSON()
	policy := new(models.Policies)
	err := json.Unmarshal(a.Ctx.Input.RequestBody, policy)
	if err != nil {
		beego.Warn("[C] Got error:", err)
		a.Data["json"] = map[string]string{
			"message": "Bad request",
			"error":   err.Error(),
		}
		a.Ctx.Output.SetStatus(http.StatusBadRequest)
		return
	}

	beego.Debug("[C] Got data:", policy)
	id, err := models.AddPolicy(policy)
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

// @Title getPolicy
// @router /:name [get]
func (a *PolicyController) Get() {
	name := a.GetString(":name")
	defer a.ServeJSON()
	beego.Debug("[C] Got name:", name)
	if name != "" {
		policy := &models.Policies{
			Name: name,
		}
		policies, err := models.GetPolicies(policy, 0, 0)
		if err != nil {
			a.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get  with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			a.Ctx.Output.SetStatus(http.StatusInternalServerError)
			return
		}
		a.Data["json"] = policies
		if len(policies) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			a.Ctx.Output.SetStatus(http.StatusNotFound)
		} else {
			a.Ctx.Output.SetStatus(http.StatusOK)
		}
	}
}

// @Title listPolicies
// @router / [get]
func (a *PolicyController) GetAll() {
	limit, _ := a.GetInt("limit", 0)
	index, _ := a.GetInt("index", 0)
	defer a.ServeJSON()
	policy := &models.Policies{}
	policies, err := models.GetPolicies(policy, limit, index)
	if err != nil {
		a.Data["json"] = map[string]string{
			"message": fmt.Sprint("Failed to get"),
			"error":   err.Error(),
		}
		beego.Warn("[C] Got error:", err)
		a.Ctx.Output.SetStatus(http.StatusInternalServerError)
		return
	}
	a.Data["json"] = policies
	if len(policies) == 0 {
		beego.Debug("[C] Got nothing")
		a.Ctx.Output.SetStatus(http.StatusNotFound)
	} else {
		a.Ctx.Output.SetStatus(http.StatusOK)
	}
}

// @Title deleteOAS
// @router /:name [delete]
func (a *PolicyController) Delete() {
	name := a.GetString(":name")
	a.ServeJSON()
	beego.Debug("[C] Got name:", name)
	if name != "" {
		policy := &models.Policies{
			Name: name,
		}
		policies, err := models.GetPolicies(policy, 0, 0)
		if err != nil {
			a.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			a.Ctx.Output.SetStatus(http.StatusInternalServerError)
			return
		}
		if len(policies) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			a.Ctx.Output.SetStatus(http.StatusNotFound)
			return
		}
		err = models.DeletePolicy(policies[0])
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
func (a *PolicyController) Put() {
	name := a.GetString(":name")
	defer a.ServeJSON()
	beego.Debug("[C] Got policy name:", name)
	if name != "" {
		policy := &models.Policies{
			Name:      name,
			AppSets:   make([]*models.AppSets, 0),
			Paths:     make([]*models.Paths, 0),
			Hosts:     make([]*models.Hosts, 0),
			BackupSet: new(models.BackupSets),
		}
		policies, err := models.GetPolicies(policy, 0, 0)
		if err != nil {
			a.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			a.Ctx.Output.SetStatus(http.StatusInternalServerError)
			return
		}
		if len(policies) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			a.Ctx.Output.SetStatus(http.StatusNotFound)
			return
		}

		err = json.Unmarshal(a.Ctx.Input.RequestBody, policy)
		policy.Id = policies[0].Id
		if err != nil {
			beego.Warn("[C] Got error:", err)
			a.Data["json"] = map[string]string{
				"message": "Bad request",
				"error":   err.Error(),
			}
			a.Ctx.Output.SetStatus(http.StatusBadRequest)
			return
		}
		beego.Debug("[C] Got policy data:", policy)
		err = models.UpdatePolicy(policy)
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
