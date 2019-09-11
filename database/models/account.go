package models

import (
	"simple_project/pkg/utils"

	"github.com/jinzhu/gorm"
)

//Account 后台用户
type Account struct {
	gorm.Model
	Account  string `gorm:"not null;unique;unique_index"` //账号
	NickName string `gorm:"not null"`                     //昵称
	Password string `gorm:"not null" json:"-"`            //密码
	Name     string `gorm:"not null"`                     //姓名
	AllowIP  string //允许登录的IP地址（IP地址以;分隔存储，例：127.0.0.1;168.134.20.140
}

//CreateTable 创建表格
func (model *Account) CreateTable(db *gorm.DB) error {
	err := db.AutoMigrate(model).Error
	//初次创建管理员账号
	if err == nil {
		db.Create(&Account{
			NickName: "管理员127",
			Account:  "admin",
			Name:     "管理员",
			Password: utils.Md5("admin")})
	}
	return err
}
