package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

/*
手机号码验证器，判断是否符号
*/
func ValidateMobile(f1 validator.FieldLevel) bool {
	mobile := f1.Field().String()
	//使用正则表达式判断是否合法,如果错误直接返回false
	ok, _ := regexp.MatchString(`^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`, mobile)
	if !ok {
		return false
	}
	return true
}
