package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"simple_project/database"
	"simple_project/pkg/setting"
	"simple_project/routers"
	"strconv"
	"time"
)

func init() {
	//初始化配置
	setting.Setup()
	fmt.Printf("%+v", *setting.AppSetting)
}

func main() {

	//初始化数据库
	database.InitDataSources(true)
	database.PGClient.LogMode(true)
	//初始化表
	database.InitTables()
	gin.SetMode(setting.ServerSetting.RunMode)
	router := routers.InitServer()
	_ = router.Run("0.0.0.0:" + strconv.Itoa(setting.ServerSetting.HttpPort)) //启动web服务

	//优雅的关闭服务
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
