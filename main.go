package main

import (
	"io"
	"monica-proxy/internal/apiserver"
	"monica-proxy/internal/config"
	"monica-proxy/internal/logger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func main() {
	// 从环境变量加载配置
	config.LoadConfig()

	// 设置日志级别
	logger.SetLevel("info")

	// 设置 Echo Server
	e := echo.New()
	e.Logger.SetOutput(io.Discard)

	// 添加基础中间件
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())

	// 注册路由
	apiserver.RegisterRoutes(e)

	// 启动服务器
	logger.Info("启动服务器", zap.String("port", "8080"))
	if err := e.Start(":8080"); err != nil {
		logger.Fatal("启动服务器失败", zap.Error(err))
	}
}
