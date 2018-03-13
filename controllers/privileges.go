package controllers

import (
	"regexp"

	"github.com/ModuleAB/ModuleAB/server/models"

	"github.com/astaxie/beego/context"
)

type Privileges struct {
	Method   string
	Pattern  string
	RoleFlag int
}

var privileges []Privileges

func init() {
	privileges = make([]Privileges, 0)
}

func CheckPrivileges(userid string, ctx *context.Context) bool {
	if userid == "" {
		return false
	}
	users, err := models.GetUser(&models.Users{Id: userid}, 1, 0)
	if err != nil {
		return false
	}
	if len(users) == 0 {
		return false
	}
	for _, v := range users[0].Roles {
		for _, p := range privileges {
			if regexp.MustCompile(p.Pattern).MatchString(
				ctx.Input.URL(),
			) &&
				ctx.Input.Method() == p.Method &&
				v.RoleFlag <= p.RoleFlag {
				return true
			}
			// Default is Operator
			if v.RoleFlag <= models.RoleFlagOperator {
				return true
			}
		}
	}
	return false
}

func AddPrivilege(method, urlpattern string, roleflag int) {
	privileges = append(
		privileges,
		Privileges{
			Method:   method,
			Pattern:  urlpattern,
			RoleFlag: roleflag,
		},
	)
}
