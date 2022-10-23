package model

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

// BaseModel
//数据通用基础模型.
//使用同样的类型，是为了统一标准，防止建立外键时类型不匹配.
type BaseModel struct {
	ID        int32          `gorm:"primarkey;type:int" json:"id"`
	CreatedAt time.Time      `gorm:"column:add_time" json:"-"`
	UpdatedAt time.Time      `gorm:"column:update_time" json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
	IsDeleted bool           `gorm:"column:is_deleted" json:"-"` //软删除，上同DeletedAt
}

// GormList (将切片类型转化为数组内可存储类型)
//通过实现 Gorm 自带的 sql.Scanner,driver.Value即可完成自定义类型定义
type GormList []string

// Value 实现 driver.Value 接口，Value 返回 json value
func (g GormList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (g *GormList) Scan(value interface{}) error {
	//使用断言，将 value 确定为 []byte 类型
	return json.Unmarshal(value.([]byte), &g)
}
