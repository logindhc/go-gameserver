package http

import (
	"fmt"
	"gameserver/appserver/api/middleware"
	"gameserver/appserver/api/routes"
	"gameserver/common/config"
	"gameserver/common/logger"
	"github.com/gin-gonic/gin"
)

func NewGinServer() {
	gin.SetMode(gin.DebugMode)
	//创建gin引擎
	engine := gin.Default()
	//添加zap日志中间件
	engine.Use(logger.GinLogger(logger.Logger), logger.GinRecovery(logger.Logger, true))
	//添加websocket服务
	if config.ServerConfig.OpenWS {
		engine.GET("/ws", middleware.WsHandler)
	}
	//添加预过滤器
	//engine.Use(middleware.MD5Auth())
	//添加预过滤器
	engine.Use(middleware.JWTAuth())
	//设置路由
	routes.New(engine)
	port := config.ServerConfig.HttpPort
	if port == "" {
		return
	}
	logger.Logger.Info(fmt.Sprintf("start http server success, port %s", port))
	//开启服务器，不填默认监听localhost:8080
	err := engine.Run(port)
	if err != nil {
		logger.Logger.Info(fmt.Sprintf("start http server error：%v", err))
		panic(err)
	}
}
