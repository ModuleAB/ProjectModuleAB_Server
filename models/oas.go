package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

//目前只做记录用，先调用oascmd来操作OAS
type Oas struct {
	Id         string `orm:"pk;size(36)"`
	Region     string
	VaultName  string        `orm:"size(32)"`
	VaultId    string        `orm:"size(32)"`
	BackupSets []*BackupSets `orm:"reverse(many)"`
	Jobs       []*OasJobs    `orm:"reverse(many)"`
}

func init() {
	if prefix := beego.AppConfig.String("mysqlprefex"); prefix != "" {
		orm.RegisterModelWithPrefix(prefix, new(Oas))
	} else {
		orm.RegisterModel(new(Oas))
	}
}
