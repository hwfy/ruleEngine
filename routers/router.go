// @APIVersion 1.0.0
// @Title Rule engine API
// @Description Rule management, binding, parsing
// @Contact yourEmail@gmail.com
package routers

import (
	"ruleEngine/controllers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	ns := beego.NewNamespace("/v1/rule",
		beego.NSNamespace("/mamt",
			beego.NSInclude(
				&controllers.RuleController{},
			),
		),
		beego.NSNamespace("/bind",
			beego.NSInclude(
				&controllers.BindController{},
			),
		),
		beego.NSNamespace("/analy",
			beego.NSInclude(
				&controllers.AnalyController{},
			),
		),
	)
	beego.AddNamespace(ns.Filter("before", func(ctx *context.Context) {
		ctx.Output.Header("Access-Control-Allow-Origin", "*")
		ctx.Output.Header("Access-Control-Allow-Headers", "Authorization,Content-Type")
		ctx.Output.Header("Access-Control-Allow-Methods", "DELETE,PUT")

		//here you can add authentication to the token
	}))
}
