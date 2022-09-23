package model

import (
	"gorm.io/gorm"
	"time"
)

/*
在定义表模型时，一定要与proto文件内定义的表结构顺序相同
不然在进行分页查询时会报：panic: runtime error: index out of range [2] with length 1
原因是表结构不同
解决方法：
1.重新定义表结构
2.在分页查询的时添加主键
 */

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
	Password string     `gorm:"type:varchar(100);not null"`
	Mobile   string     `gorm:"index:id_mobile;unique;type:varchar(11);not null"`
	NickName string     `gorm:"type:varchar(20) comment '用户名'"`
	Birthday *time.Time `gorm:"type:datetime comment '出生日期'"`
	Gender   string     `gorm:"column:gender;default:male;type:varchar(6) comment 'female表示女,male表示男'"`
	Role     int        `gorm:"column:role;default:1;type:int comment '1表示普通用户，2表示管理员'"`
}
