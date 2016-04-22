package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

//用户
type Users struct {
	Id       string `orm:"pk;size(36)"`
	Name     string `orm:"size(32)"`
	Password string
	Roles    []*Roles `orm:"rel(m2m);rel_table(users_to_roles)"`
}

func init() {
	if prefix := beego.AppConfig.String("mysqlprefex"); prefix != "" {
		orm.RegisterModelWithPrefix(prefix, new(Users))
	} else {
		orm.RegisterModel(new(Users))
	}
}
