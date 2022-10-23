package forms

type SendSmsForms struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"`           //手机号验证器
	Type   string `form:"type" json:"type" binding:"required,oneof=register login"` //注册发送登录，注册验证码(注册短信与修改密码等是不同类别的)
}
