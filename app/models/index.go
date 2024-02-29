package models

import (
	Config "EasyGin/app/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// ConnectDB 创建连接
func ConnectDB(database Config.Database) error {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		database.Host, database.User, database.Password,
		database.DBName, database.Port, database.SSLMode)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}) //Postgres数据库连接

	//sqlite 参考
	// db, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

	//sql Server 参考
	//dsn := "sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm"
	//db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})

	//mysql 参考
	//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	DB, err := db.DB()
	DB.SetMaxIdleConns(10)  //设置空闲连接池中连接的最大数量
	DB.SetMaxOpenConns(100) //设置打开数据库连接的最大数量。
	//设置了连接可复用的最大时间 DB.SetConnMaxLifetime(1000)
	fmt.Printf("成功连接Postgres数据库,地址为%s\n", database.Host)
	return nil
}

// CloseDB 关闭DB连接/*
func CloseDB() error {
	dbSQL, err := db.DB()
	if err != nil {
		return err
	}
	return dbSQL.Close()
}

// AutoMigrate 自动迁移模式
func AutoMigrate() {
	//db.AutoMigrate(&Users{})
	//...
}
