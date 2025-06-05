package middleware

import (
	"monica-proxy/internal/errors"
	"monica-proxy/internal/logger"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// ErrorHandler 创建统一的错误处理中间件
func ErrorHandler() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		// 获取请求ID
		requestID := c.Request().Header.Get(echo.HeaderXRequestID)

		// 处理应用错误
		if appErr, ok := err.(*errors.AppError); ok {
			status, response := appErr.HTTPResponse()

			// 添加请求ID到响应中
			if response["error"] != nil {
				if errMap, ok := response["error"].(map[string]any); ok {
					errMap["request_id"] = requestID
				}
			}

			// 记录错误日志
			logger.Error("应用错误",
				zap.Int("status", status),
				zap.Int("error_code", int(appErr.Code)),
				zap.String("error_msg", appErr.Message),
				zap.Error(appErr.Err),
			)

			c.JSON(status, response)
			return
		}

		// 处理Echo框架错误
		if echoErr, ok := err.(*echo.HTTPError); ok {
			status := echoErr.Code
			message := "服务器错误"
			if m, ok := echoErr.Message.(string); ok {
				message = m
			}

			// 构建响应
			response := map[string]any{
				"error": map[string]any{
					"code":       echoErr.Code,
					"message":    message,
					"request_id": requestID,
				},
			}

			// 记录错误日志
			logger.Error("框架错误",
				zap.Int("status", status),
				zap.String("error_msg", message),
				zap.Error(err),
			)

			c.JSON(status, response)
			return
		}

		// 处理其他错误
		status := http.StatusInternalServerError
		response := map[string]any{
			"error": map[string]any{
				"code":       status,
				"message":    "服务器内部错误",
				"request_id": requestID,
			},
		}

		// 记录错误日志
		logger.Error("未分类错误",
			zap.Int("status", status),
			zap.Error(err),
		)

		c.JSON(status, response)
	}
}
