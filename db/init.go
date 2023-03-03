package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"k8smanager-demo/config"
)

var (
	isInit bool
	GORM   *gorm.DB
	err    error
)

// Init db的初始化函数，与数据库建立链接
func Init() {
	//判断是否以及初始化了
	if isInit {
		return
	}

	// 组装连接配置
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&",
		config.DbUser,
		config.DbPass,
		config.DbHost,
		config.DbPort,
		config.DbName,
	)
	GORM, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("数据库链接失败")
	}
	// 设置gorm日志输出模式
	GORM.Logger.LogMode(logger.Info)
	//
	db, err := GORM.DB()
	if err != nil {
		panic("获取数据库对象错误")
	}
	// 连接池相关设置
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetConnMaxLifetime(config.MaxLifeTime)
}

// Close 关闭数据库廉价而
func Close() error {
	db, err := GORM.DB()
	if err != nil {
		panic("获取数据库对象错误")
	}
	return db.Close()
}
