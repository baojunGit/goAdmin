package main

import (
	"context"
	"fmt"
	"github.com/baojunGit/goAdmin/config"
	"github.com/baojunGit/goAdmin/user"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func main() {
	// 初始化配置
	// configDir参数为空字符串，它会使用默认的配置文件目录"./config"
	err := config.InitConfig(context.Background(), "")
	if err != nil {
		// 处理配置初始化错误
		fmt.Println("Failed to initialize config:", err)
		return
	}
	mySql := config.Config.MySql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s",
		mySql.Username, mySql.Password, mySql.Host, mySql.Port, mySql.Dbname, mySql.Timeout)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // 慢速 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别，学习的时候可以换成Info
			IgnoreRecordNotFoundError: true,        // 忽略记录器的 ErrRecordNotFound 错误
			ParameterizedQueries:      true,        // 不在 SQL 日志中包含参数
			Colorful:                  true,        // 彩色打印，学习的时候可以打开比较明显
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 打印原始的sql和错误
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
		return
	}
	fmt.Println("连接成功")

	// 使用 db 进行数据库操作
	// 不需要关闭数据库连接，Gorm v2 会自动管理

	// 创建一个 user 变量并进行更新
	result := db.Create(&user.User{Username: "Jinzhu", Password: "lbj*dddd"}) // 通过数据的指针来创建
	// 返回 error
	if result.Error != nil {
		fmt.Println("Failed to create user:", result.Error)
		return
	}
	fmt.Println("Rows affected:", result.RowsAffected)
	// 返回插入记录的条数
	if result.RowsAffected > 0 {
		fmt.Println("User created successfully!")
		return
	} else {
		fmt.Println("No rows were affected. User creation failed.")
	}
	fmt.Println("User updated successfully!")

}
