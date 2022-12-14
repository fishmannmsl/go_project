package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go_project/fish_farm/fish-api/user-web/global"
	"go_project/fish_farm/fish-api/user-web/utils"

	"go_project/fish_farm/fish-api/user-web/initialize"
	myvalidator "go_project/fish_farm/fish-api/user-web/validator"
)

func main() {
	//注册日志
	initialize.InitLogger()
	//注册配置文件
	initialize.InitConfig()
	//注册路由
	Router := initialize.Routers()
	//注册翻译
	if err := initialize.InitTrans("zh"); err != nil {
		panic(err)
	}
	//初始化 srv用户 连接
	initialize.InitSrvConn()
	//区分生产环境与开发环境，进而开发环境使用动态路由
	viper.AutomaticEnv()
	debug := viper.GetInt("ACSvcPort")
	if debug != 17532 {
		global.ServerConfig.Port, _ = utils.GetFreePort()
	}
	//注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", myvalidator.ValidateMobile)
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号码!", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	zap.S().Infof("启动服务器，端口:%d", global.ServerConfig.Port)
	if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动失败:", err.Error())
	}

}
