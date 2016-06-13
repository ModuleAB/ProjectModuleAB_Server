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
			for _, host := range p.Hosts {
				for _, path := range p.Paths {
					records, err := models.GetRecords(
						&models.Records{
							BackupSet: p.BackupSet,
							AppSet:    appSet,
							Host:      host,
							Path:      path,
							Type:      p.Target,
						},
						0, 0, models.OrderAsc, models.OrderAsc,
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
					beego.Debug("Got matched records length:", len(records))

					baseLine := records[0]
					for _, r := range records {
						oas, err := common.NewOasClient(r.BackupSet.Oas.Endpoint)
						if err != nil {
							beego.Warn("Cannot connect to OAS Service:", err)
							continue
						}
						oss, err := common.NewOssClient(r.BackupSet.Oss.Endpoint)
						if err != nil {
							beego.Warn("Cannot connect to OSS Service:", err)
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

						switch p.Action {
						case models.PolicyActionArchive:
							beego.Debug("Action type: Archive")
							switch r.Type {
							case models.RecordTypeBackup:
								beego.Debug("Record type: Backup")
								if r.ArchiveId != "" {
									beego.Debug("Record", r.Id, "have archived, skip.")
									continue
								}
								step := r.BackupTime.Sub(baseLine.BackupTime)
								beego.Debug("Step=", step)
								if step >= time.Duration(p.Step)*time.Second &&
									p.Step != models.PolicyReserveNone {
									beego.Debug("New baseline is:", r.Id)
									baseLine = r

									var reqId, jobId string
									beego.Debug(
										"ArchiveToOas:",
										r.BackupSet.Oas.VaultId,
										common.ConvertOssAddrToInternal(
											r.BackupSet.Oss.Endpoint,
										),
										r.BackupSet.Oss.BucketName,
										r.GetFullPath(),
									)
									reqId, jobId, err = oas.ArchiveToOas(
										r.BackupSet.Oas.VaultId,
										common.ConvertOssAddrToInternal(
											r.BackupSet.Oss.Endpoint,
										),
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
									}
								}

							case models.RecordTypeArchive:
								beego.Debug("Record type: Archive")
								beego.Debug("Skip archived data:", r.Id)
							}
						case models.PolicyActionDelete:
							beego.Debug("Action type: Delete")
							switch r.Type {
							case models.RecordTypeBackup:
								beego.Debug("Record type: Backup")
								step := r.BackupTime.Sub(baseLine.BackupTime)
								beego.Debug("Step=", step)
								if step >= time.Duration(p.Step)*time.Second &&
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

								// Don't delete record with ArchiveId,
								// Convert it to Archive.
								if r.ArchiveId == "" {
									err = models.DeleteRecord(r)
									if err != nil {
										beego.Warn(
											"Cannot delete record:", r.Id,
											"error:", err,
										)
									}
								} else {
									r.Type = models.RecordTypeArchive
									err = models.UpdateRecord(r)
									if err != nil {
										beego.Warn(
											"Cannot update archived record:", r.Id,
											"error:", err,
										)
									}
								}

							case models.RecordTypeArchive:
								beego.Debug("Record type: Archive")
								step := r.ArchivedTime.Sub(baseLine.ArchivedTime)
								beego.Debug("Step=", step)
								if step >= time.Duration(p.Step)*time.Second &&
									p.Step != models.PolicyReserveNone {
									beego.Debug("New baseline is:", r.Id)
									baseLine = r
									continue
								}
								beego.Debug("Will delete archive:", r.Id)
								var reqId, jobId string
								reqId, jobId, err = oas.DeleteArchive(
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
								}
							}
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
	period := beego.AppConfig.DefaultInt64("misc::checkoasjobperiod", 5)
	ticker := time.NewTicker(
		time.Duration(period) * time.Minute,
	)
	reservedays := beego.AppConfig.DefaultInt64("misc::oasjobsreservedays", 7)
	defer ticker.Stop()
	beego.Debug("checkOasJob() running...")
	defer beego.Debug("checkOasJob() STOPPED!")
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
					if jl.Completed && !job.Status {
						if jl.StatusCode == "Failed" {
							beego.Warn("Oas job failed:", jl.StatusMessage)
							continue
						}
						job.Status = jl.Completed
						err = models.UpdateOasJobs(job)
						if err != nil {
							beego.Warn("Got error on update oas jobs:", err)
							continue
						}
						record := job.Records
						switch job.JobType {
						case models.OasJobTypePushToOSS:
							beego.Debug("Job type: Push to OSS")
							record.BackupTime = time.Now()
							record.Type = models.RecordTypeBackup

							signal := models.MakeDownloadSignal(
								job.Records.GetFullPath(),
								job.Records.BackupSet.Oss.Endpoint,
								job.Records.BackupSet.Oss.BucketName,
							)
							id, err := models.AddSignal(
								job.Records.Host.Id, signal)
							if err != nil {
								beego.Warn("Got error on add signal:", err)
								continue
							}
							err = models.NotifySignal(
								job.Records.Host.Id, id)
							if err != nil {
								beego.Warn("Got error on push signal:", err)
							}
							err = models.UpdateRecord(record)
							if err != nil {
								beego.Warn(
									"Cannot update record:", record.Id,
									"error:", err,
								)
							}

						case models.OasJobTypePullFromOSS:
							beego.Debug("Job type: Pull from OSS")

							record.ArchiveId = jl.ArchiveId
							record.ArchivedTime = time.Now()

							err = models.UpdateRecord(record)
							if err != nil {
								beego.Warn(
									"Cannot update record:", record.Id,
									"error:", err,
								)
							}

						case models.OasJobTypeDeleteArchive:
							beego.Debug("Job type: Delete archive")
							err = models.DeleteRecord(record)
							if err != nil {
								beego.Warn(
									"Cannot delete record:", record.Id,
									"error:", err,
								)
							}
						}

					} else if job.Status {
						duration := time.Now().Sub(job.CreatedTime)
						if duration > time.Duration(reservedays*24)*time.Hour {
							models.DeleteOasJobs(job)
							beego.Info("Oas job record", job.Id, "is out of date, delete.")
						}
					}
				}
				beego.Info("checkOasJob() completed.")
			}
		}
	}
}
