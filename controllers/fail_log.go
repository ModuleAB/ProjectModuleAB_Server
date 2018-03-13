package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	"github.com/ModuleAB/ModuleAB/server/models"
	"github.com/astaxie/beego"
)

func init() {
	AddPrivilege("GET", "^/api/v1/records", models.RoleFlagUser)
}

type FailLogContoller struct {
	beego.Controller
}

/*
func (h *FailLogContoller) Prepare() {
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
*/

// @Title createFailLog
// @Success 201
// @router / [post]
func (h *FailLogContoller) Post() {
	failLog := new(models.FailLog)
	defer h.ServeJSON()
	err := json.Unmarshal(h.Ctx.Input.RequestBody, failLog)
	if err != nil {
		beego.Warn("[C] Got error:", err)
		h.Data["json"] = map[string]string{
			"message": "Bad request",
			"error":   err.Error(),
		}
		h.Ctx.Output.SetStatus(http.StatusBadRequest)
		return
	}
	beego.Debug("[C] Got data:", failLog, failLog.Host)
	id, err := models.AddFailLog(failLog)
	if err != nil {
		beego.Warn("[C] Got error:", err)
		h.Data["json"] = map[string]string{
			"message": "Failed to add New record",
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

	fp, err := filepath.Abs(os.Args[0])
	if err != nil {
		beego.Warn("[C] Got error:", err)
		return
	}
	dir, _ := path.Split(fp)

	cmd := exec.Command(
		path.Join(dir, "alarm"),
		failLog.Host.Name,
		failLog.Log,
	)
	err = cmd.Run()
	if err != nil {
		beego.Warn("[C] Got error:", err)
		return
	}
}

// @Title listFailLog
// @Success 200
// @router / [get]
func (h *FailLogContoller) GetAll() {
	Time, err := time.ParseInLocation("2006-01-02", h.GetString("time"), time.Local)
	beego.Debug("[C] Time:", Time)
	host := h.GetString("host")
	failLog := new(models.FailLog)
	if err == nil {
		failLog.Time = Time
	}
	if host != "" {
		failLog.Host = &models.Hosts{
			Name: host,
		}
	}
	defer h.ServeJSON()
	failLogs, err := models.GetFailLogs(failLog)
	if err != nil {
		h.Data["json"] = map[string]string{
			"message": fmt.Sprint("Failed to get"),
			"error":   err.Error(),
		}
		beego.Warn("[C] Got error:", err)
		h.Ctx.Output.SetStatus(http.StatusInternalServerError)
		return
	}
	h.Data["json"] = failLogs
	if len(failLogs) == 0 {
		beego.Debug("[C] Got nothing")
		h.Ctx.Output.SetStatus(http.StatusNotFound)
	} else {
		h.Ctx.Output.SetStatus(http.StatusOK)
	}
}
