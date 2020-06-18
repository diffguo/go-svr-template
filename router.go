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

	api := router.Group("/api")
	{
		// APP
		app := api.Group("/app")
		{
			app.POST("test", apis.PingPong)
		}

		// 管理后台
		admin := api.Group("/admin")
		{
			admin.POST("/login", apis.ApiLogin)
			admin.POST("/upload", apis.ApiUploadAvatar)
		}

		// 小程序
		miniProgram := api.Group("/minipro")
		{
			miniProgram.POST("decode_phone_number", apis.ApiDecodePhoneNumber)
		}

		// 回调，走这里不鉴权
		callback := api.Group("/callback")
		{
			miniProgram.POST("wx_pay_callback", apis.ApiWXPayCallBack)
			callback.GET("/wechat_callback", apis.UserWxCallbackHandler)  // 公众号回调
			callback.POST("/wechat_callback", apis.UserWxCallbackHandler) // 公众号回调
		}

		// 登录，走这里不鉴权
		login := api.Group("/login")
		{
			login.POST("app_login", apis.ApiLogin)
			login.POST("mini_pro_login", apis.ApiLogin)
			login.POST("admin_login", apis.ApiLogin)
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
