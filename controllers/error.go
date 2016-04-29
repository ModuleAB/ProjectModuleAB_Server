package controllers

import (
	"net/http"

	"github.com/astaxie/beego"
)

type ErrorController struct {
	beego.Controller
}

func (e *ErrorController) Error404() {
	e.Ctx.Output.SetStatus(http.StatusNotFound)
	e.Data["json"] = map[string]string{
		"message": "API not implemented.",
	}
	e.ServeJSON()
}

func (e *ErrorController) Error500() {
	e.Ctx.Output.SetStatus(http.StatusInternalServerError)
	e.Data["json"] = map[string]string{
		"message": "Got serieve error.",
	}
	e.ServeJSON()
}
