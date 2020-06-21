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

	// 测试接口都放到test下面
	test := router.Group("/test")
	{
		test.GET("/ping", apis.PingPong)
		test.GET("/transaction", apis.Transaction)
	}

	// 其它为各端提供API的接口放到api组下面
	api := router.Group("/api")
	{
		// 为APP提供的接口
		app := api.Group("/app")
		{
			app.POST("test", apis.PingPong)
		}

		// 为管理后台提供的接口
		admin := api.Group("/admin")
		{
			admin.POST("/upload", apis.ApiUploadAvatar)
		}

		// 为小程序提供的接口. 接口命名使用小写加下划线的方式，如这里不能叫：minipro
		miniProgram := api.Group("/mini-pro")
		{
			miniProgram.POST("/decode_phone_number", apis.ApiDecodePhoneNumber)
		}

		// 回调，走这里不鉴权
		callback := api.Group("/callback")
		{
			callback.POST("/wx_pay_callback", apis.ApiWXPayCallBack)
			callback.POST("/we_chat_callback", apis.UserWxCallbackHandler) // 公众号回调
		}

		// 登录，走这里不鉴权
		login := api.Group("/login")
		{
			login.POST("/app_login", apis.ApiLogin)
			login.POST("/mini_pro_login", apis.ApiLogin)
			login.POST("/admin_login", apis.ApiLogin)
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
