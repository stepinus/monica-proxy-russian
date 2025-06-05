package middleware

import (
	"monica-proxy/internal/logger"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// RequestLogger 创建一个请求日志记录中间件
func RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			req := c.Request()
			res := c.Response()

			// 生成请求ID
			requestID := req.Header.Get(echo.HeaderXRequestID)
			if requestID == "" {
				requestID = generateRequestID()
				c.Request().Header.Set(echo.HeaderXRequestID, requestID)
			}

			// 设置请求ID到响应头
			c.Response().Header().Set(echo.HeaderXRequestID, requestID)

			// 记录请求开始
			// logger.Info("请求开始",
			// 	zap.String("method", req.Method),
			// 	zap.String("uri", req.RequestURI),
			// 	zap.String("remote_addr", c.RealIP()),
			// )

			// 处理请求
			err := next(c)

			// 计算耗时
			duration := time.Since(start)

			// 记录请求结束
			fields := []zap.Field{
				zap.String("method", req.Method),
				zap.String("uri", req.RequestURI),
				zap.Int("status", res.Status),
				zap.Duration("latency", duration),
				zap.String("remote_addr", c.RealIP()),
			}

			if err != nil {
				// 记录错误
				fields = append(fields, zap.Error(err))
				logger.Error("请求失败", fields...)
			} else {
				logger.Info("请求完成", fields...)
			}

			return err
		}
	}
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	// 使用UUID生成唯一请求ID
	return "req_" + time.Now().Format("20060102150405") + "_" + randomString(8)
}

// randomString 生成指定长度的随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
