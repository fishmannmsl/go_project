package main

import (
	"crypto/sha512"
	"fmt"
	"go_project/fish_farm/fish_srv/user_srv/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

/*
用以生成模拟用户数据，用来做测试
*/
func main() {
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
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		//修改命名策略，使其直接用数据库模型名来建表
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	options := &password.Options{16, 100, 32, sha512.New} //md5.New可改为sha512.New()
	salt, encodedPwd := password.Encode("admin123", options)
	newPasswd := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	fmt.Println(newPasswd)

	for i := 0; i < 10; i++ {
		user := model.User{
			NickName: fmt.Sprintf("dsxg%d", i),
			Mobile:   fmt.Sprintf("1376162457%d", i),
			Password: newPasswd,
		}
		db.Save(&user)
	}
	////设置全局的logger，在每次执行gorm时都会将底层的sql语句输出出来
	//_ = db.AutoMigrate(&model.User{})

}
