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

type OasJobs struct {
	Id        string `orm:"pk;size(36)"`
	Vault     *Oas   `orm:"rel(fk)"`
	RequestId string
	JobId     string
	JobType   int
	Records   *Records `orm:"rel(fk);null"`
}

func init() {
	if prefix := beego.AppConfig.String("database::mysqlprefex"); prefix != "" {
		orm.RegisterModelWithPrefix(prefix, new(OasJobs))
	} else {
		orm.RegisterModel(new(OasJobs))
	}
}
