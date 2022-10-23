package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"net/http"
)

/*
图片验证码
*/
var store = base64Captcha.DefaultMemStore

//定义返回图片验证码
func GetCaptcha(ctx *gin.Context) {
	//创建并配置驱动
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	//创建实例
	cp := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err := cp.Generate()
	if err != nil {
		zap.S().Errorf("生成验证码错误", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成验证码错误",
		})
		return
	}
	//将图片 id 与 图片数值文本 返回给前端
	ctx.JSON(http.StatusOK, gin.H{
		"captcha": id,
		"picPath": b64s,
	})
}
