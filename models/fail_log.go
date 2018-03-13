package models

import (
	"fmt"
	"time"

	"github.com/ModuleAB/ModuleAB/server/common"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"github.com/pborman/uuid"
)

type FailLog struct {
	Id   string    `orm:"pk;size(36)" jsong:"id"  valid:"Match(/^[A-Fa-f0-9]{8}-([A-Fa-f0-9]{4}-){3}[A-Fa-f0-9]{12}$/)"`
	Time time.Time `orm:"type(datetime)" json:"time" valid:"Required"`
	Log  string    `json:"log" valid:"Required"`
	Host *Hosts    `orm:"rel(fk)" json:"host"`
}

func init() {
	if prefix := beego.AppConfig.String("database::mysqlprefex"); prefix != "" {
		orm.RegisterModelWithPrefix(prefix, new(FailLog))
	} else {
		orm.RegisterModel(new(FailLog))
	}
}

func AddFailLog(failLog *FailLog) (string, error) {
	beego.Debug("[M] Got data:", failLog)
	failLog.Time = time.Now()
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		return "", err
	}
	failLog.Id = uuid.New()
	validator := new(validation.Validation)
	valid, err := validator.Valid(failLog)
	if err != nil {
		o.Rollback()
		return "", err
	}
	if !valid {
		o.Rollback()
		var errS string
		for _, err := range validator.Errors {
			errS = fmt.Sprintf("%s, %s:%s", errS, err.Key, err.Message)
		}
	}
	_, err = o.Insert(failLog)
	if err != nil {
		o.Rollback()
		return "", err
	}
	o.Commit()
	return failLog.Id, nil
}

func GetFailLogs(cond *FailLog) ([]*FailLog, error) {
	r := make([]*FailLog, 0)
	o := orm.NewOrm()
	q := o.QueryTable("fail_log")
	if cond.Host != nil {
		if cond.Host.Name != "" {
			host := &Hosts{
				Name: cond.Host.Name,
			}
			hosts, err := GetHosts(host, 1, 0)
			if err == nil && len(hosts) != 0 {
				q = q.Filter("host_id", hosts[0].Id)
			}
		}
	}
	if !cond.Time.IsZero() {
		q = q.Filter("time__gte", cond.Time).Filter("time__lt", cond.Time.AddDate(0, 0, 1))
	}
	_, err := q.RelatedSel(common.RelDepth).OrderBy("-time").All(&r)
	if err != nil {
		return nil, err
	}
	return r, nil
}
