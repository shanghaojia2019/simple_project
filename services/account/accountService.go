package services

import (
	"errors"
	"simple_project/database"
	"simple_project/database/models"
	"simple_project/pkg/utils"
)

// Login 后台用户登录
func Login(username string, password string) (*models.Account, error) {
	if utils.IsEmpty(username) || utils.IsEmpty(password) {
		return nil, errors.New("用户名和密码不能为空")
	}
	accountModel := &models.Account{}
	err := database.PGClient.Where("account= ? AND password= ? ", username, utils.Md5(password)).First(accountModel).Error
	if err != nil || accountModel == nil {
		return nil, errors.New("用户名或密码错误")
	}
	return accountModel, nil
}

// GetInfo 获取当前登录用户信息
func GetInfo(account string) (*models.Account, error) {
	accountModel := &models.Account{}
	err := database.PGClient.Where("account = ?", account).First(accountModel).Error
	return accountModel, err
}

// UpdatePassword 修改密码
func UpdatePassword(account string, oldpassword string, password string) error {
	accountModel := &models.Account{}
	err := database.PGClient.Where("account= ? AND password= ? ", account, utils.Md5(password)).First(accountModel).Error
	if err != nil || accountModel == nil {
		return errors.New("原密码错误")
	}
	err = database.PGClient.Model(&models.Account{}).Where("account = ?", account).Update("password", utils.Md5(password)).Error
	return err
}
