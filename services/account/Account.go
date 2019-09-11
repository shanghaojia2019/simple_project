package account

import (
	"github.com/gin-gonic/gin"
	"simple_project/database"
	"simple_project/database/models"
	"simple_project/pkg/e"
	"simple_project/pkg/resp"
	"simple_project/pkg/utils"
)

// @Summary 登录
// @Produce  json
// @Param username path string true "username"
// @Param password path string true "password"
// @Success 200 {json} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /manager/account/login [post]
func Login(ctx *gin.Context) {
	ret := &resp.ResponseModel{Context: ctx, Code: e.ERROR}
	defer resp.ResponseJSON(ret)
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	if utils.IsEmpty(username) || utils.IsEmpty(password) {
		ret.Msg = "用户名和密码不能为空"
		return
	}
	accountModel := &models.Account{}
	err := database.PGClient.Where("account= ? AND password= ? ", username, utils.Md5(password)).First(accountModel).Error
	if err != nil || accountModel == nil {
		ret.Msg = "用户名密码错误"
		return
	}
	token, err := utils.GenerateToken(accountModel.Account, accountModel.NickName)
	if err != nil {
		ret.Msg = "颁发token失败"
	}
	ret.Code = e.SUCCESS
	ret.Data = gin.H{"token": token, "nickname": accountModel.NickName}
}
