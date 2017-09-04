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
	AddPrivilege("GET", "^/api/v1/paths", models.RoleFlagUser)
}

type PathsController struct {
	beego.Controller
}

func (h *PathsController) Prepare() {
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

// @Title createPath
// @Description create Path
// @Param	path 	body 	object true	"path"
// @Success 201 {object} models.Paths
// @Failure 400 Pathname or IP missing
// @Failure 500 Failure on writing database
// @router / [post]
func (h *PathsController) Post() {
	defer h.ServeJSON()
	path := new(models.Paths)
	err := json.Unmarshal(h.Ctx.Input.RequestBody, path)
	if err != nil {
		beego.Warn("[C] Got error:", err)
		h.Data["json"] = map[string]string{
			"message": "Bad request",
			"error":   err.Error(),
		}
		h.Ctx.Output.SetStatus(http.StatusBadRequest)
		return
	}
	beego.Debug("[C] Got data:", path)
	id, err := models.AddPath(path)
	if err != nil {
		beego.Warn("[C] Got error:", err)
		h.Data["json"] = map[string]string{
			"message": "Failed to add New path",
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

// @Title getPath
// @Description get Path info, only id works
// @Param	body		body 	models.Paths	true		"body for path content"
// @Success 200 {string} models.Paths.Id
// @Failure 403 body is empty
// @router /:id [get]
func (h *PathsController) Get() {
	id := h.GetString(":id")
	defer h.ServeJSON()
	beego.Debug("[C] Got id:", id)
	if id != "" {
		path := &models.Paths{
			Id: id,
		}
		paths, err := models.GetPaths(path, 0, 0)
		if err != nil {
			h.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with id:", id),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			h.Ctx.Output.SetStatus(http.StatusInternalServerError)
			return
		}
		h.Data["json"] = paths
		if len(paths) == 0 {
			beego.Debug("[C] Got nothing with id:", id)
			h.Ctx.Output.SetStatus(http.StatusNotFound)
		} else {
			h.Ctx.Output.SetStatus(http.StatusOK)
		}
	}
}

// @Title listPaths
// @Description get all Path info
// @Success 200
// @router / [get]
func (h *PathsController) GetAll() {
	limit, _ := h.GetInt("limit", 0)
	index, _ := h.GetInt("index", 0)

	defer h.ServeJSON()

	path := &models.Paths{}
	paths, err := models.GetPaths(path, limit, index)
	if err != nil {
		h.Data["json"] = map[string]string{
			"message": fmt.Sprint("Failed to get"),
			"error":   err.Error(),
		}
		beego.Warn("[C] Got error:", err)
		h.Ctx.Output.SetStatus(http.StatusInternalServerError)
		return
	}
	h.Data["json"] = paths
	if len(paths) == 0 {
		beego.Debug("[C] Got nothing")
		h.Ctx.Output.SetStatus(http.StatusNotFound)
	} else {
		h.Ctx.Output.SetStatus(http.StatusOK)
	}
}

// @Title deletePath
// @Description delete path
// @Success 204
// @Failure 404
// @router /:id [delete]
func (h *PathsController) Delete() {
	id := h.GetString(":id")
	defer h.ServeJSON()
	beego.Debug("[C] Got id:", id)
	if id != "" {
		path := &models.Paths{
			Id: id,
		}
		paths, err := models.GetPaths(path, 0, 0)
		if err != nil {
			h.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with id:", id),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			h.Ctx.Output.SetStatus(http.StatusInternalServerError)
			return
		}
		if len(paths) == 0 {
			beego.Debug("[C] Got nothing with id:", id)
			h.Ctx.Output.SetStatus(http.StatusNotFound)
			return
		}
		err = models.DeletePath(paths[0])
		if err != nil {
			h.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to delete with id:", id),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			h.Ctx.Output.SetStatus(http.StatusInternalServerError)
			return
		}
		h.Ctx.Output.SetStatus(http.StatusNoContent)
	}
}

// @Title updatePath
// @Description update path
// @Success 204
// @Failure 404
// @router /:id [put]
func (h *PathsController) Put() {
	id := h.GetString(":id")
	defer h.ServeJSON()
	beego.Debug("[C] Got id:", id)
	if id != "" {
		path := &models.Paths{
			Id: id,
			// 外键关系也需要初始化，否则会出现问题，反向关系则不用
			BackupSet: new(models.BackupSets),
			AppSet:    make([]*models.AppSets, 0),
		}
		paths, err := models.GetPaths(path, 0, 0)
		if err != nil {
			h.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with id:", id),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			h.Ctx.Output.SetStatus(http.StatusInternalServerError)
			return
		}
		if len(paths) == 0 {
			beego.Debug("[C] Got nothing with id:", id)
			h.Ctx.Output.SetStatus(http.StatusNotFound)
			return
		}

		err = json.Unmarshal(h.Ctx.Input.RequestBody, path)
		if err != nil {
			beego.Warn("[C] Got error:", err)
			h.Data["json"] = map[string]string{
				"message": "Bad request",
				"error":   err.Error(),
			}
			h.Ctx.Output.SetStatus(http.StatusBadRequest)
			return
		}
		path.Id = paths[0].Id
		if path.AppSet == nil || len(path.AppSet) == 0 {
			path.AppSet = paths[0].AppSet
		}
		beego.Debug("[C] Got path data:", path)
		err = models.UpdatePath(path)
		if err != nil {
			h.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to update with id:", id),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			h.Ctx.Output.SetStatus(http.StatusInternalServerError)
			return
		}
		h.Ctx.Output.SetStatus(http.StatusAccepted)
	}
}
