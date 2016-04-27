package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

//用户
type Users struct {
	Id       string `orm:"pk;size(36)"`
	Name     string `orm:"size(32);unique;index"`
	Password string
	Roles    []*Roles `orm:"rel(m2m)"`
}

func init() {
	if prefix := beego.AppConfig.String("database::mysqlprefex"); prefix != "" {
		orm.RegisterModelWithPrefix(prefix, new(Users))
	} else {
		orm.RegisterModel(new(Users))
	}
}
