package controllers

import (
	"fmt"
	"moduleab_server/common"
	"moduleab_server/models"
	"net/http"

	"github.com/astaxie/beego"
)

func init() {
	AddPrivilege("GET", "^/api/v1/roles", models.RoleFlagUser)
}

type RolesController struct {
	beego.Controller
}

func (h *RolesController) Prepare() {
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

// @Title getAllRoles
// @route / [get]
func (h *RolesController) GetAll() {
	limit, _ := h.GetInt("limit", 0)
	index, _ := h.GetInt("index", 0)

	role := &models.Roles{}
	roles, err := models.GetRole(role, limit, index)
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
	h.Data["json"] = roles
	if len(roles) == 0 {
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
