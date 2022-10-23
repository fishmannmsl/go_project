package initialize

import "go.uber.org/zap"

func InitLogger() {
	/*
		 配置全局logger变量
		1.S()可以获取一个全局的sugar,也可以让我们自己设置一个全局的logger
		2.日志是分级别的：debug，info，warn，error，fetal
		3.S函数和L函数提供了一个安全访问logger的途径
	*/
	logger, err := zap.NewProduction()
	if err != nil {
		zap.S().Panic("日志配置失败：", err.Error())
	}
	zap.ReplaceGlobals(logger)
}
