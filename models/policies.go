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
	Name        string      `orm:"size(32);uniq" json:"name" valid:"Required"`
	Desc        string      `orm:"size(128);null" json:"description"`
	BackupSet   *BackupSets `orm:"rel(fk)" json:"backupset"`
	AppSets     []*AppSets  `orm:"rel(m2m);null" json:"appsets"` // null means all
	Hosts       []*Hosts    `orm:"rel(m2m);null" json:"hosts"`
	Paths       []*Paths    `orm:"rel(m2m);null" json:"paths"`
	Target      int         `json:"target"`
	Action      int         `json:"action"`
	TargetStart int         `orm:"default(0)" json:"starttime"` // Seconds, 0 means now
	TargetEnd   int         `orm:"default(-1)" json:"endtime"`  // Seconds, -1 means long long ago
	Step        int         `orm:"default(-1)" json:"step"`     // Seconds, 0 means reserve none, -1 means all
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
	if a.AppSets != nil && len(a.AppSets) != 0 {
		_, err = o.QueryM2M(a, "AppSets").Add(a.AppSets)
		if err != nil {
			o.Rollback()
			return "", err
		}
	}
	if a.Hosts != nil && len(a.Hosts) != 0 {
		_, err = o.QueryM2M(a, "Hosts").Add(a.Hosts)
		if err != nil {
			o.Rollback()
			return "", err
		}
	}
	if a.Paths != nil && len(a.Paths) != 0 {
		_, err = o.QueryM2M(a, "Paths").Add(a.Paths)
		if err != nil {
			o.Rollback()
			return "", err
		}
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
	_, err = o.QueryM2M(a, "AppSets").Clear()
	if err != nil {
		o.Rollback()
		return err
	}
	_, err = o.QueryM2M(a, "Hosts").Clear()
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
	if a.AppSets != nil {
		_, err = o.QueryM2M(a, "AppSets").Clear()
		if err != nil {
			o.Rollback()
			return err
		}
		_, err = o.QueryM2M(a, "AppSets").Add(a.AppSets)
		if err != nil {
			o.Rollback()
			return err
		}
	}
	if a.Hosts != nil {
		_, err = o.QueryM2M(a, "Hosts").Clear()
		if err != nil {
			o.Rollback()
			return err
		}
		_, err = o.QueryM2M(a, "Hosts").Add(a.Hosts)
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
		o.LoadRelated(v, "BackupSet", common.RelDepth)
		o.LoadRelated(v, "AppSets", common.RelDepth)
		o.LoadRelated(v, "Hosts", common.RelDepth)
		o.LoadRelated(v, "Paths", common.RelDepth)
	}
	return r, nil
}
