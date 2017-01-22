package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ModuleAB/ModuleAB/server/common"
	"github.com/ModuleAB/ModuleAB/server/models"

	"github.com/astaxie/beego"
)

func init() {
	AddPrivilege("GET", "^/api/v1/records", models.RoleFlagUser)
}

type RecordsController struct {
	beego.Controller
}

func (h *RecordsController) Prepare() {
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

// @Title createRecord
// @Description create Record
// @Param	record 	body 	object true	"record"
// @Success 201 {object} models.Records
// @Failure 400 Recordname or IP missing
// @Failure 500 Failure on writing database
// @router / [post]
func (h *RecordsController) Post() {
	record := new(models.Records)
	defer h.ServeJSON()
	err := json.Unmarshal(h.Ctx.Input.RequestBody, record)
	if err != nil {
		beego.Warn("[C] Got error:", err)
		h.Data["json"] = map[string]string{
			"message": "Bad request",
			"error":   err.Error(),
		}
		h.Ctx.Output.SetStatus(http.StatusBadRequest)
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
		return
	}

	beego.Debug("[C] Got id:", id)
	h.Data["json"] = map[string]string{
		"id": id,
	}
	h.Ctx.Output.SetStatus(http.StatusCreated)
}

// @Title listRecords
// @Description get all Record info
// @Success 200
// @router / [get]
func (h *RecordsController) GetAll() {
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
		models.OrderDesc, models.OrderDesc,
		tBtStart, tBtEnd, tAtStart, tAtEnd)
	defer h.ServeJSON()
	if err != nil {
		h.Data["json"] = map[string]string{
			"message": fmt.Sprint("Failed to get"),
			"error":   err.Error(),
		}
		beego.Warn("[C] Got error:", err)
		h.Ctx.Output.SetStatus(http.StatusInternalServerError)
		return
	}
	h.Data["json"] = records
	if len(records) == 0 {
		beego.Debug("[C] Got nothing")
		h.Ctx.Output.SetStatus(http.StatusNotFound)
	} else {
		h.Ctx.Output.SetStatus(http.StatusOK)
	}
}

// @Title deleteRecord
// @Description delete record
// @Success 204
// @Failure 404
// @router /:id [delete]
func (h *RecordsController) Delete() {
	id := h.GetString(":id")
	beego.Debug("[C] Got id:", id)
	defer h.ServeJSON()
	if id != "" {
		record := &models.Records{
			Id: id,
		}
		records, err := models.GetRecords(record, 0, 0,
			models.OrderAsc, models.OrderAsc)
		if err != nil {
			h.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with id:", id),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			h.Ctx.Output.SetStatus(http.StatusInternalServerError)
			return
		}
		if len(records) == 0 {
			beego.Debug("[C] Got nothing with id:", id)
			h.Ctx.Output.SetStatus(http.StatusNotFound)
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
			return
		}
		h.Ctx.Output.SetStatus(http.StatusNoContent)
		return
	}
}

// @router /:id/recover [get]
func (h *RecordsController) Recover() {
	id := h.GetString(":id")
	beego.Debug("[C] Got id:", id)
	defer h.ServeJSON()
	if id != "" {
		record := &models.Records{
			Id: id,
		}
		records, err := models.GetRecords(record, 0, 0,
			models.OrderAsc, models.OrderAsc)
		if err != nil {
			h.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with id:", id),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			h.Ctx.Output.SetStatus(http.StatusInternalServerError)
			return
		}
		if len(records) == 0 {
			beego.Debug("[C] Got nothing with id:", id)
			h.Ctx.Output.SetStatus(http.StatusNotFound)
			return
		}

		switch records[0].Type {
		case models.RecordTypeArchive:
			oasJob := &models.OasJobs{
				JobType: models.OasJobTypePushToOSS,
				Vault:   records[0].BackupSet.Oas,
				Records: records[0],
			}

			oasClient, err := common.NewOasClient(
				records[0].BackupSet.Oas.Endpoint,
			)
			if err != nil {
				h.Data["json"] = map[string]string{
					"message": fmt.Sprint("Failed to connect to OAS"),
					"error":   err.Error(),
				}
				beego.Warn("[C] Got error:", err)
				h.Ctx.Output.SetStatus(http.StatusInternalServerError)
				return
			}

			oasJob.RequestId, oasJob.JobId, err = oasClient.RecoverToOss(
				records[0].BackupSet.Oas.VaultId,
				records[0].ArchiveId,
				common.ConvertOssAddrToInternal(
					records[0].BackupSet.Oss.Endpoint,
				),
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
			id, err := models.AddSignal(
				records[0].Host.Id,
				signal,
			)
			if err != nil {
				beego.Warn("[C] Got error:", err)
			}

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
				return
			}
			h.Data["json"] = map[string]string{
				"job_id":  id,
				"message": "This is backup, so agent should be downloading now.",
			}
			h.Ctx.Output.SetStatus(http.StatusOK)
		}
	}
}
