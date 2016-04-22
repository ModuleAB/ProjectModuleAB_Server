package models

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/pborman/uuid"
)

// 当Agent运行时，自动注册相关信息，如有则跳过
type Hosts struct {
	Id      string   `orm:"pk;size(36)"`
	Name    string   `orm:"index;unique;size(64)"`
	IpAddr  string   `orm:"index;unique;size(15)"`
	AppSets *AppSets `orm:"rel(fk);on_delete(set_null);null"`
}

func init() {
	if prefix := beego.AppConfig.String("mysqlprefex"); prefix != "" {
		orm.RegisterModelWithPrefix(prefix, new(Hosts))
	} else {
		orm.RegisterModel(new(Hosts))
	}
}

func AddHost(name string, ipAddr string) (string, error) {
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		return "", err
	}

	host := new(Hosts)
	host.Id = uuid.New()
	host.Name = name
	host.IpAddr = ipAddr
	_, err = o.Insert(host)
	if err != nil {
		o.Rollback()
		return "", err
	}

	o.Commit()
	return host.Id, nil

}

func DeleteHost(h *Hosts) error {
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		return err
	}

	_, err = o.Delete(h)
	if err != nil {
		o.Rollback()
		return err
	}
	o.Commit()
	return nil

}

// If get all, just use &Host{}
func GetHosts(cond *Host) ([]*Hosts, error) {
	r := make([]*Hosts, 0)
	q := orm.NewOrm().QueryTable("hosts")
	if cond.Id != "" {
		q = q.Filter("id", cond.Id)
	}
	if cond.Name != "" {
		q = q.Filter("name", cond.Name)
	}
	if cond.IpAddr != "" {
		q = q.Filter("ip_addr", cond.IpAddr)
	}
	_, err = q.All(&r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

/*
func GetAllHosts() ([]*Hosts, error) {
	var s interface{}
	o := orm.NewOrm()
	err := o.QueryTable("hosts").All(&s)
	if err != nil {
		return nil, err
	}
	h, ok := s.([]*Hosts)
	if !ok {
		return nil, fmt.Errorf("Bad data of Hosts")
	}
	return h, nil
}
*/
