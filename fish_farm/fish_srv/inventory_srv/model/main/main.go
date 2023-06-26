package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"go_project/fish_farm/fish_srv/inventory_srv/model"
)

func genMd5(code string) string {
	Md5 := md5.New()
	_, _ = io.WriteString(Md5, code)
	return hex.EncodeToString(Md5.Sum(nil))
}

func main() {
	dsn := "root:root@tcp(192.168.0.104:3306)/mxshop_inventory_srv?charset=utf8mb4&parseTime=True&loc=Local"

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

	//_ = db.AutoMigrate(&model.Inventory{}, &model.StockSellDetail{})
	//插入一条数据
	//orderDetail := model.StockSellDetail{
	//	OrderSn: "imooc-bobby",
	//	Status:  1,
	//	Detail:  []model.GoodsDetail{{1,2},{2,3}},
	//}
	//db.Create(&orderDetail)

	var sellDetail model.StockSellDetail
	db.Where(model.StockSellDetail{OrderSn: "imooc-bobby"}).First(&sellDetail)
	fmt.Println(sellDetail.Detail)
}
