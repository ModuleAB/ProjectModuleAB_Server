package controllers

import (
	"encoding/json"
	"fmt"
	"moduleab_server/models"
	"net/http"

	"github.com/astaxie/beego"
)

type LoginController struct {
	beego.Controller
}

// @router /login [post]
func (h *LoginController) Login() {
	defer h.ServeJSON()
	user := new(models.Users)
	err := json.Unmarshal(h.Ctx.Input.RequestBody, user)
	if err != nil {
		beego.Warn("[C] Got error:", err)
		h.Data["json"] = map[string]string{
			"message": "Bad request",
			"error":   err.Error(),
		}
		h.Ctx.Output.SetStatus(http.StatusBadRequest)
		return
	}
	beego.Debug("[C] Got data:", user)
	users, err := models.GetUser(user, 0, 0)
	if err != nil {
		h.Data["json"] = map[string]string{
			"message": fmt.Sprint("Failed to get with name:", user.Name),
			"error":   err.Error(),
		}
		beego.Warn("[C] Got error:", err)
		h.Ctx.Output.SetStatus(http.StatusInternalServerError)
		return
	}
	if len(users) == 0 || user.Password == "" {
		beego.Debug("[C] Got nothing with name:", user.Name)
		h.Ctx.Output.SetStatus(http.StatusForbidden)
		return
	} else if len(users) > 1 {
		beego.Debug("[C] Got duplicate user with name:", user.Name)
		h.Ctx.Output.SetStatus(http.StatusForbidden)
		return
	}
	h.SetSession("id", users[0].Id)
	h.SetSession("name", users[0].Name)
	h.SetSession("show_name", users[0].ShowName)
	h.Ctx.Output.SetStatus(http.StatusOK)
}

// @router /logout [get]
func (h *LoginController) Logout() {
	defer h.ServeJSON()
	if h.GetSession("id") == nil {
		h.Data["json"] = map[string]string{
			"error": "You need login first.",
		}
		h.Ctx.Output.SetStatus(http.StatusUnauthorized)
		return
	}

	h.DelSession("id")
	h.Ctx.Output.SetStatus(http.StatusOK)
}

// @router /check [get]
func (h *LoginController) Check() {
	id := h.GetSession("id")
	defer h.ServeJSON()
	if id == nil {
		h.Data["json"] = map[string]string{
			"error": "You need login first.",
		}
		h.Ctx.Output.SetStatus(http.StatusUnauthorized)
		return
	} else if _, ok := id.(string); !ok {
		h.Data["json"] = map[string]string{
			"error": "Invalid user id.",
		}
		h.Ctx.Output.SetStatus(http.StatusUnauthorized)
		return
	}
	h.Data["json"] = map[string]interface{}{
		"id":        id,
		"name":      h.GetSession("name"),
		"show_name": h.GetSession("show_name"),
	}
	h.Ctx.Output.SetStatus(http.StatusOK)
}
