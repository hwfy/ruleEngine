package controllers

import (
	"encoding/json"
	"ruleEngine/models"

	"github.com/astaxie/beego"
)

// Rule binding controller
type BindController struct {
	beego.Controller
}

// @Title Bind
// @Description Bind multiple rules based on form name
// @Param	Authorization	header string        true "token for login api"
// @Param	name			path   string        true "the form name"
// @Param	body			body   models.Form   true "the rule content and form field"
// @Success 200 {string} OK
// @Failure 400 request parameter is incorrect
// @Failure 401 the token is invalid
// @Failure 500 database failure
// @router / [options]
// @router /:name [post]
func (r *BindController) Post() {
	name := r.Ctx.Input.Param(":name")
	if name != "{name}" {
		var rules map[string]map[string]models.Field
		if err := json.Unmarshal(r.Ctx.Input.RequestBody, &rules); err == nil {
			err = models.Binding(name, rules)
			if err == nil {
				r.Data["json"] = "OK"
			} else {
				r.Ctx.Output.SetStatus(500)
				r.Data["json"] = err.Error()
			}
		} else {
			r.Ctx.Output.SetStatus(400)
			r.Data["json"] = err.Error()
		}
	} else {
		r.Ctx.Output.SetStatus(400)
		r.Data["json"] = "The form name is empty"
	}
	r.ServeJSON()
}

// @Title GetAll
// @Description Obtain all bound rules based on the form name
// @Param	Authorization	header string true	"token for login api"
// @Param	name			path   string true	"the form name"
// @Success 200 {object} models.Form
// @Failure 400 the form name is empty
// @Failure 401 the token is invalid
// @Failure 404 the rule does not exist
// @router /:name [get]
func (r *BindController) GetAll() {
	name := r.Ctx.Input.Param(":name")
	if name != "{name}" {
		rules, err := models.LookUp(name)
		if err == nil {
			r.Data["json"] = rules
		} else {
			r.Ctx.Output.SetStatus(404)
			r.Data["json"] = err.Error()
		}
	} else {
		r.Ctx.Output.SetStatus(400)
		r.Data["json"] = "The form name is empty"
	}
	r.ServeJSON()
}

// @Title Delete
// @Description Delete the bound fields based on the form and rule names
// @Param	Authorization	header string true "token for login api"
// @Param	name			query  string true "the form name"
// @Param   rule		    query  string true "the rule name"
// @Success 200 {string} OK
// @Failure 400 the request parameter is empty
// @Failure 401 the token is invalid
// @Failure 404 the rule does not exist
// @router / [delete]
func (r *BindController) Delete() {
	name := r.GetString("name")
	rule := r.GetString("rule")
	if name != "" && rule != "" {
		err := models.UnBind(name, rule)
		if err == nil {
			r.Data["json"] = "OK"
		} else {
			r.Ctx.Output.SetStatus(404)
			r.Data["json"] = err.Error()
		}
	} else {
		r.Ctx.Output.SetStatus(400)
		r.Data["json"] = "The form or rule name is empty"
	}
	r.ServeJSON()
}
