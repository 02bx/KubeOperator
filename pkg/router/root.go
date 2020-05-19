package router

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "ko3-gin/docs"
	"ko3-gin/pkg/middleware"
	v1 "ko3-gin/pkg/router/v1"
	"net/http"
	"os"
)

func Server() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	f, _ := os.Open(os.DevNull)
	gin.DefaultWriter = f
	server := gin.Default()
	server.StaticFS("static", http.Dir("resource/static"))
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	server.Use(middleware.LoggerMiddleware())
	server.Use(middleware.PagerMiddleware())
	api := server.Group("/api")
	{
		v1.V1(api)
	}

	//jwtMiddleware := middleware.JWTMiddleware()
	//auth := server.Group("/auth")
	//{
	//	auth.POST("/login", jwtMiddleware.LoginHandler)
	//	auth.GET("/refresh", jwtMiddleware.RefreshHandler)
	//}
	//api := server.Group("/api")
	//api.Use(jwtMiddleware.MiddlewareFunc())
	//{
	//	pkg_api.V1(api)
	//}
	return server
}
