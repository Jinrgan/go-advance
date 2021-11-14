package main

import (
	"go-advance/week04/app/goods/internal/data"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func main() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "root:root@tcp(127.0.0.1:3306)/mxshop_goods?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // 禁用彩色打印
		},
	)

	// 全局模式
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&data.Category{}, &data.Brand{}, &data.Banner{}, &data.Goods{})
	if err != nil {
		panic(err)
	}

	//options := &password.Options{SaltLen: 10, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	//salt, encodedPwd := password.Encode("admin123", options)
	//newPwd := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	//
	//for i := 0; i < 10; i++ {
	//	user := User{
	//		Mobile:   fmt.Sprintf("1876666888%d", i),
	//		Password: newPwd,
	//		NickName: fmt.Sprintf("bobby%d", i),
	//	}
	//	db.Create(&user)
	//}
}
