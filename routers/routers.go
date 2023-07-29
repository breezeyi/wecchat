package routers

import (
	"wexin/controllers"

	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {

	r := gin.Default()
	// 微信公众号接口配置验证
	r.GET("/wechat", controllers.Wechat)
	//接受用户发送的消息
	r.POST("/wechat", controllers.PostWecaht)
	return r
}
