package main

import (
	"github.com/arong/dean/models"
	_ "github.com/arong/dean/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/plugins/cors"
)

var FilterUser = func(ctx *context.Context) {
	uid := ctx.Input.Session("uid")
	if uid == nil {
		ctx.Redirect(302, "/login")
		return
	}

	_, ok := uid.(int)
	if !ok && ctx.Request.RequestURI != "/login" {
		ctx.Redirect(302, "/login")
	}
}

func main() {
	conf := models.DBConfig{}
	err := conf.GetConf()
	if err != nil {
		logs.Warn("[main] GetConf failed ", err)
		return
	}

	models.Init(&conf)

	logs.SetLogger(logs.AdapterFile, `{"filename":"dean.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`)

	logs.Info("server start...")
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = false
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))

	//beego.InsertFilter("/*", beego.BeforeRouter, FilterUser)

	beego.Run(":2008")
	logs.Info("server stopped.")
}
