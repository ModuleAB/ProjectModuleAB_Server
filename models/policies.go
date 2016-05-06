package models

import (
	//"time"
	"fmt"
	"moduleab_server/common"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"github.com/pborman/uuid"
)

const (
	PolicyActionAll = iota
	PolicyActionArchive
	PolicyActionDelete
)

const (
	PolicyTargetAll = iota
	PolicyTargetBackup
	PolicyTargetArchive
)

const (
	PolicyReserveAll  = -1
	PolicyReserveNone = 0
)

const (
	PolicyTargetTimeLongLongAgo = -1
	PolicyTargetTimeNow         = 0
)

//策略
type Policies struct {
	Id          string      `orm:"pk;size(36)" json:"id" valid:"Match(/^[A-Fa-f0-9]{8}-([A-Fa-f0-9]{4}-){3}[A-Fa-f0-9]{12}$/)"`
	Name        string      `orm:"size(32)" json:"name" valid:"Required"`
	Desc        string      `orm:"size(128);null" json:"description"`
	BackupSet   *BackupSets `orm:"rel(fk)" json:"backup_set"`
	AppSet      *AppSets    `orm:"rel(fk);null" json:"app_set"` // null means all
	Target      int         `json:"target"`
	Action      int         `json:"action"`
	TargetStart int         `orm:"default(0)" json:"start_time"` // Seconds, 0 means now
	TargetEnd   int         `orm:"default(-1)" json:"end_tile"`  // Seconds, -1 means long long ago
	Step        int         `orm:"default(-1)" json:"step"`      // Seconds, 0 means reserve none, -1 means all
}

func init() {
	if prefix := beego.AppConfig.String("database::mysqlprefex"); prefix != "" {
		orm.RegisterModelWithPrefix(prefix, new(Policies))
	} else {
		orm.RegisterModel(new(Policies))
	}
}

func AddPolicy(a *Policies) (string, error) {
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
	beego.Debug("[M] Policy info saved")
	o.Commit()
	return a.Id, nil
}

func DeletePolicy(a *Policies) error {
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

func UpdatePolicy(a *Policies) error {
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

// If get all, just use &Policies{}
func GetPolicies(cond *Policies, limit, index int) ([]*Policies, error) {
	r := make([]*Policies, 0)
	o := orm.NewOrm()
	q := o.QueryTable("policies")
	if cond.Id != "" {
		q = q.Filter("id", cond.Id)
	}
	if cond.Name != "" {
		q = q.Filter("name", cond.Name)
	}
	if cond.Target != PolicyTargetAll {
		q = q.Filter("target", cond.Target)
	}
	if cond.Action != PolicyActionAll {
		q = q.Filter("action", cond.Action)
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
		o.LoadRelated(v, "BackupSets", common.RelDepth)
		o.LoadRelated(v, "Jobs", common.RelDepth)
	}
	return r, nil
}
