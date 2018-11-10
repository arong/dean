package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/arong/dean/auth"
	"github.com/arong/dean/controllers"
	"github.com/arong/dean/models"
	_ "github.com/arong/dean/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/astaxie/beego/session"
	"time"
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

type baseReq struct {
	Token     string      `json:"token"`
	Timestamp int64       `json:"timestamp"`
	Check     string      `json:"check"`
	Data      interface{} `json:"data"`
}

func (b baseReq) IsValid() bool {
	md5str := fmt.Sprintf("%x", md5.Sum([]byte(b.Token+fmt.Sprintf("%d", b.Timestamp))))
	if b.Check != md5str {
		return false
	}

	if b.Timestamp+30 < time.Now().Unix() {
		return false
	}
	return true
}

var FilterUser = func(ctx *context.Context) {
	request := baseReq{}

	err := json.Unmarshal(ctx.Input.RequestBody, &request)
	if err != nil {
		logs.Debug("bad request found", ctx.Input.IP())
		ctx.Output.JSON(controllers.CommResp{Code: -3, Msg: "bad request"}, false, true)
		return
	}

	if !request.IsValid() {
		ctx.Output.JSON(controllers.CommResp{Code: -3, Msg: "invalid request"}, false, true)
		return
	}

	if !auth.VerifyToken(request.Token) {
		if ctx.Request.URL.Path != "/api/v1/dean/login/" {
			ctx.Output.JSON(controllers.CommResp{Code: -2, Msg: "invalid token"}, false, true)
			return
		}
	}

	ctx.Input.RequestBody, _ = json.Marshal(request.Data)
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
