package response

import (
	"fmt"
	"time"
)

type JsonTime time.Time

/*
直接实现time.Time内置方法 MarshalJSON 来完成类型转换
*/
func (j JsonTime) MarshalJSON() ([]byte, error) {
	var stmp = fmt.Sprintf("\"%s\"", time.Time(j).Format("1990-01-01"))
	return []byte(stmp), nil
}

type UserResponse struct {
	Id       int32  `json:"id"`
	PassWord string `json:"pass_word"`
	Mobile   string `json:"mobile"`
	NickName string `json:"nick_name"`
	//Birthday string `json:"birthday"`
	Birthday JsonTime //直接调用JsonTime内置方法完成类型转换，实现用 time.Time 类型
	Gender   string   `json:"gender"`
	Role     uint32   `json:"role"`
}
