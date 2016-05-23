package models

import (
	"fmt"
	"moduleab_server/common"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"github.com/pborman/uuid"
)

//用户
type Users struct {
	Id        string   `orm:"pk;size(36)" json:"id" valid:"Match(/^[A-Fa-f0-9]{8}-([A-Fa-f0-9]{4}-){3}[A-Fa-f0-9]{12}$/)"`
	Name      string   `orm:"size(32);unique;index" json:"loginname" valid:"Required"`
	ShowName  string   `json:"name" valid:"Required"`
	Password  string   `valid:"Required" json:"password" valid:"Base64"`
	Roles     []*Roles `orm:"rel(m2m)" valid:"Required"`
	Removable bool     `orm:"default(1)" json:"removable"`
}

func init() {
	if prefix := beego.AppConfig.String("database::mysqlprefex"); prefix != "" {
		orm.RegisterModelWithPrefix(prefix, new(Users))
	} else {
		orm.RegisterModel(new(Users))
	}
}

func AddUser(a *Users) (string, error) {
	beego.Debug("[M] Got data:", a)
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		return "", err
	}

	a.Id = uuid.New()
	beego.Debug("[M] Got new id:", a.Id)
	a.Password = common.EncryptPassword(a.Password)
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
	_, err = o.QueryM2M(a, "Roles").Add(a.Roles)
	if err != nil {
		o.Rollback()
		return "", err
	}
	beego.Debug("[M] User info saved")
	o.Commit()
	return a.Id, nil
}

func DeleteUser(a *Users) error {
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
	_, err = o.QueryM2M(a, "Roles").Clear()
	if err != nil {
		o.Rollback()
		return err
	}
	_, err = o.QueryTable("users").Filter("removable", true).
		Filter("id", a.Id).Filter("name", a.Name).Delete()
	if err != nil {
		o.Rollback()
		return err
	}
	o.Commit()
	return nil
}

func UpdateUser(a *Users) error {
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
	if a.Roles != nil && len(a.Roles) != 0 {
		_, err = o.QueryM2M(a, "Roles").Clear()
		if err != nil {
			o.Rollback()
			return err
		}
		_, err = o.QueryM2M(a, "Roles").Add(a.Roles)
		if err != nil {
			o.Rollback()
			return err
		}
	}
	o.Commit()
	return nil
}

// If get all, just use &User{}
func GetUser(cond *Users, limit, index int) ([]*Users, error) {
	r := make([]*Users, 0)
	o := orm.NewOrm()
	q := o.QueryTable("users")
	if cond.Id != "" {
		q = q.Filter("id", cond.Id)
	}
	if cond.Name != "" {
		q = q.Filter("name", cond.Name)
	}
	if cond.ShowName != "" {
		q = q.Filter("show_name", cond.ShowName)
	}
	if cond.Password != "" {
		q = q.Filter("password", common.EncryptPassword(cond.Password))
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
		o.LoadRelated(v, "Roles", common.RelDepth)
	}
	return r, nil
}
