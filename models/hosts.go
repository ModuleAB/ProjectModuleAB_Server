package models

import (
	"fmt"
	"moduleab_server/common"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"github.com/pborman/uuid"
)

// 当Agent运行时，自动注册相关信息，如有则跳过
type Hosts struct {
	Id         string        `orm:"pk;size(36)" json:"id" valid:"Match(/^[A-Fa-f0-9]{8}-([A-Fa-f0-9]{4}-){3}[A-Fa-f0-9]{12}$/)"`
	Name       string        `orm:"index;unique;size(64)" json:"name" valid:"Required"`
	IpAddr     string        `orm:"index;unique;size(15)" json:"ip" valid:"Required;IP"`
	AppSet     *AppSets      `orm:"rel(fk);on_delete(set_null);null" json:"appset"`
	BackupSets []*BackupSets `orm:"rel(m2m);on_delete(set_null)" json:"backupsets"`
	Paths      []*Paths      `orm:"rel(m2m);on_delete(set_null)" json:"path"`
	ClientJobs []*ClientJobs `orm:"rel(m2m);on_delete(set_null)" json:"jobs"`
}

func init() {
	if prefix := beego.AppConfig.String("database::mysqlprefex"); prefix != "" {
		orm.RegisterModelWithPrefix(prefix, new(Hosts))
	} else {
		orm.RegisterModel(new(Hosts))
	}
}

func AddHost(host *Hosts) (string, error) {
	beego.Debug("[M] Got data:", host)
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		return "", err
	}

	host.Id = uuid.New()
	beego.Debug("[M] Got id:", host.Id)
	validator := new(validation.Validation)
	valid, err := validator.Valid(host)
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

	beego.Debug("[M] Got new data:", host)
	_, err = o.Insert(host)
	if err != nil {
		o.Rollback()
		return "", err
	}
	if host.Paths != nil {
		err = AddHostPaths(host, host.Paths)
		if err != nil {
			o.Rollback()
			return "", err
		}
	}
	if host.ClientJobs != nil {
		err = AddHostClientJobs(host, host.ClientJobs)
		if err != nil {
			o.Rollback()
			return "", err
		}
	}
	if host.BackupSets != nil {
		err = AddHostBackupSets(host, host.BackupSets)
		if err != nil {
			o.Rollback()
			return "", err
		}
	}
	beego.Debug("[M] Host data saved")
	o.Commit()
	return host.Id, nil

}

func DeleteHost(h *Hosts) error {
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
	err = ClearHostPaths(h)
	if err != nil {
		o.Rollback()
		return err
	}
	err = ClearHostClientJobs(h)
	if err != nil {
		o.Rollback()
		return err
	}
	err = ClearHostBackupSets(h)
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

func UpdateHost(h *Hosts) error {
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
	if h.Paths != nil {
		err = ClearHostPaths(h)
		if err != nil {
			o.Rollback()
			return err
		}
		err = AddHostPaths(h, h.Paths)
		if err != nil {
			o.Rollback()
			return err
		}
	}
	if h.BackupSets != nil {
		err = ClearHostBackupSets(h)
		if err != nil {
			o.Rollback()
			return err
		}
		err = AddHostBackupSets(h, h.BackupSets)
		if err != nil {
			o.Rollback()
			return err
		}
	}
	o.Commit()
	return nil
}

// If get all, just use &Host{}
func GetHosts(cond *Hosts, limit, index int) ([]*Hosts, error) {
	r := make([]*Hosts, 0)
	o := orm.NewOrm()
	q := o.QueryTable("hosts")
	if cond.Id != "" {
		q = q.Filter("id", cond.Id)
	}
	if cond.Name != "" {
		q = q.Filter("name", cond.Name)
	}
	if cond.IpAddr != "" {
		q = q.Filter("ip_addr", cond.IpAddr)
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
		o.LoadRelated(v, "BackupSets", common.RelDepth)
		o.LoadRelated(v, "Paths", common.RelDepth+5)
		o.LoadRelated(v, "ClientJobs", common.RelDepth)
	}
	return r, nil
}

func AddHostPaths(host *Hosts, paths []*Paths) error {
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		return err
	}

	if paths != nil {
		_, err = o.QueryM2M(host, "Paths").Add(paths)
		if err != nil {
			o.Rollback()
			return err
		}
	}
	o.Commit()
	return nil
}

func DeleteHostPaths(host *Hosts, paths []*Paths) error {
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		return err
	}

	if paths != nil {
		_, err = o.QueryM2M(host, "Paths").Remove(paths)
		if err != nil {
			o.Rollback()
			return err
		}
	}
	o.Commit()
	return nil
}

func ClearHostPaths(host *Hosts) error {
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		o.Rollback()
		return err
	}
	_, err = o.QueryM2M(host, "Paths").Clear()
	if err != nil {
		o.Rollback()
		return err
	}
	o.Commit()
	return nil
}

func AddHostClientJobs(host *Hosts, jobs []*ClientJobs) error {
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		return err
	}

	if jobs != nil {
		_, err = o.QueryM2M(host, "ClientJobs").Add(jobs)
		if err != nil {
			o.Rollback()
			return err
		}
	}
	o.Commit()
	return nil
}

func DeleteHostClientJobs(host *Hosts, jobs []*ClientJobs) error {
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		return err
	}

	if jobs != nil {
		_, err = o.QueryM2M(host, "ClientJobs").Remove(jobs)
		if err != nil {
			o.Rollback()
			return err
		}
	}
	o.Commit()
	return nil
}

func ClearHostClientJobs(host *Hosts) error {
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		o.Rollback()
		return err
	}
	_, err = o.QueryM2M(host, "ClientJobs").Clear()
	if err != nil {
		o.Rollback()
		return err
	}
	o.Commit()
	return nil
}

func AddHostBackupSets(host *Hosts, backupSets []*BackupSets) error {
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		o.Rollback()
		return err
	}

	if backupSets != nil {
		_, err = o.QueryM2M(host, "BackupSets").Add(backupSets)
		if err != nil {
			o.Rollback()
			return err
		}
	}
	o.Commit()
	return nil
}

func ClearHostBackupSets(host *Hosts) error {
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		o.Rollback()
		return err
	}
	_, err = o.QueryM2M(host, "BackupSets").Clear()
	if err != nil {
		o.Rollback()
		return err
	}
	o.Commit()
	return nil
}
