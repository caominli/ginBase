package models //定义在models包下

import (
	config "gin_scaffold/config" //导入配置包
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// 定义一个常量用来操作数据库，常量就是初始化后返回的db
var DB = Init()

// 初始化数据库连接
func Init() *gorm.DB {
	log.Print("初始化数据库连接")
	dsn := "host=" + config.Config.DBhost + " user=" + config.Config.DBuser + " password=" + config.Config.DBpassword + " dbname=" + config.Config.DBname + " port=" + config.Config.DBport + " sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("无法连接数据库")
	}
	// 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		panic("无法设置连接池")
	}
	sqlDB.SetMaxIdleConns(8)                  // 设置最大空闲连接数
	sqlDB.SetMaxOpenConns(30)                 // 设置最大打开连接数
	sqlDB.SetConnMaxIdleTime(time.Minute * 2) // 设置连接空闲时间
	sqlDB.SetConnMaxLifetime(time.Hour)       // 设置连接的最大存活时间
	// 自动创建表
	db.AutoMigrate(&Users{}, &Captchas{}, &Items{})
	log.Print("数据库连接初始化成功")
	return db
}
