package initialize

import (
	"go.uber.org/zap"
)

func InitLogger() {
	//zap为第三方包 "go.uber.org/zap"
	logger, err := zap.NewProduction()
	if err != nil {
		zap.S().Panic("日志配置失败：", err.Error())
	}
	zap.ReplaceGlobals(logger)
}
