package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

//应用集
type AppSets struct {
	Id       string      `orm:"pk;size(36)"`
	Name     string      `orm:"size(32)"`
	Desc     string      `orm:"size(128);null"`
	Policies []*Policies `orm:"reverse(many)"`
	Hosts    []*Hosts    `orm:"reverse(many)"`
}

func init() {
	if prefix := beego.AppConfig.String("mysqlprefex"); prefix != "" {
		orm.RegisterModelWithPrefix(prefix, new(AppSets))
	} else {
		orm.RegisterModel(new(AppSets))
	}
}
