// record.go
package models

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

const (
	RecordTypeBackup = iota
	RecordTypeArchive
)

type Records struct {
	Id           string      `orm:"pk;size(36)"`
	Host         *Hosts      `orm:"rel(fk)"`
	BackupSets   *BackupSets `orm:"rel(fk)"`
	AppSets      *AppSets    `orm:"rel(fk)"`
	Path         string
	Type         int       // 0 - Backup, 1 - Archive
	ArchiveId    string    `orm:"null"` // 如果Type是1（归档）时，这里应该有数据
	BackupTime   time.Time `orm:"type(datetime)"`
	ArchivedTime time.Time `orm:"type(datatime);null"`
}

func (r *Records) TableEngine() string {
	return "TokuDB"
}

func init() {
	if prefix := beego.AppConfig.String("mysqlprefex"); prefix != "" {
		orm.RegisterModelWithPrefix(prefix, new(Records))
	} else {
		orm.RegisterModel(new(Records))
	}
}
