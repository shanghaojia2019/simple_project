package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
	_ "simple_project/docs"
	"simple_project/middleware"
	"simple_project/pkg/utils"
	"simple_project/services/account"
)

//InitServer 绑定路由，创建服务
func InitServer() *gin.Engine {
	server := utils.CreateDefaultServer(true, nil)

	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	server.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "server start...")
	})
	//管理后台
	apiV1 := server.Group("/manager")
	apiV1.POST("/account/login", account.Login)
	apiV1.Use(jwt.AccountJWT())
	{
		apiV1.POST("/test", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, "auth ok")
		})
	}
	return server
}
