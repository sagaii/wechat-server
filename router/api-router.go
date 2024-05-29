package router

import (
	"wechat-server/controller"
	"wechat-server/middleware"

	"github.com/gin-gonic/gin"
)

func SetApiRouter(router *gin.Engine) {
	apiRouter := router.Group("/api")
	apiRouter.Use(middleware.GlobalAPIRateLimit())
	{
		apiRouter.GET("/status", controller.GetStatus)
		apiRouter.GET("/notice", controller.GetNotice)
		apiRouter.GET("/wechat", controller.WeChatVerification)
		apiRouter.POST("/wechat", controller.ProcessWeChatMessage)
		apiRouter.GET("/verification", middleware.CriticalRateLimit(), controller.SendEmailVerification)
		apiRouter.GET("/reset_password", middleware.CriticalRateLimit(), controller.SendPasswordResetEmail)
		apiRouter.GET("/user/reset", controller.SendNewPasswordEmail)
		apiRouter.GET("/oauth/github", controller.GitHubOAuth)

		userRoute := apiRouter.Group("/user")
		{
			userRoute.POST("/register", middleware.CriticalRateLimit(), controller.Register)
			userRoute.POST("/login", middleware.CriticalRateLimit(), controller.Login)
			userRoute.GET("/logout", controller.Logout)

			selfRoute := userRoute.Group("/")
			selfRoute.Use(middleware.UserAuth(), middleware.NoTokenAuth())
			{
				selfRoute.GET("/self", controller.GetSelf)
				selfRoute.PUT("/self", controller.UpdateSelf)
				selfRoute.DELETE("/self", controller.DeleteSelf)
				selfRoute.GET("/token", controller.GenerateToken)
			}

			adminRoute := userRoute.Group("/")
			adminRoute.Use(middleware.AdminAuth(), middleware.NoTokenAuth())
			{
				adminRoute.GET("/", controller.GetAllUsers)
				adminRoute.GET("/:id", controller.GetUser)
				adminRoute.POST("/", controller.CreateUser)
				adminRoute.POST("/manage", controller.ManageUser)
				adminRoute.PUT("/", controller.UpdateUser)
				adminRoute.DELETE("/:id", controller.DeleteUser)
			}
		}

		logRoute := apiRouter.Group("/logs")
		logRoute.Use(middleware.RootAuth(), middleware.NoTokenAuth())
		{
			logRoute.GET("/common", controller.GetCommonLogs)
			logRoute.GET("/error", controller.GetErrorLogs)
		}

		optionRoute := apiRouter.Group("/option")
		optionRoute.Use(middleware.RootAuth(), middleware.NoTokenAuth())
		{
			optionRoute.GET("/", controller.GetOptions)
			optionRoute.PUT("/", controller.UpdateOption)
		}
		fileRoute := apiRouter.Group("/file")
		{
			fileRoute.GET("/:id", middleware.DownloadRateLimit(), controller.DownloadFile)
			fileRoute.POST("/", middleware.UserAuth(), middleware.UploadRateLimit(), controller.UploadFile)
			fileRoute.DELETE("/:id", middleware.UserAuth(), controller.DeleteFile)
		}
		wechatRoute := apiRouter.Group("/wechat")
		wechatRoute.Use(middleware.AdminAuth(), middleware.TokenOnlyAuth())
		{
			wechatRoute.GET("/access_token", controller.GetAccessToken)
			wechatRoute.GET("/user", controller.GetUserIDByCode)
		}
	}
}
