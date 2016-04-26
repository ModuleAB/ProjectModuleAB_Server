package models

import (
	"fmt"

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
	AppSet     *AppSets      `orm:"rel(fk);on_delete(set_null);null" json:"app_set"`
	BackupSets []*BackupSets `orm:"reverse(many)" json:"backup_sets"`
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
	o.Commit()
	return nil
}

// If get all, just use &Host{}
func GetHosts(cond *Hosts) ([]*Hosts, error) {
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
	_, err := q.RelatedSel().All(&r)
	if err != nil {
		return nil, err
	}
	for _, v := range r {
		o.QueryTable("backup_sets").RelatedSel().All(&v.BackupSets)
	}
	return r, nil
}
