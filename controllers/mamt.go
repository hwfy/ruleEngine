package controllers

import (
	"encoding/json"
	"ruleEngine/models"

	"github.com/astaxie/beego"
)

// Rule management controller
type RuleController struct {
	beego.Controller
}

// @Title Create
// @Description Add or modify multiple rule sets
// @Param	Authorization	header	string			true "token for login api"
// @Param	body			body	models.Rules	true "the rules content"
// @Success 200 {string} OK
// @Failure 400 request parameter is incorrect
// @Failure 401 the token is invalid
// @Failure 500 database failure
// @router / [options]
// @router / [post]
func (r *RuleController) PostAll() {
	var rules []models.Rules
	if err := json.Unmarshal(r.Ctx.Input.RequestBody, &rules); err == nil {
		if len(rules) != 0 {
			err = models.SetAll(rules)
			if err == nil {
				r.Data["json"] = "OK"
			} else {
				r.Ctx.Output.SetStatus(500)
				r.Data["json"] = err.Error()
			}
		} else {
			r.Ctx.Output.SetStatus(400)
			r.Data["json"] = "Parameter is empty"
		}
	} else {
		r.Ctx.Output.SetStatus(400)
		r.Data["json"] = err.Error()
	}
	r.ServeJSON()
}

// @Title GetAll
// @Description Get all rule names
// @Param	Authorization header string true "token for login api"
// @Success 200 {object} []string
// @Failure 401 the token is invalid
// @Failure 404 the rule does not exist
// @router / [get]
func (r *RuleController) GetKeys() {
	names, err := models.GetNames()
	if err == nil {
		r.Data["json"] = names
	} else {
		r.Ctx.Output.SetStatus(404)
		r.Data["json"] = err.Error()
	}
	r.ServeJSON()
}

// @Title GetAll
// @Description Get all rules according to the rule name, and the value field in the sample data is not shown
// @Param	Authorization	header string true	"token for login api"
// @Param	name			path   string true	"the rule name"
// @Success 200 {object} models.Result
// @Failure 400 the rule name is empty
// @Failure 401 the token is invalid
// @Failure 404 the rule does not exist
// @router /:name [get]
func (r *RuleController) GetAll() {
	name := r.Ctx.Input.Param(":name")
	if name != "{name}" {
		rules, err := models.GetAll(name)
		if err == nil {
			r.Data["json"] = rules
		} else {
			r.Ctx.Output.SetStatus(404)
			r.Data["json"] = err.Error()
		}
	} else {
		r.Ctx.Output.SetStatus(400)
		r.Data["json"] = "The rule name is empty"
	}
	r.ServeJSON()
}

// @Title DeleteAll
// @Description Delete multiple rules in different rule sets
// @Param	Authorization	header string				true	"token for login api"
// @Param	body			body   map[string][]string	true	"{"ruleName":["id1", "id2" ...] }"
// @Success 200 {string} OK
// @Failure 400 the request parameter is empty
// @Failure 401 the token is invalid
// @Failure 404 the rule does not exist
// @router / [delete]
func (r *RuleController) DeleteAll() {
	var rules map[string][]string
	if err := json.Unmarshal(r.Ctx.Input.RequestBody, &rules); err == nil {
		err := models.DeleteAll(rules)
		if err == nil {
			r.Data["json"] = "OK"
		} else {
			r.Ctx.Output.SetStatus(404)
			r.Data["json"] = err.Error()
		}
	} else {
		r.Ctx.Output.SetStatus(400)
		r.Data["json"] = err.Error()
	}
	r.ServeJSON()
}

// @Title Get
// @Description Get rules according to the rule name and id, and the value field in the sample data is not shown
// @Param	Authorization 	header string	true	"token for login api"
// @Param	name			path   string	true	"the rule name"
// @Param	id			  	path   string	true	"the rule id"
// @Success 200 {object} models.Rule
// @Failure 400 the request parameter is empty
// @Failure 401 the token is invalid
// @Failure 404 the rule does not exist
// @router /:name/:id [get]
func (r *RuleController) Get() {
	name := r.Ctx.Input.Param(":name")
	id := r.Ctx.Input.Param(":id")
	if name != "{name}" && id != "{id}" {
		rule, err := models.GetOne(name, id)
		if err == nil {
			r.Data["json"] = rule
		} else {
			r.Ctx.Output.SetStatus(404)
			r.Data["json"] = err.Error()
		}
	} else {
		r.Ctx.Output.SetStatus(400)
		r.Data["json"] = "The rule name or ID is empty"
	}
	r.ServeJSON()
}

// @Title Update
// @Description According to the rule name and id update, conditions need to add value field
// @Param	Authorization	header	string 		true	"token for login api"
// @Param	name			path  	string 		true	"the rule name"
// @Param	id			  	path  	string		true	"the rule id"
// @Param	body			body 	models.Rule	true	"the rule content"
// @Success 200 {string} OK
// @Failure 400 request parameter is incorrect
// @Failure 401 the token is invalid
// @Failure 404 the rule does not exist
// @router /:name/:id [put]
func (r *RuleController) Put() {
	name := r.Ctx.Input.Param(":name")
	id := r.Ctx.Input.Param(":id")
	if name != "{name}" && id != "{id}" {
		var rule models.Rule
		if err := json.Unmarshal(r.Ctx.Input.RequestBody, &rule); err == nil {
			err = models.Update(name, id, rule)
			if err == nil {
				r.Data["json"] = "OK"
			} else {
				r.Ctx.Output.SetStatus(404)
				r.Data["json"] = err.Error()
			}
		} else {
			r.Ctx.Output.SetStatus(400)
			r.Data["json"] = err.Error()
		}
	} else {
		r.Ctx.Output.SetStatus(400)
		r.Data["json"] = "The rule name or ID is empty"
	}
	r.ServeJSON()
}

// @Title Delete
// @Description Delete the rule according to rule name and id
// @Param	Authorization	header	string true	"token for login api"
// @Param	name			path 	string true	"the rule name"
// @Param	id			  	path	string true	"the rule id"
// @Success 200 {string} OK
// @Failure 400 the request parameter is empty
// @Failure 401 the token is invalid
// @Failure 404 the rule does not exist
// @router /:name/:id [delete]
func (r *RuleController) Delete() {
	name := r.Ctx.Input.Param(":name")
	id := r.Ctx.Input.Param(":id")
	if name != "{name}" && id != "{id}" {
		err := models.Delete(name, id)
		if err == nil {
			r.Data["json"] = "OK"
		} else {
			r.Ctx.Output.SetStatus(404)
			r.Data["json"] = err.Error()
		}
	} else {
		r.Ctx.Output.SetStatus(400)
		r.Data["json"] = "The rule name or ID is empty"
	}
	r.ServeJSON()
}
