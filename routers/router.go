package routers

import (
	"net/http"

	"simple_project/common/utils"

	"github.com/gin-gonic/gin"
)

//InitServer 绑定路由，创建服务
func InitServer() *gin.Engine {
	server := utils.CreateDefaultServer(true, nil)

	server.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "server start...")
	})

	return server
}
