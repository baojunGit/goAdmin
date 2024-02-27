package initialize

import (
	"context"
	"fmt"
	"github.com/baojunGit/goAdmin/config"
	"github.com/baojunGit/goAdmin/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

func InitDB() {
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

	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 打印原始的sql和错误
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 表明用英文单数形式的表名
		},
	})
	if err != nil {
		// panic可以记录日志，如果没有捕获和处理panic，panic会终止程序
		panic("数据库连接失败" + err.Error())
		return
	}
	// 创建account表，第一次用
	err = DB.AutoMigrate(&model.Account{})
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("连接成功")
}
