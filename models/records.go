// record.go
package models

import (
	"fmt"
	"moduleab_server/common"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"github.com/pborman/uuid"
)

const (
	RecordTypeAll     = PolicyTargetAll
	RecordTypeBackup  = PolicyTargetBackup
	RecordTypeArchive = PolicyTargetArchive
)

const (
	backupTimeStart = iota
	backupTimeEnd
	archiveTimeStart
	archiveTimeEnd
)

type Records struct {
	Id           string      `orm:"pk;size(36)" json:"id" valid:"Match(/^[A-Fa-f0-9]{8}-([A-Fa-f0-9]{4}-){3}[A-Fa-f0-9]{12}$/)"`
	Host         *Hosts      `orm:"rel(fk)" json:"host"`
	BackupSet    *BackupSets `orm:"rel(fk)" json:"backupset"`
	AppSet       *AppSets    `orm:"rel(fk)" json:"appset"`
	Path         *Paths      `orm:"rel(fk)" json:"path" valid:"Required"`
	Filename     string      `orm:"key" json:"filename" valid:"Required"`
	Type         int         `json:"type" valid:"Required"` // 0 - All, 1 - Backup, 2 - Archive
	ArchiveId    string      `orm:"null" json:"archiveid"`  // 如果Type是1（归档）时，这里应该有数据
	BackupTime   time.Time   `orm:"type(datetime)" json:"backuptime"`
	ArchivedTime time.Time   `orm:"type(datatime);null" json:"archivedtime"`
	Jobs         []*OasJobs  `orm:"reverse(many);null" json:"jobs"`
}

func (r *Records) GetFullPath() string {
	return fmt.Sprintf("%s/%s%s/%s",
		r.AppSet.Name,
		r.Host.Name,
		r.Path.Path,
		r.Filename,
	)
}

func init() {
	if prefix := beego.AppConfig.String("database::mysqlprefex"); prefix != "" {
		orm.RegisterModelWithPrefix(prefix, new(Records))
	} else {
		orm.RegisterModel(new(Records))
	}
}

func AddRecord(record *Records) (string, error) {
	beego.Debug("[M] Got data:", record)
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		return "", err
	}

	record.Id = uuid.New()
	beego.Debug("[M] Got id:", record.Id)
	validator := new(validation.Validation)
	valid, err := validator.Valid(record)
	if err != nil {
		o.Rollback()
		return "", err
	}
	if !valid {
		o.Rollback()
		var errS string
		for _, err := range validator.Errors {
			errS = fmt.Sprintf("%s, %s:%s", errS, err.Key, err.Message)
		}
		return "", fmt.Errorf("Bad info: %s", errS)
	}

	beego.Debug("[M] Got new data:", record)
	_, err = o.Insert(record)
	if err != nil {
		o.Rollback()
		return "", err
	}
	beego.Debug("[M] Record data saved")
	o.Commit()
	return record.Id, nil

}

func DeleteRecord(h *Records) error {
	beego.Debug("[M] Got data:", h)
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		return err
	}
	validator := new(validation.Validation)
	valid, err := validator.Valid(h)
	if err != nil {
		o.Rollback()
		return err
	}
	if !valid {
		o.Rollback()
		var errS string
		for _, err := range validator.Errors {
			errS = fmt.Sprintf("%s, %s:%s", errS, err.Key, err.Message)
		}
		return fmt.Errorf("Bad info: %s", errS)
	}
	_, err = o.Delete(h)
	if err != nil {
		o.Rollback()
		return err
	}
	o.Commit()
	return nil
}

func UpdateRecord(h *Records) error {
	beego.Debug("[M] Got data:", h)
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		return err
	}
	validator := new(validation.Validation)
	valid, err := validator.Valid(h)
	if err != nil {
		o.Rollback()
		return err
	}
	if !valid {
		o.Rollback()
		var errS string
		for _, err := range validator.Errors {
			errS = fmt.Sprintf("%s, %s:%s", errS, err.Key, err.Message)
		}
		return fmt.Errorf("Bad info: %s", errS)
	}
	_, err = o.Update(h)
	if err != nil {
		o.Rollback()
		return err
	}
	o.Commit()
	return nil
}

// If get all, just use &Record{}
func GetRecords(cond *Records, limit, index int,
	times ...time.Time) ([]*Records, error) {
	r := make([]*Records, 0)
	o := orm.NewOrm()
	q := o.QueryTable("records")
	if cond.Id != "" {
		q = q.Filter("id", cond.Id)
	}
	if cond.Filename != "" {
		q = q.Filter("filename__icontains", cond.Filename)
	}
	if cond.ArchiveId != "" {
		q = q.Filter("archive_id", cond.ArchiveId)
	}
	if cond.Path != nil {
		if cond.Path.Path != "" {
			path := &Paths{
				Path: cond.Path.Path,
			}
			paths, err := GetPaths(path, 1, 0)
			if err == nil && len(paths) != 0 {
				cond.Path.Id = paths[0].Id
			}
		}
		if cond.Path.Id != "" {
			q = q.Filter("path_id", cond.Path.Id)
		}
	}
	if cond.Host != nil {
		if cond.Host.Name != "" {
			host := &Hosts{
				Name: cond.Host.Name,
			}
			hosts, err := GetHosts(host, 1, 0)
			if err == nil && len(hosts) != 0 {
				cond.Host.Id = hosts[0].Id
			}
		}
		if cond.Host.Id != "" {
			q = q.Filter("host_id", cond.Host.Id)
		}
	}

	if cond.AppSet != nil {
		if cond.AppSet.Name != "" {
			appSet := &AppSets{
				Name: cond.AppSet.Name,
			}
			appSets, err := GetAppSets(appSet, 1, 0)
			if err == nil && len(appSets) != 0 {
				cond.AppSet.Id = appSets[0].Id
			}
		}
		if cond.AppSet.Id != "" {
			q = q.Filter("app_set_id", cond.AppSet.Id)
		}
	}

	if cond.BackupSet != nil {
		if cond.BackupSet.Name != "" {
			backupSet := &BackupSets{
				Name: cond.BackupSet.Name,
			}
			backupSets, err := GetBackupSets(backupSet, 1, 0)
			if err == nil && len(backupSets) != 0 {
				cond.AppSet.Id = backupSets[0].Id
			}
		}
		if cond.BackupSet.Id != "" {
			q = q.Filter("backup_set_id", cond.BackupSet.Id)
		}
	}
	if len(times) != 0 {
		if !times[backupTimeStart].IsZero() {
			q = q.Filter("backup_time__gte", times[backupTimeStart])
		}
		if !times[backupTimeEnd].IsZero() {
			q = q.Filter("backup_time__lte", times[backupTimeEnd])
		}
		if !times[archiveTimeStart].IsZero() {
			q = q.Filter("archived_time__gte", times[backupTimeStart])
		}
		if !times[archiveTimeEnd].IsZero() {
			q = q.Filter("archived_time__lte", times[backupTimeEnd])
		}
	}

	if limit > 0 {
		q = q.Limit(limit)
	}
	if index > 0 {
		q = q.Offset(index)
	}
	_, err := q.OrderBy("backup_time", "archived_time").
		RelatedSel(common.RelDepth).All(&r)
	if err != nil {
		return nil, err
	}
	for _, v := range r {
		if !v.BackupTime.IsZero() {
			v.BackupTime = v.BackupTime.Local()
		}
		if !v.ArchivedTime.IsZero() {
			v.ArchivedTime = v.ArchivedTime.Local()
		}

		o.LoadRelated(v, "Jobs", common.RelDepth)
	}
	return r, nil
}
