package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Oss struct {
	Id         string `orm:"pk;size(36)"`
	Region     string
	BucketName string        `orm:"size(32)"`
	BackupSets []*BackupSets `orm:"reverse(many)"`
}

func init() {
	if prefix := beego.AppConfig.String("mysqlprefex"); prefix != "" {
		orm.RegisterModelWithPrefix(prefix, new(Oss))
	} else {
		orm.RegisterModel(new(Oss))
	}
}
