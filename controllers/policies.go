package controllers

import (
	"encoding/json"
	"fmt"
	"moduleab_server/common"
	"moduleab_server/models"
	"net/http"

	"github.com/astaxie/beego"
)

type PolicyController struct {
	beego.Controller
}

// @Title createOAS
// @router / [post]
func (a *PolicyController) Post() {
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
			a.Ctx.Output.SetStatus(http.StatusUnauthorized)
			a.ServeJSON()
		}
		if models.CheckPrivileges(
			a.GetSession("id").(string),
			models.RoleFlagOperator,
		) {
			a.Data["json"] = map[string]string{
				"error": "No privilege",
			}
			a.Ctx.Output.SetStatus(http.StatusForbidden)
			a.ServeJSON()
		}
	}

	policy := new(models.Policies)
	err := json.Unmarshal(a.Ctx.Input.RequestBody, policy)
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

	beego.Debug("[C] Got data:", policy)
	id, err := models.AddPolicy(policy)
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

// @Title getPolicy
// @router /:name [get]
func (a *PolicyController) Get() {
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
			a.Ctx.Output.SetStatus(http.StatusUnauthorized)
			a.ServeJSON()
		}
		if models.CheckPrivileges(
			a.GetSession("id").(string),
			models.RoleFlagUser,
		) {
			a.Data["json"] = map[string]string{
				"error": "No privilege",
			}
			a.Ctx.Output.SetStatus(http.StatusForbidden)
			a.ServeJSON()
		}
	}

	name := a.GetString(":name")
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
			a.ServeJSON()
			return
		}
		a.Data["json"] = policies
		if len(policies) == 0 {
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

// @Title listPolicies
// @router / [get]
func (a *PolicyController) GetAll() {
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
			a.Ctx.Output.SetStatus(http.StatusUnauthorized)
			a.ServeJSON()
		}
		if models.CheckPrivileges(
			a.GetSession("id").(string),
			models.RoleFlagUser,
		) {
			a.Data["json"] = map[string]string{
				"error": "No privilege",
			}
			a.Ctx.Output.SetStatus(http.StatusForbidden)
			a.ServeJSON()
		}
	}

	limit, _ := a.GetInt("limit", 0)
	index, _ := a.GetInt("index", 0)

	policy := &models.Policies{}
	policies, err := models.GetPolicies(policy, limit, index)
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
	a.Data["json"] = policies
	if len(policies) == 0 {
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
func (a *PolicyController) Delete() {
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
			a.Ctx.Output.SetStatus(http.StatusUnauthorized)
			a.ServeJSON()
		}
		if models.CheckPrivileges(
			a.GetSession("id").(string),
			models.RoleFlagOperator,
		) {
			a.Data["json"] = map[string]string{
				"error": "No privilege",
			}
			a.Ctx.Output.SetStatus(http.StatusForbidden)
			a.ServeJSON()
		}
	}

	name := a.GetString(":name")
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
			a.ServeJSON()
			return
		}
		if len(policies) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			a.Ctx.Output.SetStatus(http.StatusNotFound)
			a.ServeJSON()
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
func (a *PolicyController) Put() {
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
			a.Ctx.Output.SetStatus(http.StatusUnauthorized)
			a.ServeJSON()
		}
		if models.CheckPrivileges(
			a.GetSession("id").(string),
			models.RoleFlagOperator,
		) {
			a.Data["json"] = map[string]string{
				"error": "No privilege",
			}
			a.Ctx.Output.SetStatus(http.StatusForbidden)
			a.ServeJSON()
		}
	}

	name := a.GetString(":name")
	beego.Debug("[C] Got policy name:", name)
	if name != "" {
		policy := &models.Policies{
			Name:      name,
			AppSet:    new(models.AppSets),
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
			a.ServeJSON()
			return
		}
		if len(policies) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			a.Ctx.Output.SetStatus(http.StatusNotFound)
			a.ServeJSON()
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
			a.ServeJSON()
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
			a.ServeJSON()
			return

		}
		a.Ctx.Output.SetStatus(http.StatusAccepted)
		a.ServeJSON()
		return
	}
}
