package http

import (
	"fmt"
	"gameserver/api/middleware"
	"gameserver/api/routes"
	"gameserver/core/logger"
	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
)

func NewGinServer() {
	gin.SetMode(gin.DebugMode)
	//创建gin引擎
	engine := gin.Default()
	//添加zap日志中间件
	engine.Use(logger.GinLogger(logger.Logger), logger.GinRecovery(logger.Logger, true))
	//添加websocket服务

	if viper.GetBool("app.openWs") {
		engine.GET("/ws", middleware.WsHandler)
	}
	//添加预过滤器
	//engine.Use(middleware.MD5Auth())
	//添加预过滤器
	engine.Use(middleware.JWTAuth())
	//设置路由
	routes.New(engine)
	port := viper.GetString("app.httpPort")
	if port == "" {
		port = ":8888"
	}
	logger.Logger.Info(fmt.Sprintf("start http server success, port %s", port))
	//开启服务器，不填默认监听localhost:8080
	err := engine.Run(port)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("start http server error：%v", err))
		panic(err)
	}
}
