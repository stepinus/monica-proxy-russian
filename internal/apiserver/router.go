package apiserver

import (
	"monica-proxy/internal/middleware"
	"monica-proxy/internal/monica"
	"monica-proxy/internal/types"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sashabaranov/go-openai"
)

// RegisterRoutes 注册 Echo 路由
func RegisterRoutes(e *echo.Echo) {
	// 添加Bearer Token认证中间件
	e.Use(middleware.BearerAuth())

	// ChatGPT 风格的请求转发到 /v1/chat/completions
	e.POST("/v1/chat/completions", handleChatCompletion)
	// 获取支持的模型列表
	e.GET("/v1/models", handleListModels)
}

// handleChatCompletion 接收 ChatGPT 形式的对话请求并转发给 Monica
func handleChatCompletion(c echo.Context) error {
	var req openai.ChatCompletionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error": "Invalid request payload",
		})
	}

	// 检查请求是否包含消息
	if len(req.Messages) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error": "No messages found",
		})
	}

	// marshalIndent, err := json.MarshalIndent(req, "", "  ")
	// if err != nil {
	// 	return err
	// }
	// log.Printf("Received completion request: \n%s\n", marshalIndent)
	// 将 ChatGPTRequest 转换为 MonicaRequest
	monicaReq, err := types.ChatGPTToMonica(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
	}

	// 调用 Monica 并获取 SSE Stream
	stream, err := monica.SendMonicaRequest(c.Request().Context(), monicaReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
	}
	// Resty 不会自动关闭 Body，需要我们自己来处理
	defer stream.RawBody().Close()

	// 根据 stream 参数决定是否使用流式响应
	if req.Stream {
		// 使用流式响应
		c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
		c.Response().Header().Set("Cache-Control", "no-cache")
		c.Response().Header().Set("Transfer-Encoding", "chunked")
		c.Response().WriteHeader(http.StatusOK)

		// 将 Monica 的 SSE 数据逐行读出，再以 SSE 格式返回给调用方
		if err := monica.StreamMonicaSSEToClient(req.Model, c.Response().Writer, stream.RawBody()); err != nil {
			return err
		}
	} else {
		// 使用非流式响应
		c.Response().Header().Set(echo.HeaderContentType, "application/json")

		// 收集所有的 SSE 数据并转换为完整的响应
		response, err := monica.CollectMonicaSSEToCompletion(req.Model, stream.RawBody())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"error": err.Error(),
			})
		}

		// 返回完整的响应
		return c.JSON(http.StatusOK, response)
	}

	return nil
}

// handleListModels 返回支持的模型列表
func handleListModels(c echo.Context) error {
	models := types.GetSupportedModels()
	return c.JSON(http.StatusOK, models)
}
