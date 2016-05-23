package policies

import (
	"moduleab_server/common"
	"moduleab_server/models"
	"os"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/pborman/uuid"
	"github.com/robfig/cron"
)

func init() {
	cronSpec := beego.AppConfig.String("misc::policyrun")
	cron := cron.New()
	err := cron.AddFunc(cronSpec, RunPolicies)
	if err != nil {
		beego.Warn("Policy may not be executed for error:", err)
		return

	}
	cron.Start()
}

func RunPolicies() {
	beego.Debug("Policy running...")
	policies, err := models.GetPolicies(&models.Policies{}, 0, 0)
	if err != nil {
		beego.Warn("Run policies error:", err)
		return
	}
	for _, p := range policies {
		beego.Info("Run policy id:", p.Id)
		var backupStart, backupEnd, archiveStart, archiveEnd time.Time
		now := time.Now()
		switch p.Target {
		case models.PolicyTargetBackup:
			if p.TargetEnd != models.PolicyTargetTimeLongLongAgo {
				backupStart = now.Add(
					time.Duration(-p.TargetEnd) * time.Second,
				)
			}
			backupEnd = now.Add(
				time.Duration(-p.TargetStart) * time.Second,
			)
			if backupEnd.Before(backupStart) {
				beego.Warn("End time is before than start time")
				continue
			}
		case models.PolicyTargetArchive:
			if p.TargetEnd != models.PolicyTargetTimeLongLongAgo {
				archiveStart = now.Add(
					time.Duration(-p.TargetEnd) * time.Second,
				)
			}
			archiveEnd = now.Add(
				time.Duration(-p.TargetStart) * time.Second,
			)
			if archiveEnd.Before(archiveStart) {
				beego.Warn("End time is before than start time")
				continue
			}
		}
		for _, appSet := range p.AppSets {
			records, err := models.GetRecords(
				&models.Records{
					BackupSet: p.BackupSet,
					AppSet:    appSet,
					Type:      p.Target,
				},
				0, 0,
				backupStart, backupEnd,
				archiveStart, archiveEnd,
			)
			if err != nil {
				beego.Warn("Get records failed", err)
				continue
			}

			if len(records) == 0 {
				beego.Debug("No records.")
				continue
			}

			baseLine := records[0]
			for _, r := range records {
				oas, err := common.NewOasClient(r.BackupSet.Oas.Endpoint)
				if err != nil {
					beego.Warn("Cannot connect OAS Service:", err)
					continue
				}
				switch p.Action {
				case models.PolicyActionArchive:
					switch r.Type {
					case models.RecordTypeBackup:
						oss, err := common.NewOssClient(r.BackupSet.Oss.Endpoint)
						if err != nil {
							beego.Warn("Cannot connect OSS Service:", err)
							continue
						}
						bucket, err := oss.Bucket(r.BackupSet.Oss.BucketName)
						if err != nil {
							beego.Warn(
								"Cannot get bucket:", r.BackupSet.Oss.BucketName,
								"error:", err,
							)
							continue
						}

						step := r.BackupTime.Sub(baseLine.BackupTime)
						if (step >= time.Duration(p.Step)*time.Second ||
							p.Step == models.PolicyReserveAll) &&
							p.Step != models.PolicyReserveNone {
							baseLine = r
							reqId, jobId, err := oas.ArchiveToOas(
								r.BackupSet.Oas.VaultId,
								r.BackupSet.Oss.Endpoint,
								r.BackupSet.Oss.BucketName,
								r.GetFullPath(),
								r.GetFullPath(),
							)
							if err != nil {
								beego.Warn("Cannot make job to archive:", err)
								continue
							}
							_, err = models.AddOasJobs(
								&models.OasJobs{
									Vault:     r.BackupSet.Oas,
									RequestId: reqId,
									JobId:     jobId,
									JobType:   models.OasJobTypePullFromOSS,
									Status:    models.OasJobStatusIncomplete,
									Records:   r,
								},
							)
							if err != nil {
								beego.Warn("Cannot make oas job:", err)
								continue
							}
						}
						err = bucket.DeleteObject(r.GetFullPath())
						if err != nil {
							beego.Warn(
								"Cannot delete backup:", r.GetFullPath(),
								"error:", err,
							)
							continue
						}
						r.Type = models.RecordTypeArchive
						r.ArchivedTime = time.Now()
						err = models.UpdateRecord(r)
						if err != nil {
							beego.Warn(
								"Cannot update record:", r.Id,
								"error:", err,
							)
							continue
						}

					case models.RecordTypeArchive:
						beego.Debug("Skip archived data:", r.Id)
						continue
					}
				case models.PolicyActionDelete:
					switch r.Type {
					case models.RecordTypeBackup:
						oss, err := common.NewOssClient(r.BackupSet.Oss.Endpoint)
						if err != nil {
							beego.Warn("Cannot connect OSS Service:", err)
							continue
						}
						bucket, err := oss.Bucket(r.BackupSet.Oss.BucketName)
						if err != nil {
							beego.Warn(
								"Cannot get bucket:", r.BackupSet.Oss.BucketName,
								"error:", err,
							)
							continue
						}
						step := r.BackupTime.Sub(baseLine.BackupTime)
						if (step >= time.Duration(p.Step)*time.Second ||
							p.Step == models.PolicyReserveAll) &&
							p.Step != models.PolicyReserveNone {
							baseLine = r
							continue
						}
						err = bucket.DeleteObject(r.GetFullPath())
						if err != nil {
							beego.Warn(
								"Cannot delete backup:", r.GetFullPath(),
								"error:", err,
							)
							continue
						}
						err = models.DeleteRecord(r)
						if err != nil {
							beego.Warn(
								"Cannot delete record:", r.Id,
								"error:", err,
							)
							continue
						}

					case models.RecordTypeArchive:
						step := r.ArchivedTime.Sub(baseLine.ArchivedTime)
						if (step >= time.Duration(p.Step)*time.Second ||
							p.Step == models.PolicyReserveAll) &&
							p.Step != models.PolicyReserveNone {
							baseLine = r
							continue
						}
						reqId, jobId, err := oas.DeleteArchive(
							r.BackupSet.Oas.VaultId,
							r.ArchiveId,
						)
						if err != nil {
							beego.Warn("Cannot make job to delete archive:", err)
							continue
						}
						_, err = models.AddOasJobs(
							&models.OasJobs{
								Vault:     r.BackupSet.Oas,
								RequestId: reqId,
								JobId:     jobId,
								JobType:   models.OasJobTypeDeleteArchive,
								Status:    models.OasJobStatusIncomplete,
								Records:   r,
							},
						)
						if err != nil {
							beego.Warn("Cannot make oas job:", err)
							continue
						}
					}
				}
			}
		}
		beego.Info("Policy id", p.Id, "Done.")
	}
}

