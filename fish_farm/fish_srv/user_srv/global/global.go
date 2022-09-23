package global

import (
	"go_project/fish_farm/user_srv/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

//将DB作为全局变量
var (
	DB *gorm.DB
)

func init()  {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "root:123456@tcp(192.168.0.108:3306)/fish_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
	//配置日志输出，以便在执行gorm时可以查看底层的sql执行过程
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), //io writer输出到控制台
		logger.Config{
			SlowThreshold: time.Second, //sql语句输出间隔
			LogLevel:      logger.Info, //输出sql语句的重要等级
			Colorful:      true,        //输出sql语句是否禁用颜色
		},
	)
	//全局模式
	var err error
	DB,err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		//修改命名策略，使其直接用数据库模型名来建表
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	//设置全局的logger，在每次执行gorm时都会将底层的sql语句输出出来
	_ = DB.AutoMigrate(&model.User{})
}
