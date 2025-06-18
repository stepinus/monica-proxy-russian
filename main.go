package main

import (
	"fmt"
	"io"
	"monica-proxy/internal/apiserver"
	"monica-proxy/internal/config"
	"monica-proxy/internal/logger"
	"monica-proxy/internal/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// 设置日志级别
	logger.SetLevel(cfg.Logging.Level)

	// 创建应用实例
	app := newApp(cfg)

	// 启动服务器
	logger.Info("启动服务器", zap.String("address", cfg.GetAddress()))

	if err := app.Start(); err != nil {
		logger.Fatal("启动服务器失败", zap.Error(err))
	}
}

// App 应用实例
type App struct {
	config *config.Config
	server *echo.Echo
}

// newApp 创建应用实例
func newApp(cfg *config.Config) *App {
	// 初始化HTTP客户端
	utils.InitHTTPClients(cfg)

	// 设置 Echo Server
	e := echo.New()
	e.Logger.SetOutput(io.Discard)
	e.HideBanner = true

	// 配置服务器
	e.Server.ReadTimeout = cfg.Server.ReadTimeout
	e.Server.WriteTimeout = cfg.Server.WriteTimeout
	e.Server.IdleTimeout = cfg.Server.IdleTimeout

	// 添加基础中间件
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())

	// 注册路由
	apiserver.RegisterRoutes(e, cfg)

	return &App{
		config: cfg,
		server: e,
	}
}

// Start 启动应用
func (a *App) Start() error {
	return a.server.Start(a.config.GetAddress())
}
