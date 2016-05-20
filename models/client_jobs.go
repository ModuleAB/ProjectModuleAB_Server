package models

import (
	"fmt"
	"moduleab_server/common"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"github.com/pborman/uuid"
)

const (
	ClientJobsTypeAll = iota
	ClientJobsTypeDelete
)

type ClientJobs struct {
	Id           string   `orm:"pk;size(36)" json:"id" valid:"Match(/^[A-Fa-f0-9]{8}-([A-Fa-f0-9]{4}-){3}[A-Fa-f0-9]{12}$/)"`
	Period       int      `json:"period" valid:"Required;Min(10)"` // Second
	Type         int      `json:"type" valid:"Required"`
	ReservedTime int      `json:"reservedtime" valid:"Required"` // Second
	Host         []*Hosts `orm:"rel(m2m);on_delete(set_null)" json:"hosts"`
	Paths        []*Paths `orm:"rel(m2m);on_delete(set_null)" json:"paths"`
}

func init() {
	if prefix := beego.AppConfig.String("database::mysqlprefex"); prefix != "" {
		orm.RegisterModelWithPrefix(prefix, new(ClientJobs))
	} else {
		orm.RegisterModel(new(ClientJobs))
	}
}

func AddClientJob(a *ClientJobs) (string, error) {
	beego.Debug("[M] Got data:", a)
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		return "", err
	}

	a.Id = uuid.New()
	beego.Debug("[M] Got new id:", a.Id)
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
	if a.Host != nil {
		_, err = o.QueryM2M(a, "Host").Add(a.Host)
		if err != nil {
			o.Rollback()
			return "", err
		}
	}
	if a.Paths != nil {
		_, err = o.QueryM2M(a, "Paths").Add(a.Paths)
		if err != nil {
			o.Rollback()
			return "", err
		}
	}
	beego.Debug("[M] App set saved")
	o.Commit()
	return a.Id, nil
}

func DeleteClientJob(a *ClientJobs) error {
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
	_, err = o.QueryM2M(a, "Host").Clear()
	if err != nil {
		o.Rollback()
		return err
	}
	_, err = o.QueryM2M(a, "Paths").Clear()
	if err != nil {
		o.Rollback()
		return err
	}
	_, err = o.Delete(a)
	if err != nil {
		o.Rollback()
		return err
	}
	o.Commit()
	return nil
}

func UpdateClientJob(a *ClientJobs) error {
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
	if a.Host != nil {
		_, err = o.QueryM2M(a, "Host").Clear()
		if err != nil {
			o.Rollback()
			return err
		}
		_, err = o.QueryM2M(a, "Host").Add(a.Host)
		if err != nil {
			o.Rollback()
			return err
		}
	}
	if a.Paths != nil {
		_, err = o.QueryM2M(a, "Paths").Clear()
		if err != nil {
			o.Rollback()
			return err
		}
		_, err = o.QueryM2M(a, "Paths").Add(a.Paths)
		if err != nil {
			o.Rollback()
			return err
		}
	}
	o.Commit()
	return nil
}

// If get all, just use &Host{}
func GetClientJobs(cond *ClientJobs, limit, index int) ([]*ClientJobs, error) {
	r := make([]*ClientJobs, 0)
	o := orm.NewOrm()
	q := o.QueryTable("client_jobs")
	if cond.Id != "" {
		q = q.Filter("id", cond.Id)
	}
	if limit > 0 {
		q = q.Limit(limit)
	}
	if index > 0 {
		q = q.Offset(index)
	}
	_, err := q.All(&r)

	if err != nil {
		return nil, err
	}
	for _, v := range r {
		o.LoadRelated(v, "Host", common.RelDepth)
		o.LoadRelated(v, "Paths", common.RelDepth)
	}
	return r, nil
}
