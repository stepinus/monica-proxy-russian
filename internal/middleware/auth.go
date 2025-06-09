package middleware

import (
	"monica-proxy/internal/config"
	"monica-proxy/internal/logger"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// BearerAuth 创建一个Bearer Token认证中间件
func BearerAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 获取Authorization header
			auth := c.Request().Header.Get("Authorization")

			// 检查header格式
			if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
				// 添加详细日志记录
				logger.Warn("无效的授权头",
					zap.String("method", c.Request().Method),
					zap.String("uri", c.Request().RequestURI),
					zap.String("remote_addr", c.RealIP()),
					zap.String("auth_header", auth),
				)
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header")
			}

			// 提取token
			token := strings.TrimPrefix(auth, "Bearer ")

			// 验证token
			if token != config.MonicaConfig.BearerToken || token == "" {
				// 使用结构化日志替代简单的log.Printf
				logger.Warn("无效的Token",
					zap.String("method", c.Request().Method),
					zap.String("uri", c.Request().RequestURI),
					zap.String("remote_addr", c.RealIP()),
					zap.String("token", token[0:4]+"..."), // 只显示部分token以保护安全
				)
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			return next(c)
		}
	}
}
