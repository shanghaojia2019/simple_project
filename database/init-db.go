package database

import (
	"fmt"
	"log"
	"simple_project/database/models"
	"simple_project/pkg/setting"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" //引用postgre驱动
)

var (
	//PGClient pg客户端
	PGClient *gorm.DB
)

//InitDataSources  初始化数据库
func InitDataSources(openDebug bool) {
	InitNewPG()
	//是否开启log
	PGClient.LogMode(openDebug)
}
func makeTables(m models.DBModel, db *gorm.DB) {
	if err := m.CreateTable(db); err != nil {
		panic(err)
	}
}

//InitTables 初始化数据库表格
func InitTables() {
	tx := PGClient.Begin()

	makeTables(&models.Account{}, tx)
	if err := tx.Commit().Error; err != nil {
		panic(err)
	}
}

//InitNewPG 初始化PG
func InitNewPG() {

	var err error
	if setting.DatabaseSetting.Type == "sqlite" {
		PGClient, err = gorm.Open("sqlite3", setting.DatabaseSetting.Name)
	} else {
		PGClient, err = gorm.Open(setting.DatabaseSetting.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			setting.DatabaseSetting.User,
			setting.DatabaseSetting.Password,
			setting.DatabaseSetting.Host,
			setting.DatabaseSetting.Name))
	}
	if err != nil {
		log.Fatalf("创建数据库链接出错: %v", err)
	}
	//defer PGClient.Close()

	//PGClient = db
	// if IsRemote() {
	// 	PGClient = NewPGDB(&PGOption{
	// 		DBName:     "DBName",
	// 		DBUser:     "postgres",
	// 		DBPort:     5432,
	// 		DBHost:     "127.0.0.1", //正式地址
	// 		DBPassword: "long121",
	// 	})
	// } else {
	// 	PGClient = NewPGDB(&PGOption{
	// 		DBName:     "DBName",
	// 		DBUser:     "postgres",
	// 		DBPort:     5432,
	// 		DBHost:     "127.0.0.1", //本地地址
	// 		DBPassword: "long121",
	// 	})
	// }
}
