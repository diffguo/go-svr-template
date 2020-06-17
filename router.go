package main

import (
	"github.com/diffguo/gocom"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // io-swagger middleware
	"go-svr-template/apis"
)

func addRoute(router *gin.Engine) {

	gocom.AddPProfHandler(router)

	test := router.Group("/test")
	{
		test.GET("/ping", apis.PingPong)
	}

	admin := router.Group("/admin")
	{
		admin.POST("/login", apis.ApiLogin)
		admin.POST("/upload", apis.ApiUploadAvatar)
	}

	miniProgram := router.Group("/minipro")
	{
		miniProgram.POST("decode_phone_number", apis.ApiDecodePhoneNumber)
		miniProgram.POST("wx_pay_callback", apis.ApiWXPayCallBack)
	}

	// 公众号
	wechat := router.Group("/wechat")
	{
		wechat.GET("/wx_callback", apis.UserWxCallbackHandler)
		wechat.POST("/wx_callback", apis.UserWxCallbackHandler)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
