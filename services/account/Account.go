package services

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"simple_project/database"
	"simple_project/database/models"
	"simple_project/pkg/e"
	"simple_project/pkg/resp"
)

func Login(ctx *gin.Context) {
	ret:= &resp.ResponseModel{}
	accountModel :=&models.Account{}
	username:=ctx.PostForm("username")
	password:=ctx.PostForm("password")
	err :=database.PGClient.Where("account= ? AND password= ? ", username, password).First(accountModel).Error
	if err!=nil||accountModel==nil {
		ret.Code=e.ERROR
		ret.Msg="用户名密码错误"
	}
	ret.Data=accountModel
	ret.Context.JSON(http.StatusOK,accountModel)
}
