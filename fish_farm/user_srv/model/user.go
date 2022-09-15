package model

import (
	"time"

	"gorm.io/gorm"
)

//数据表通用基础模型
type BaseModel struct {
	ID        int32     `gorm:"primarkey"`
	CreatedAt time.Time `gorm:"column:add_time"`
	UpdatedAt time.Time `gorm:"column:update_time"`
	DeletedAt gorm.DeletedAt
	IsDeleted bool `gorm:"column:is_deleted"` //软删除，上同DeletedAt
}

/*
用户表模型
1.密码采用密文保存，且不可反解
2.非对称加密：md5(信息摘要算法)
3.用户找回密码时通过密文比对来判断是否可修改
*/
type User struct {
	BaseModel
	Mobile   string     `gorm:"index:id_mobile;unique;type:varchar(11);not null"`
	Password string     `gorm:"type:varchar(100);not null"`
	Nickname string     `gorm:"type:varchar(20) comment '用户名'"`
	Birthday *time.Time `gorm:"type:datetime comment '出生日期'"`
	Gender   string     `gorm:"column:gender;default:male;type:varchar(6) comment 'female表示女,male表示男'"`
	Role     int        `gorm:"column:role;default:1;type:int comment '1表示普通用户，2表示管理员'"`
}
