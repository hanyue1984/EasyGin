package models

import (
	Config "EasyGin/app/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

var db *gorm.DB

// ConnectDB 创建连接
func ConnectDB(database Config.Database) error {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		database.Host, database.User, database.Password,
		database.DBName, database.Port, database.SSLMode)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Silent), //关闭日志
	})

	//Postgres数据库连接
	//sqlite 参考
	// db, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

	//sql Server 参考
	//dsn := "sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm"
	//db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})

	//mysql 参考
	//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("connect db %s error: %v", database.DBName, err)
	}
	DB, err := db.DB()
	if err != nil {
		return fmt.Errorf("connect db %s error: %v", database.DBName, err)
	}
	DB.SetMaxIdleConns(10)  //设置空闲连接池中连接的最大数量
	DB.SetMaxOpenConns(100) //设置打开数据库连接的最大数量。
	//设置了连接可复用的最大时间
	DB.SetConnMaxLifetime(1000)
	fmt.Printf("成功连接Postgres数据库,地址为%s\n等待自动迁移数据库 \n", database.Host)

	err = db.AutoMigrate(&Users{}, &Accounts{}, &RbacRoles{}, &RbacMenu{}, &RbacMenuMeta{}, &RbacDept{})
	if err != nil {
		return fmt.Errorf("autoMigrate error: %v", err)
	}
	fmt.Printf("数据库迁移完成 \n")
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

func ConversionListData(count int64, page int, pageSize int, rows []gin.H) gin.H {
	return gin.H{
		"total":    count,
		"page":     page,
		"pageSize": pageSize,
		"rows":     rows,
	}
}
