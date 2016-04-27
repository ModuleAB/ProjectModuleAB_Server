package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

const (
	RoleFlagAdmin = iota
	RoleFlagOperator
	RoleFlagUser
)

//角色
type Roles struct {
	Id       string `orm:"pk;size(36)"`
	Name     string `orm:"size(32);unique"`
	RoleFlag int
	Users    []*Users `orm:"reverse(many)"`
}

func init() {
	if prefix := beego.AppConfig.String("database::mysqlprefex"); prefix != "" {
		orm.RegisterModelWithPrefix(prefix, new(Roles))
	} else {
		orm.RegisterModel(new(Roles))
	}
}
