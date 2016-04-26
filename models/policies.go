package models

import (
	//"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

const (
	PolicyActionArchive = iota
	PolicyActionDelete
)

const (
	PolicyTargetBackup = iota
	PolicyTargetArchive
)

const (
	PolicyReserveAll  = -1
	PolicyReserveNone = 0
)

//策略
type Policies struct {
	Id            string      `orm:"pk;size(36)"`
	Name          string      `orm:"size(32)"`
	Desc          string      `orm:"size(128);null"`
	BackupSet     *BackupSets `orm:"rel(fk)"`
	AppSet        *AppSets    `orm:"rel(fk);null"` // null means all
	Target        int
	Action        int
	TargetStart   int `orm:"default(0)"`  // Seconds, 0 means now
	TargetEnd     int `orm:"default(-1)"` // Seconds, -1 means long long ago
	ReservePeriod int `orm:"default(-1)"` // Seconds, 0 means reserve none, -1 means all
}

func init() {
	if prefix := beego.AppConfig.String("database::mysqlprefex"); prefix != "" {
		orm.RegisterModelWithPrefix(prefix, new(Policies))
	} else {
		orm.RegisterModel(new(Policies))
	}
}
