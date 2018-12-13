package main

import (
	"encoding/json"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/arong/dean/base"
	"github.com/arong/dean/controllers"
	"github.com/arong/dean/models"
	_ "github.com/arong/dean/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/dgraph-io/badger"
)

var filterUser = func(ctx *context.Context) {
	request := base.BaseRequest{}
	var err error

	if ctx.Input.IsPost() {
		err := json.Unmarshal(ctx.Input.RequestBody, &request)
		if err != nil {
			logs.Debug("[filterUser] bad request found", ctx.Input.IP())
			goto Out
		}

		if !request.IsValid() {
			goto Out
		}
	} else if ctx.Input.IsGet() {
		v := ctx.Input.Query("token")
		if v == "" {
			logs.Info("[filterUser] token not found")
			goto Out
		}
		request.Token = v

		v = ctx.Input.Query("timestamp")
		if v == "" {
			logs.Info("[filterUser] timestamp not found")
			goto Out
		}
		request.Timestamp, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			goto Out
		}

		v = ctx.Input.Query("check")
		if v == "" {
			logs.Info("[filterUser] check sum not found")
			goto Out
		}
		request.Check = v

		if !request.IsValid() {
			goto Out
		}
	}

	// store login info to request context
	{
		loginInfo, ok := models.Ac.VerifyToken(request.Token)
		if !ok {
			if ctx.Request.URL.Path != "/api/v1/dean/auth/login" {
				goto Out
			}
		}
		ctx.Input.SetData(base.Private, loginInfo)
	}

	if ctx.Request.URL.Path == "/api/v1/dean/auth/logout" {
		return
	}
	ctx.Input.RequestBody, _ = json.Marshal(request.Data)
	return
Out:
	ctx.Output.JSON(controllers.BaseResponse{Code: -2, Msg: "invalid token"}, false, true)
}

func signalHandler(db *badger.DB) {
	c := make(chan os.Signal)
	signal.Notify(c)
	for {
		s := <-c
		if s == syscall.SIGINT || s == syscall.SIGQUIT || s == syscall.SIGTERM {
			logs.Info("[signalHandler] quit for signal", s)
			db.Close()
			os.Exit(1)
		}
	}
}

func main() {

	// read config
	conf := models.DBConfig{}
	err := conf.GetConf()
	if err != nil {
		logs.Warn("[main] GetConf failed ", err)
		return
	}

	// config log
	logs.SetLogger(logs.AdapterFile, `{"filename":"dean.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`)

	// init modules
	models.Init(&conf)

	// start local storage
	opts := badger.DefaultOptions
	opts.Dir = "./badger"
	opts.ValueDir = "./badger"
	db, err := badger.Open(opts)
	if err != nil {
		logs.Error("[main] badger start failed", err)
		return
	}
	defer db.Close()

	// register signal handler
	go signalHandler(db)

	models.Ac.SetStore(db)

	models.Ac.LoadToken()

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

	beego.InsertFilter("/*", beego.BeforeRouter, filterUser)

	// 开启平滑升级
	beego.BConfig.Listen.Graceful = true

	beego.Run("127.0.0.1:2008")
	logs.Info("server stopped.")
}
