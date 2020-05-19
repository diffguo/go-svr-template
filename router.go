package main

import (
	"github.com/gin-gonic/gin"
	"go-svr-template/views"
)

func addRoute(router *gin.Engine) {
	test := router.Group("/test")
	{
		test.GET("/ping", views.PingPong)
	}

	admin := router.Group("/admin")
	{
		admin.POST("/login", views.ApiLogin)
		admin.POST("/upload", views.ApiUploadAvatar)
	}

	miniProgram := router.Group("/minipro")
	{
		miniProgram.POST("decode_phone_number", views.ApiDecodePhoneNumber)
		miniProgram.POST("wx_pay_callback", views.ApiWXPayCallBack)
	}

	// 公众号
	wechat := router.Group("/wechat")
	{
		wechat.GET("/wx_callback", views.UserWxCallbackHandler)
		wechat.POST("/wx_callback", views.UserWxCallbackHandler)
	}
}
