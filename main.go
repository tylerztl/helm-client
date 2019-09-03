package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"helm-client/commons"
	_ "helm-client/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/plugins/cors"
)

func init() {
	initLogger()
}

func initLogger() {
	config := make(map[string]interface{})
	logPath := beego.AppConfig.String("logPath")
	if logPath == "" {
		logPath = commons.GetConfig().Home.Path("logs")
	}
	if fi, err := os.Stat(logPath); err != nil {
		if err := os.MkdirAll(logPath, 0755); err != nil {
			panic("Invalid log path")
		}
	} else if !fi.IsDir() {
		panic(fmt.Sprintf("%s must be a directory", logPath))
	}
	logFile := filepath.Join(logPath, "helm-client.log")
	if _, err := os.Stat(logFile); err != nil {
		if err = ioutil.WriteFile(logFile, nil, 0644); err != nil {
			panic(err)
		}
	}
	config["filename"] = logFile
	logLevel, err := beego.AppConfig.Int("logLevel")
	if nil != err {
		logLevel = beego.LevelDebug
	}
	config["level"] = logLevel

	configStr, err := json.Marshal(config)
	if nil != err {
		panic(err)
	}
	err = beego.SetLogger(logs.AdapterFile, string(configStr))
	beego.SetLogFuncCall(true)
}

func main() {
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins: true,
		AllowOrigins:    []string{"*"},
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers"},
		ExposeHeaders:   []string{"Origin", "Content-Type", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers"},
	}))
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
