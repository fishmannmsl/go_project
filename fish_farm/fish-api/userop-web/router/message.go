package router

import (
	"github.com/gin-gonic/gin"
	
	"go_project/fish_farm/fish-api/userop-web/api/message"
	"go_project/fish_farm/fish-api/userop-web/middlewares"
)

func InitMessageRouter(Router *gin.RouterGroup) {
	MessageRouter := Router.Group("message").Use(middlewares.JWTAuth())
	{
		MessageRouter.GET("", message.List) // 轮播图列表页
		MessageRouter.POST("", message.New) //新建轮播图
	}
}
