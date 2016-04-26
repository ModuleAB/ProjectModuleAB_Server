package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

const (
	OasJobTypeArchiveRetrieval = iota
	OasJobTypeInventoryRetrieval
	OasJobTypePullFromOSS
	OasJobTypePushToOSS
	OasJobTypeDeleteArchive
)

//目前只做记录用，先调用oascmd来操作OAS
type OasJobs struct {
	Id        string `orm:"pk;size(36)"`
	Vault     *Oas   `orm:"rel(fk)"`
	RequestId string
	JobId     string
	JobType   int
	Record    *Records `orm:"rek(fk);null"`
}

func init() {
	if prefix := beego.AppConfig.String("mysqlprefex"); prefix != "" {
		orm.RegisterModelWithPrefix(prefix, new(OasJobs))
	} else {
		orm.RegisterModel(new(OasJobs))
	}
}
