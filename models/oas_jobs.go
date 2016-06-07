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
	OasJobTypeAll = iota
	OasJobTypeArchiveRetrieval
	OasJobTypeInventoryRetrieval
	OasJobTypePullFromOSS
	OasJobTypePushToOSS
	OasJobTypeDeleteArchive
)

const (
	OasJobStatusComplete   = true
	OasJobStatusIncomplete = false
)

type OasJobs struct {
	Id          string    `orm:"pk;size(36)" json:"id" valid:"Match(/^[A-Fa-f0-9]{8}-([A-Fa-f0-9]{4}-){3}[A-Fa-f0-9]{12}$/)"`
	Vault       *Oas      `orm:"rel(fk)" json:"vault" valid:"Required"`
	RequestId   string    `json:"request_id valid:"Required"`
	JobId       string    `json:"job_id" valid:"Required"`
	JobType     int       `json:"job_type" valid:"Required"`
	Status      bool      `orm:"default(0)"`
	Records     *Records  `orm:"rel(fk);null" valid:"Required"`
	CreatedTime time.Time `orm:"type(datetime)"`
}

func init() {
	if prefix := beego.AppConfig.String("database::mysqlprefex"); prefix != "" {
		orm.RegisterModelWithPrefix(prefix, new(OasJobs))
	} else {
		orm.RegisterModel(new(OasJobs))
	}
}

func AddOasJobs(a *OasJobs) (string, error) {
	beego.Debug("[M] Got data:", a)
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		return "", err
	}

	a.Id = uuid.New()
	beego.Debug("[M] Got new id:", a.Id)
	a.CreatedTime = time.Now()

	validator := new(validation.Validation)
	valid, err := validator.Valid(a)
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
	beego.Debug("[M] Got new data:", a)
	_, err = o.Insert(a)
	if err != nil {
		o.Rollback()
		return "", err
	}
	beego.Debug("[M] OasJobs info saved")
	o.Commit()
	return a.Id, nil
}

func DeleteOasJobs(a *OasJobs) error {
	beego.Debug("[M] Got data:", a)
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		return err
	}
	validator := new(validation.Validation)
	valid, err := validator.Valid(a)
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
	_, err = o.Delete(a)
	if err != nil {
		o.Rollback()
		return err
	}
	o.Commit()
	return nil
}

func UpdateOasJobs(a *OasJobs) error {
	beego.Debug("[M] Got data:", a)
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		return err
	}
	validator := new(validation.Validation)
	valid, err := validator.Valid(a)
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
	_, err = o.Update(a)
	if err != nil {
		o.Rollback()
		return err
	}
	o.Commit()
	return nil
}

// If get all, just use &OasJobs{}
func GetOasJobs(cond *OasJobs, limit, index int) ([]*OasJobs, error) {
	r := make([]*OasJobs, 0)
	o := orm.NewOrm()
	q := o.QueryTable("oas_jobs")
	if cond.Id != "" {
		q = q.Filter("id", cond.Id)
	}
	if cond.Vault != nil {
		q = q.Filter("vault_id", cond.Vault.Id)
	}
	if cond.RequestId != "" {
		q = q.Filter("request_id", cond.RequestId)
	}
	if cond.JobId != "" {
		q = q.Filter("job_id", cond.JobId)
	}
	if cond.JobType != OasJobTypeAll {
		q = q.Filter("job_type", cond.JobType)
	}
	if limit > 0 {
		q = q.Limit(limit)
	}
	if index > 0 {
		q = q.Offset(index)
	}
	_, err := q.RelatedSel(common.RelDepth).All(&r)

	if err != nil {
		return nil, err
	}
	return r, nil
}
