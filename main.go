package main

import (
	"encoding/json"
	"github.com/arong/dean/base"
	"github.com/arong/dean/controllers"
	"github.com/arong/dean/models"
	_ "github.com/arong/dean/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/plugins/cors"
	"strconv"
)

//var globalSessions *session.Manager

var FilterUser = func(ctx *context.Context) {
	request := base.BaseRequest{}
	var err error

	if ctx.Input.IsPost() {
		err := json.Unmarshal(ctx.Input.RequestBody, &request)
		if err != nil {
			logs.Debug("bad request found", ctx.Input.IP())
			goto Out
		}

		if !request.IsValid() {
			goto Out
		}
	} else if ctx.Input.IsGet() {
		v := ctx.Input.Query("token")
		if v == "" {
			logs.Info("[] token not found")
			goto Out
		}
		request.Token = v

		v = ctx.Input.Query("timestamp")
		if v == "" {
			logs.Info("[] timestamp not found")
			goto Out
		}
		request.Timestamp, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			goto Out
		}

		v = ctx.Input.Query("check")
		if v == "" {
			logs.Info("[] check sum not found")
			goto Out
		}
		request.Check = v

		if !request.IsValid() {
			goto Out
		}
	}

	if !models.Ac.VerifyToken(request.Token) {
		if ctx.Request.URL.Path != "/api/v1/dean/auth/login" {
			goto Out
		}
	}

	if ctx.Request.URL.Path == "/api/v1/dean/auth/logout" {
		return
	}
	ctx.Input.RequestBody, _ = json.Marshal(request.Data)
	return
Out:
	ctx.Output.JSON(controllers.BaseResponse{Code: -2, Msg: "invalid token"}, false, true)
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

	beego.Run("127.0.0.1:2008")
	logs.Info("server stopped.")
}
