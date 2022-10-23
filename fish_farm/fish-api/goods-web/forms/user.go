package forms

/*
使用密码登录模板
*/
type PassWordLoginForms struct {
	Mobile    string `form:"mobile" json:"mobile" binding:"required,mobile"`           //加入自定义手机号验证器
	Password  string `form:"password" json:"password" binding:"required,min=3,max=10"` //密码、约束条件不能加空格
	Captcha   string `form:"captcha" json:"captcha" binding:"required,min=5,max=5"`    //图片验证码数字
	CaptchaId string `form:"captcha_id" json:"captcha_id" binding:"required"`          //图形验证码ID
}

/*
用户注册信息模板
*/
type RegisterForm struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile"`           //加入自定义手机号验证器
	Password string `form:"password" json:"password" binding:"required,min=3,max=10"` //密码
	Code     string `form:"code" json:"code" binding:"required,min=4,max=4"`          //信息验证码
}