func InitDb() {
	o := orm.NewOrm()

	role := []models.Roles{
		models.Roles{
			Id:        uuid.New(),
			Name:      "Administrator",
			RoleFlag:  models.RoleFlagAdmin,
			Removable: false,
		},
		models.Roles{
			Id:        uuid.New(),
			Name:      "Operator",
			RoleFlag:  models.RoleFlagOperator,
			Removable: false,
		},
		models.Roles{
			Id:        uuid.New(),
			Name:      "User",
			RoleFlag:  models.RoleFlagUser,
			Removable: false,
		},
	}
	o.Begin()
	_, err := o.InsertMulti(1, role)
	if err != nil {
		o.Rollback()
		beego.Alert("Error on inserting roles:", err)
		os.Exit(1)
	}
	o.Commit()

	user := &models.Users{
		Id:       uuid.New(),
		Name:     "admin",
		ShowName: "Administrator",
		Password: "admin",
		Roles: []*models.Roles{
			&role[0],
		},
		Removable: false,
	}
	_, err = models.AddUser(user)
	if err != nil {
		beego.Alert("Error on inserting user:", err)
		os.Exit(1)
	}

	appSet := &models.AppSets{
		Name: "Default",
		Desc: "Default app set",
	}
	_, err = models.AddAppSet(appSet)
	if err != nil {
		beego.Alert("Error on inserting default application set:", err)
		os.Exit(1)
	}

	backupSet := &models.BackupSets{
		Id:   uuid.New(),
		Name: "Default",
		Desc: "Default backup set",
	}
	_, err = models.AddBackupSet(backupSet)
	if err != nil {
		o.Rollback()
		beego.Alert("Error on inserting default backup set:", err)
		os.Exit(1)
	}
}

func CheckOasJob() {
	period, err := beego.AppConfig.Int64("misc::checkoasjobperiod")
	if err != nil {
		period = 5
	}
	ticker := time.NewTicker(
		time.Duration(period) * time.Minute,
	)
	defer ticker.Stop()
	beego.Debug("checkOasJob() running...")
	for {
		select {
		case <-ticker.C:
			beego.Info("checkOasJob() start.")
			oas, err := models.GetOas(&models.Oas{}, 0, 0)
			if err != nil {
				beego.Warn("Got error on retrieving OAS records:", err)
				continue
			}
			for _, v := range oas {
				beego.Debug("Got oas:", v)
				o, err := common.NewOasClient(v.Endpoint)
				if err != nil {
					beego.Warn("Got error on connecting to OAS:", err)
					continue
				}

				jobCond := &models.OasJobs{
					Vault: v,
				}
				jobs, err := models.GetOasJobs(jobCond, 0, 0)
				if err != nil {
					beego.Warn("Got error on retrieving oas jobs:", err)
					continue
				}

				for _, job := range jobs {
					beego.Debug("Got job:", job)
					_, jl, err := o.GetJobInfo(
						job.Vault.VaultId,
						job.JobId,
					)
					if err != nil {
						beego.Warn("Got error on retrieving job info:", err)
						continue
					}
					if jl.Completed {
						job.Status = jl.Completed
						err = models.UpdateOasJobs(job)
						if err != nil {
							beego.Warn("Got error on update oas jobs:", err)
							continue
						}
						switch job.JobType {
						case models.OasJobTypePushToOSS:
							record := job.Records
							record.BackupTime = time.Now()
							err := models.UpdateRecord(record)
							if err != nil {
								beego.Warn("Got error on updating record:", err)
								continue
							}

							signal := models.MakeDownloadSignal(
								job.Records.GetFullPath(),
								job.Records.BackupSet.Oss.Endpoint,
								job.Records.BackupSet.Oss.BucketName,
							)
							id, _ := models.AddSignal(
								job.Records.Host.Id, signal)
							err = models.NotifySignal(
								job.Records.Host.Id, id)
							if err != nil {
								beego.Warn("Got error on push signal:", err)
							}
						}
					}
				}
				beego.Info("checkOasJob() completed.")
			}
		}
	}
	defer beego.Debug("checkOasJob() STOPPED!")
}
