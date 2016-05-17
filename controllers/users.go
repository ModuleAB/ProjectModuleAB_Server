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
	AddPrivilege("GET", "^/api/v1/users", models.RoleFlagUser)
	AddPrivilege("PUT", "^/api/v1/users", models.RoleFlagUser)
}

type UserController struct {
	beego.Controller
}

func (h *UserController) Prepare() {
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

// @Title createUser
// @router / [post]
func (h *UserController) Post() {
	user := new(models.Users)
	err := json.Unmarshal(h.Ctx.Input.RequestBody, user)
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
	beego.Debug("[C] Got data:", user)
	id, err := models.AddUser(user)
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

// @Title getUser
// @router /:name [get]
func (h *UserController) Get() {
	name := h.GetString(":name")
	beego.Debug("[C] Got name:", name)
	if name != "" {
		user := &models.Users{
			Name: name,
		}
		users, err := models.GetUser(user, 0, 0)
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
		h.Data["json"] = users
		if len(users) == 0 {
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

// @Title listUser
// @router / [get]
func (h *UserController) GetAll() {
	limit, _ := h.GetInt("limit", 0)
	index, _ := h.GetInt("index", 0)

	user := &models.Users{}
	users, err := models.GetUser(user, limit, index)
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
	h.Data["json"] = users
	if len(users) == 0 {
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

// @Title deleteUser
// @router /:name [delete]
func (h *UserController) Delete() {
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
		err = models.DeleteUser(users[0])
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

// @Title updateUser
// @router /:name [put]
func (h *UserController) Put() {
	name := h.GetString(":name")
	beego.Debug("[C] Got user name:", name)
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

		sessionId := h.GetSession("id")
		if sessionId != nil {
			userNow := &models.Users{
				Id: sessionId.(string),
			}
			userNows, err := models.GetUser(userNow, 0, 0)
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
			if len(userNows) == 0 {
				beego.Debug("[C] Invalid user id:", sessionId)
				h.Ctx.Output.SetStatus(http.StatusNotFound)
				h.ServeJSON()
				return
			}
		}

		err = json.Unmarshal(h.Ctx.Input.RequestBody, user)
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
		user.Id = users[0].Id
		user.Removable = users[0].Removable // Removable should not be changed.
		if user.Password != users[0].Password {
			user.Password = common.EncryptPassword(user.Password)
		}
		beego.Debug("[C] Got user data:", user)
		err = models.UpdateUser(user)
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

// @router /:name/roles [post]
func (h *UserController) AddUserRoles() {
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
