package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"simple_project/database"
	"simple_project/pkg/setting"
	"simple_project/routers"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
)

func init() {
	//初始化配置
	setting.Setup()
	fmt.Printf("%+v", *setting.AppSetting)
}

// @title Swagger Example API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8989
// @BasePath /

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @securitydefinitions.oauth2.application OAuth2Application
// @tokenUrl https://example.com/oauth/token
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.implicit OAuth2Implicit
// @authorizationurl https://example.com/oauth/authorize
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.password OAuth2Password
// @tokenUrl https://example.com/oauth/token
// @scope.read Grants read access
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl https://example.com/oauth/token
// @authorizationurl https://example.com/oauth/authorize
// @scope.admin Grants read and write access to administrative information

// @x-extension-openapi {"example": "value on a json format"}
func main() {

	//初始化数据库
	database.InitDataSources(true)
	database.PGClient.LogMode(true)
	//初始化表
	database.InitTables()
	gin.SetMode(setting.ServerSetting.RunMode)
	router := routers.InitServer()
	//_ = router.Run("0.0.0.0:" + strconv.Itoa(setting.ServerSetting.HttpPort)) //启动web服务
	go Task()
	//优雅的关闭服务
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(setting.ServerSetting.HttpPort),
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

func Task() {
	crontask := cron.New()
	crontask.AddFunc("0/2 * * * * ? ", func() {
		println("task run")
	})
	crontask.Start()
}
