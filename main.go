package main

import (
	"github.com/arong/dean/controllers"
	"github.com/arong/dean/models"
	_ "github.com/arong/dean/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/astaxie/beego/session"
)

var globalSessions *session.Manager

func init() {
	sessionConfig := &session.ManagerConfig{
		CookieName:      "dean",
		EnableSetCookie: true,
		Gclifetime:      3600,
		Maxlifetime:     3600,
		Secure:          false,
		CookieLifeTime:  2400,
		ProviderConfig:  "./tmp",
	}
	globalSessions, _ = session.NewManager("memory", sessionConfig)
	go globalSessions.GC()
}

var FilterUser = func(ctx *context.Context) {
	_, ok := ctx.Input.Session("uid").(int)
	path := ctx.Request.URL.Path
	if !ok && (path != "/api/v1/dean/user/login") && (path != "/api/v1/dean/teacher/login") {
		// return a invalid
		ctx.Output.JSON(controllers.CommResp{Code: -1, Msg: "invalid token"},
			false,
			true)
		return
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

	beego.InsertFilter("/*", beego.BeforeRouter, FilterUser)

	beego.Run(":2008")
	logs.Info("server stopped.")
}
