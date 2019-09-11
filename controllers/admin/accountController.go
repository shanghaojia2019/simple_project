package admin

import (
	"github.com/gin-gonic/gin"
	"simple_project/pkg/e"
	"simple_project/pkg/resp"
	"simple_project/pkg/setting"
	"simple_project/pkg/utils"
	"simple_project/services/account"
)

// Login 后台用户登录
// @Summary 登录
// @Tags 后台用户管理
// @Description  后台用户登录接口
// @Accept  multipart/form-data
// @Produce  json
// @Param username formData string true "用户名"
// @Param password formData string true "密码"
// @Success 200 {object} resp.ResponseModel
// @Router /manager/account/login/ [post]
func Login(ctx *gin.Context) {
	ret := resp.CreateResponse(ctx)
	defer ret.ResponseJSON()
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	if utils.IsEmpty(username) || utils.IsEmpty(password) {
		ret.Msg = "用户名和密码不能为空"
		return
	}
	accountModel, err := services.Login(username, password)
	if err != nil {
		ret.Msg = err.Error()
		return
	}
	token, err := utils.GenerateToken(accountModel.Account, accountModel.NickName)
	if err != nil {
		ret.Msg = "颁发token失败"
	}
	ret.Code = e.SUCCESS
	ret.Data = gin.H{"token": token, "nickname": accountModel.NickName}
}

// GetInfo 获取当前登录用户信息
// @Summary 获取当前登录用户信息
// @Tags 后台用户管理
// @Description  获取当前用户信息
// @Accept  json
// @Produce  json
// @Success 200 {object} resp.ResponseModel
// @Router /manager/account/getinfo/ [get]
func GetInfo(ctx *gin.Context) {
	ret := resp.CreateResponse(ctx)
	defer ret.ResponseJSON()
	claims, _ := utils.ParseToken(ctx.GetHeader(setting.AppSetting.AuthKey))
	accountModel, err := services.GetInfo(claims.Username)
	if err != nil {
		ret.Msg = err.Error()
	}
	ret.Code = e.SUCCESS
	ret.Data = gin.H{"account": accountModel.Account, "nickname": accountModel.NickName, "name": accountModel.Name}
}

// UpdatePassword 修改密码
// @Summary 修改当前用户密码
// @Tags 后台用户管理
// @Description  修改当前登录用户密码
// @Accept  json
// @Produce  json
// @Param oldPassword formData string true "原密码"
// @Param password formData string true "新密码"
// @Success 200 {object} resp.ResponseModel
// @Router /manager/account/updatePassword/ [post]
func UpdatePassword(ctx *gin.Context) {
	ret := resp.CreateResponse(ctx)
	defer ret.ResponseJSON()
	password := ctx.PostForm("password")
	oldPassword := ctx.PostForm("oldPassword")
	claims, _ := utils.ParseToken(ctx.GetHeader(setting.AppSetting.AuthKey))
	err := services.UpdatePassword(claims.Username, oldPassword, password)
	if err != nil {
		ret.Msg = err.Error()
		ret.Error()
		return
	}
	ret.Ok()
}
