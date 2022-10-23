package initialize

import (
	"fmt"
	"go_project/fish_farm/user_srv/global"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"go_project/fish_farm/user_srv/model"
)

//将数据库作为全局变量初始化
func InitDB() {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	mysqlConfig := global.ServerConfig.MysqlInfo
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlConfig.User, mysqlConfig.Password, mysqlConfig.Host, mysqlConfig.Port, mysqlConfig.Name)
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
	//此处 mysql 是第三方包 gorm.io/driver/mysql
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
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
	_ = global.DB.AutoMigrate(&model.User{})
}
