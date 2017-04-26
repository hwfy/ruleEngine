package controllers

import (
	"ruleEngine/models"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/hwfy/expression"
)

// Rule resolution controller
type AnalyController struct {
	beego.Controller
}

// @Title Analy
// @Description Find match rules based on business data
// @Param	Authorization header string      true "token for login api"
// @Param   body  		  body   models.data true "the rule name and business data"
// @Success 200 {string} numerical
// @Success 201 {string} json
// @Success 203 the expression is incorrect
// @Failure 401 the token is invalid
// @Failure 404 rules do not match
// @router / [options]
// @router / [post]
func (o *AnalyController) Post() {
	result, err := models.GetResult(o.Ctx.Input.RequestBody)
	if err == nil {
		number, err := strconv.ParseFloat(result, 64)
		if err == nil {
			o.Ctx.Output.SetStatus(200)
			o.Data["json"] = number
		} else {
			tokens, err := expression.Parse(result)
			if err == nil {
				o.Ctx.Output.SetStatus(201)
				o.Data["json"] = tokens
			} else {
				o.Ctx.Output.SetStatus(203)
				o.Data["json"] = err
			}
		}
	} else {
		o.Ctx.Output.SetStatus(404)
		o.Data["json"] = err.Error()
	}
	o.ServeJSON()
}
