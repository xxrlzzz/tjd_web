package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"traffic_jam_direction/models"
	"traffic_jam_direction/pkg/gredis"
	"traffic_jam_direction/pkg/grpc"
	"traffic_jam_direction/pkg/logging"
	"traffic_jam_direction/pkg/setting"
	"traffic_jam_direction/pkg/util"
	"traffic_jam_direction/routers"
)

func init() {
	setting.Setup(*flag.String("f", "conf/app.ini", "config file path"))

	if setting.AppSetting.EnableBasicServer() {
		logging.Setup(setting.AppSetting.IsDev())
		util.Setup()
	} else {
		log.Fatalln("Disabled basic server")
		return
	}
	if setting.AppSetting.EnableDatabase() {
		models.Setup()
		err := gredis.Setup()
		if err != nil {
			logging.WarnF("fail to setup redis with error %#v", err)
		}
	}
	if setting.AppSetting.EnableGrpc() {
		err := grpc.SetUp()
		if err != nil {
			logging.WarnF("fail to setup grpc with error %#v", err)
		}
	}

}

func main() {
	if setting.AppSetting.EnableDatabase() {
		defer models.CloseDB()
	}
	gin.SetMode(setting.ServerSetting.RunMode)
	gin.DisableConsoleColor()

	routersInit := routers.InitRouter()
	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	logging.InfoF("start http server listening %s", endPoint)
	_ = server.ListenAndServe()
}
