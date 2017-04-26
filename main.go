package main

import (
	_ "ruleEngine/routers"

	"github.com/astaxie/beego"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/apidoc"] = "swagger"
	}
	beego.Run()
}
