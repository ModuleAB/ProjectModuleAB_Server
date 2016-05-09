package controllers

import (
	"encoding/json"
	"fmt"
	"moduleab_server/common"
	"moduleab_server/models"
	"net/http"

	"github.com/astaxie/beego"
)

type UserController struct {
	beego.Controller
}

// @Title createUser
// @router / [post]
func (a *UserController) Post() {
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

	user := new(models.Users)
	err := json.Unmarshal(a.Ctx.Input.RequestBody, user)
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
	beego.Debug("[C] Got data:", user)
	id, err := models.AddUser(user)
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
	return
}

// @Title getUser
// @router /:name [get]
func (a *UserController) Get() {
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
		user := &models.Users{
			Name: name,
		}
		users, err := models.GetUser(user, 0, 0)
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
		a.Data["json"] = users
		if len(users) == 0 {
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

// @Title listUser
// @router / [get]
func (a *UserController) GetAll() {
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

	user := &models.Users{}
	users, err := models.GetUser(user, limit, index)
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
	a.Data["json"] = users
	if len(users) == 0 {
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

// @Title deleteUser
// @router /:name [delete]
func (a *UserController) Delete() {
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
		user := &models.Users{
			Name: name,
		}
		users, err := models.GetUser(user, 0, 0)
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
		if len(users) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			a.Ctx.Output.SetStatus(http.StatusNotFound)
			a.ServeJSON()
			return
		}
		err = models.DeleteUser(users[0])
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

// @Title updateUser
// @router /:name [put]
func (a *UserController) Put() {
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
	beego.Debug("[C] Got user name:", name)
	if name != "" {
		user := &models.Users{
			Name: name,
		}
		users, err := models.GetUser(user, 0, 0)
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
		if len(users) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			a.Ctx.Output.SetStatus(http.StatusNotFound)
			a.ServeJSON()
			return
		}

		err = json.Unmarshal(a.Ctx.Input.RequestBody, user)
		user.Id = users[0].Id
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
		beego.Debug("[C] Got user data:", user)
		err = models.UpdateUser(user)
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

/*************************************/

// @router /:name/roles [post]
func (h *UserController) AddUserRoles() {
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
			h.Ctx.Output.SetStatus(http.StatusForbidden)
			h.ServeJSON()
		}
	}

	name := h.GetString(":name")
	beego.Debug("[C] Got name:", name)
	if name != "" {
		user := &models.Users{
			Name: name,
		}
		users, err := models.GetUser(user, 0, 0)
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
		if len(users) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			h.Ctx.Output.SetStatus(http.StatusNotFound)
			h.ServeJSON()
			return
		}

		roles := make([]*models.Roles, 0)
		err = json.Unmarshal(h.Ctx.Input.RequestBody, roles)
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
		beego.Debug("[C] Got data:", roles)
		err = models.AddUsersRoles(users[0], roles)
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

// @router /:name/roles [delete]
func (h *UserController) DeleteUserRoles() {
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
			h.Ctx.Output.SetStatus(http.StatusForbidden)
			h.ServeJSON()
		}
	}

	name := h.GetString(":name")
	beego.Debug("[C] Got name:", name)
	if name != "" {
		user := &models.Users{
			Name: name,
		}
		users, err := models.GetUser(user, 0, 0)
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
		if len(users) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			h.Ctx.Output.SetStatus(http.StatusNotFound)
			h.ServeJSON()
			return
		}

		roles := make([]*models.Roles, 0)
		err = json.Unmarshal(h.Ctx.Input.RequestBody, roles)
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
		beego.Debug("[C] Got data:", roles)
		err = models.DeleteUsersRoles(users[0], roles)
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

// @router /login [post]
func (a *UserController) Login() {
	user := new(models.Users)
	err := json.Unmarshal(a.Ctx.Input.RequestBody, user)
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
	beego.Debug("[C] Got data:", user)
	users, err := models.GetUser(user, 0, 0)
	if err != nil {
		a.Data["json"] = map[string]string{
			"message": fmt.Sprint("Failed to get with name:", user.Name),
			"error":   err.Error(),
		}
		beego.Warn("[C] Got error:", err)
		a.Ctx.Output.SetStatus(http.StatusInternalServerError)
		a.ServeJSON()
		return
	}
	if len(users) == 0 {
		beego.Debug("[C] Got nothing with name:", user.Name)
		a.Ctx.Output.SetStatus(http.StatusForbidden)
		a.ServeJSON()
		return
	} else if len(users) > 1 {
		beego.Debug("[C] Got duplicate user with name:", user.Name)
		a.Ctx.Output.SetStatus(http.StatusForbidden)
		a.ServeJSON()
		return
	}
	a.SetSession("id", users[0].Id)
	a.SetSession("name", users[0].Name)
	a.SetSession("show_name", users[0].ShowName)
	a.Ctx.Output.SetStatus(http.StatusOK)
	a.ServeJSON()
}
