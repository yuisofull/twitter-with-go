package main

import (
	"github.com/gin-gonic/gin"
	"twitter/component/appctx"
	"twitter/memcache"
	"twitter/middleware"
	"twitter/modules/tweet/transport/gintwitter"
	"twitter/modules/upload/uploadtransport/ginupload"
	"twitter/modules/user/transport/ginuser"
	"twitter/modules/user/userstore"
)

func setupRoute(appCtx appctx.AppContext, v1 *gin.RouterGroup) {
	userStore := userstore.NewSQLStore(appCtx.GetMyDBConnection())
	userCachingStore := memcache.NewUserCaching(memcache.NewCaching(), userStore)

	//POST /v1/upload
	v1.POST("/upload", ginupload.Upload(appCtx))

	v1.POST("/register", ginuser.Register(appCtx))

	v1.POST("/authenticate", ginuser.Login(appCtx))

	v1.GET("/profile", middleware.RequireAuth(appCtx, userStore), ginuser.GetProfile(appCtx))
	{
		tweets := v1.Group("/tweets")
		tweets.POST("", middleware.RequireAuth(appCtx, userCachingStore), gintwitter.CreateTweet(appCtx))
		tweets.GET("", gintwitter.ListTweet(appCtx))
	}

}
