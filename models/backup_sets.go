package models

import (
	"fmt"
	"moduleab_server/common"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"github.com/pborman/uuid"
)

//备份集
type BackupSets struct {
	Id       string      `orm:"pk;size(36)" json:"id" valid:"Match(/^[A-Fa-f0-9]{8}-([A-Fa-f0-9]{4}-){3}[A-Fa-f0-9]{12}$/)"`
	Name     string      `orm:"size(32);unique;index" json:"name" valid:"Required"`
	Desc     string      `orm:"size(128);null" json:"description"`
	Oss      *Oss        `orm:"null;rel(fk);on_delete(set_null)" json:"oss"`
	Oas      *Oas        `orm:"null;rel(fk);on_delete(set_null)" json:"oas"`
	Policies []*Policies `orm:"reverse(many)" json:"policies"`
	Paths    []*Paths    `orm:"reverse(many)" json:"paths"`
}

func init() {
	if prefix := beego.AppConfig.String("database::mysqlprefex"); prefix != "" {
		orm.RegisterModelWithPrefix(prefix, new(BackupSets))
	} else {
		orm.RegisterModel(new(BackupSets))
	}
}

func AddBackupSet(a *BackupSets) (string, error) {
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
	beego.Debug("[M] App set saved")
	o.Commit()
	return a.Id, nil
}

func DeleteBackupSet(a *BackupSets) error {
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

func UpdateBackupSet(a *BackupSets) error {
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

// If get all, just use &Host{}
func GetBackupSets(cond *BackupSets, limit, index int) ([]*BackupSets, error) {
	r := make([]*BackupSets, 0)
	o := orm.NewOrm()
	q := o.QueryTable("backup_sets")
	if cond.Id != "" {
		q = q.Filter("id", cond.Id)
	}
	if cond.Name != "" {
		q = q.Filter("name", cond.Name)
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
	for _, v := range r {
		o.LoadRelated(v, "Policies", common.RelDepth)
		o.LoadRelated(v, "Paths", common.RelDepth)
	}
	return r, nil
}
