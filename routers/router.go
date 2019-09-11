package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
	"simple_project/controllers/admin"
	_ "simple_project/docs"
	"simple_project/middleware"
	"simple_project/pkg/utils"
)

//InitServer 绑定路由，创建服务
func InitServer() *gin.Engine {
	server := utils.CreateDefaultServer(true, nil)

	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	server.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "server start...")
	})
	//管理后台
	manager := server.Group("/manager")
	manager.POST("/account/login", admin.Login)
	manager.Use(jwt.ManagerJWT())
	{
		manager.GET("/account/getInfo", admin.GetInfo)
		manager.POST("/account/updatePassword", admin.UpdatePassword)
		manager.POST("/test", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, "auth ok")
		})
	}
	//前台用户
	return server
}
