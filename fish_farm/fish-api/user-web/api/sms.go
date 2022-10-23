package api

import (
	"context"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"go.uber.org/zap"
	"go_project/fish_farm/fish-api/user-web/forms"
	"go_project/fish_farm/fish-api/user-web/global"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

/*
 * 发送短信验证码
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		// 您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}

/*
 生成随机验证数字
 witdh 验证码长度
*/
func GenerateSmsCode(witdh int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < witdh; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

/*
由于接受信息手机号码是需要在阿里云绑定才能使用所以，
这里的注册手机号码并不是所有的号码都能接收短信，
统一由一个已绑定的号码来接收短信。
该功能只是为了测试调用AliSms可实现
*/
func SendSms(ctx *gin.Context) {
	//表单验证，格式是否正确，使用了第三方包 validator ,内置验证规则(正则)
	sendsms := forms.SendSmsForms{}
	if err := ctx.ShouldBindJSON(&sendsms); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	//初始化阿里云 SMS
	client, _err := CreateClient(tea.String(global.ServerConfig.AliSmsInfo.AliKey), tea.String(global.ServerConfig.AliSmsInfo.AliSecret))
	if _err != nil {
		zap.S().Errorf("使用AK&SK初始化账号 Client 异常", _err.Error())
	}
	smsCode := GenerateSmsCode(4)
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:      tea.String(global.ServerConfig.AliSmsInfo.AliSignName),
		TemplateCode:  tea.String(global.ServerConfig.AliSmsInfo.AliTemplateCode),
		PhoneNumbers:  tea.String(global.ServerConfig.AliSmsInfo.AliPhoneNumbers),
		TemplateParam: tea.String("{\"code\":" + smsCode + "}"),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		_, _err = client.SendSmsWithOptions(sendSmsRequest, runtime)
		if _err != nil {
			return _err
		}

		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		// 如有需要，请打印 error
		_, _err = util.AssertAsString(error.Message)
		if _err != nil {
			zap.S().Errorf("", _err.Error())
		}
	}

	//初始化一个redis客户端
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	//将短信有效时长，设置为5分钟
	rdb.Set(context.Background(), sendsms.Mobile, smsCode, time.Duration(global.ServerConfig.RedisInfo.Expire)*time.Minute)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "发送成功",
	})
}
