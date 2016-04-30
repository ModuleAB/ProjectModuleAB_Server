package models

import (
	"fmt"
	"moduleab_server/common"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"github.com/pborman/uuid"
)

type Paths struct {
	Id         string        `orm:"pk;size(36)" json:"id" valid:"Match(/^[A-Fa-f0-9]{8}-([A-Fa-f0-9]{4}-){3}[A-Fa-f0-9]{12}$/)"`
	Path       string        `orm:"unique;index" json:"path" valid:"Required"`
	Host       []*Hosts      `orm:"reverse(many)" json:"host"`
	AppSet     []*AppSets    `orm:"rel(m2m)" json:"app_set"`
	BackupSet  *BackupSets   `orm:"rel(fk)" json:"backup_set"`
	ClientJobs []*ClientJobs `orm:"rel(m2m)" json:"jobs"`
	Records    []*Records    `orm:"reverse(many)" json:"records"`
}

func init() {
	if prefix := beego.AppConfig.String("database::mysqlprefex"); prefix != "" {
		orm.RegisterModelWithPrefix(prefix, new(Paths))
	} else {
		orm.RegisterModel(new(Paths))
	}
}

func AddPath(path *Paths) (string, error) {
	beego.Debug("[M] Got data:", path)
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		return "", err
	}

	path.Id = uuid.New()
	beego.Debug("[M] Got id:", path.Id)
	validator := new(validation.Validation)
	valid, err := validator.Valid(path)
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

	beego.Debug("[M] Got new data:", path)
	_, err = o.Insert(path)
	if err != nil {
		o.Rollback()
		return "", err
	}
	err = AddPathsAppSets(path, path.AppSet)
	if err != nil {
		o.Rollback()
		return "", err
	}
	err = AddPathsClientJobs(path, path.ClientJobs)
	if err != nil {
		o.Rollback()
		return "", err
	}
	beego.Debug("[M] Path data saved")
	o.Commit()
	return path.Id, nil

}

func DeletePath(h *Paths) error {
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
	err = ClearPathsAppSets(h)
	if err != nil {
		o.Rollback()
		return err
	}
	err = ClearPathsClientJobs(h)
	if err != nil {
		o.Rollback()
		return err
	}
	_, err = o.Delete(h)
	if err != nil {
		o.Rollback()
		return err
	}
	o.Commit()
	return nil
}

func UpdatePath(h *Paths) error {
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

// If get all, just use &Path{}
func GetPaths(cond *Paths, limit, index int) ([]*Paths, error) {

	r := make([]*Paths, 0)
	o := orm.NewOrm()
	q := o.QueryTable("paths")
	if cond.Id != "" {
		q = q.Filter("id", cond.Id)
	}
	if cond.Path != "" {
		q = q.Filter("path", cond.Path)
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
		o.LoadRelated(v, "Host", common.RelDepth)
		o.LoadRelated(v, "AppSet", common.RelDepth)
		o.LoadRelated(v, "ClientJobs", common.RelDepth)
		o.LoadRelated(v, "Records", common.RelDepth)
	}
	return r, nil
}

func AddPathsAppSets(path *Paths, appSets []*AppSets) error {
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		o.Rollback()
		return err
	}

	if appSets != nil {
		_, err = o.QueryM2M(path, "AppSets").Add(appSets)
		if err != nil {
			o.Rollback()
			return err
		}
	}
	o.Commit()
	return nil
}

func DeletePathsAppSets(path *Paths, appSets []*AppSets) error {
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		o.Rollback()
		return err
	}

	if appSets != nil {
		_, err = o.QueryM2M(path, "AppSets").Remove(appSets)
		if err != nil {
			o.Rollback()
			return err
		}
	}
	o.Commit()
	return nil
}

func ClearPathsAppSets(path *Paths) error {
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		o.Rollback()
		return err
	}
	_, err = o.QueryM2M(path, "AppSets").Clear()
	if err != nil {
		o.Rollback()
		return err
	}
	o.Commit()
	return nil
}

func AddPathsClientJobs(path *Paths, clientJobs []*ClientJobs) error {
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		o.Rollback()
		return err
	}

	if clientJobs != nil {
		_, err = o.QueryM2M(path, "ClientJobs").Add(clientJobs)
		if err != nil {
			o.Rollback()
			return err
		}
	}
	o.Commit()
	return nil
}

func DeletePathsClientJobs(path *Paths, clientJobs []*ClientJobs) error {
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		o.Rollback()
		return err
	}
	if clientJobs != nil {
		_, err = o.QueryM2M(path, "ClientJobs").Remove(clientJobs)
		if err != nil {
			o.Rollback()
			return err
		}
	}
	o.Commit()
	return nil
}

func ClearPathsClientJobs(path *Paths) error {
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		o.Rollback()
		return err
	}
	_, err = o.QueryM2M(path, "ClientJobs").Clear()
	if err != nil {
		o.Rollback()
		return err
	}
	o.Commit()
	return nil
}
