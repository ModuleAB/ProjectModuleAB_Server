package controllers

import (
	"encoding/json"
	"fmt"
	"moduleab_server/common"
	"moduleab_server/models"
	"net/http"
	"time"

	"github.com/astaxie/beego"
)

type RecordsController struct {
	beego.Controller
}

// @Title createRecord
// @Description create Record
// @Param	record 	body 	object true	"record"
// @Success 201 {object} models.Records
// @Failure 400 Recordname or IP missing
// @Failure 500 Failure on writing database
// @router / [post]
func (h *RecordsController) Post() {
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

	record := new(models.Records)
	err := json.Unmarshal(h.Ctx.Input.RequestBody, record)
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
	beego.Debug("[C] Got data:", record)
	id, err := models.AddRecord(record)
	if err != nil {
		beego.Warn("[C] Got error:", err)
		h.Data["json"] = map[string]string{
			"message": "Failed to add New record",
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

// @Title listRecords
// @Description get all Record info
// @Success 200
// @router / [get]
func (h *RecordsController) GetAll() {
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
			models.RoleFlagUser,
		) {
			h.Data["json"] = map[string]string{
				"error": "No privilege",
			}
			h.Ctx.Output.SetStatus(http.StatusForbidden)
			h.ServeJSON()
		}
	}

	limit, _ := h.GetInt("limit", 50)
	index, _ := h.GetInt("index", 0)
	filename := h.GetString("filename")
	path := h.GetString("path")
	archiveId := h.GetString("archiveId")
	appSet := h.GetString("appSet")
	backupSet := h.GetString("backupSet")
	host := h.GetString("host")
	// Format: RFC3339
	btStart := h.GetString("btStart")
	btEnd := h.GetString("btEnd")
	atStart := h.GetString("atStart")
	atEnd := h.GetString("atEnd")

	record := &models.Records{
		Path: &models.Paths{
			Path: path,
		},
		Filename:  filename,
		ArchiveId: archiveId,
		Host: &models.Hosts{
			Name: host,
		},
		AppSet: &models.AppSets{
			Name: appSet,
		},
		BackupSet: &models.BackupSets{
			Name: backupSet,
		},
	}

	tBtStart, _ := time.Parse(time.RFC3339, btStart)
	tBtEnd, _ := time.Parse(time.RFC3339, btEnd)
	tAtStart, _ := time.Parse(time.RFC3339, atStart)
	tAtEnd, _ := time.Parse(time.RFC3339, atEnd)
	records, err := models.GetRecords(record, limit, index,
		tBtStart, tBtEnd, tAtStart, tAtEnd)
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
	h.Data["json"] = records
	if len(records) == 0 {
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

// @Title deleteRecord
// @Description delete record
// @Success 204
// @Failure 404
// @router /:id [delete]
func (h *RecordsController) Delete() {
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

	id := h.GetString(":id")
	beego.Debug("[C] Got id:", id)
	if id != "" {
		record := &models.Records{
			Id: id,
		}
		records, err := models.GetRecords(record, 0, 0)
		if err != nil {
			h.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with id:", id),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			h.Ctx.Output.SetStatus(http.StatusInternalServerError)
			h.ServeJSON()
			return
		}
		if len(records) == 0 {
			beego.Debug("[C] Got nothing with id:", id)
			h.Ctx.Output.SetStatus(http.StatusNotFound)
			h.ServeJSON()
			return
		}
		err = models.DeleteRecord(records[0])
		if err != nil {
			h.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to delete with id:", id),
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

// @router /:id/recover [get]
func (h *RecordsController) Recover() {
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
			models.RoleFlagUser,
		) {
			h.Data["json"] = map[string]string{
				"error": "No privilege",
			}
			h.Ctx.Output.SetStatus(http.StatusForbidden)
			h.ServeJSON()
		}
	}

	id := h.GetString(":id")
	beego.Debug("[C] Got id:", id)
	if id != "" {
		record := &models.Records{
			Id: id,
		}
		records, err := models.GetRecords(record, 0, 0)
		if err != nil {
			h.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with id:", id),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			h.Ctx.Output.SetStatus(http.StatusInternalServerError)
			h.ServeJSON()
			return
		}
		if len(records) == 0 {
			beego.Debug("[C] Got nothing with id:", id)
			h.Ctx.Output.SetStatus(http.StatusNotFound)
			h.ServeJSON()
			return
		}

		switch records[0].Type {
		case models.RecordTypeArchive:
			oasJob := &models.OasJobs{
				JobType: models.OasJobTypePushToOSS,
				Vault:   records[0].BackupSet.Oas,
				Records: records[0],
			}
			oasJob.RequestId, oasJob.JobId, err = common.DefaultOasClient.ArchiveToOas(
				records[0].BackupSet.Oas.VaultId,
				records[0].BackupSet.Oss.Endpoint,
				records[0].BackupSet.Oss.BucketName,
				records[0].GetFullPath(),
				records[0].GetFullPath(),
			)
			if err != nil {
				h.Data["json"] = map[string]string{
					"message": fmt.Sprint("Failed commit job:", id),
					"error":   err.Error(),
				}
				beego.Warn("[C] Got error:", err)
				h.Ctx.Output.SetStatus(http.StatusInternalServerError)
				h.ServeJSON()
				return
			}
			id, err := models.AddOasJobs(oasJob)
			if err != nil {
				h.Data["json"] = map[string]string{
					"message": fmt.Sprint("Failed to record added oas job"),
					"error":   err.Error(),
				}
				beego.Warn("[C] Got error:", err)
				h.Ctx.Output.SetStatus(http.StatusInternalServerError)
				h.ServeJSON()
				return
			}
			h.Data["json"] = map[string]string{
				"job_id":  id,
				"message": "This is archive, so some waiting is necessary.",
			}
			h.Ctx.Output.SetStatus(http.StatusAccepted)
		case models.RecordTypeBackup:
			signal := models.MakeDownloadSignal(
				records[0].GetFullPath(),
				records[0].BackupSet.Oss.Endpoint,
				records[0].BackupSet.Oss.BucketName,
			)
			id, _ := models.AddSignal(
				records[0].Host.Id,
				signal,
			)
			err = models.NotifySignal(
				records[0].Host.Id,
				id,
			)
			if err != nil {
				h.Data["json"] = map[string]string{
					"message": fmt.Sprint("Failed to recover file with record:", id),
					"error":   err.Error(),
				}
				beego.Warn("[C] Got error:", err)
				h.Ctx.Output.SetStatus(http.StatusInternalServerError)
				h.ServeJSON()
				return
			}
			h.Data["json"] = map[string]string{
				"job_id":  id,
				"message": "This is backup, so agent should be downloading now.",
			}
			h.Ctx.Output.SetStatus(http.StatusOK)
		}

		h.ServeJSON()
		return
	}
}
