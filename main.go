package main

import (
	"flag"
	"github.com/arong/dean/models"
	_ "github.com/arong/dean/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/plugins/cors"
)

var confPath string
func main() {
	flag.StringVar(&confPath, "c", "server.conf", "set configuration `file`")
	if confPath == "" {
		logs.Warn("Please specify a config file")
		return
	}
	conf := models.DBConfig{}
	err := conf.GetConf(confPath)
	if err != nil {
		logs.Warn("[main] GetConf failed ", err)
		return
	}

	models.Init(&conf)

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

	beego.Run(":2008")
	logs.Info("server stopped.")
}
