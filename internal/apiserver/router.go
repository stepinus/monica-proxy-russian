package apiserver

import (
	"io"
	"monica-proxy/internal/config"
	"monica-proxy/internal/errors"
	"monica-proxy/internal/middleware"
	"monica-proxy/internal/monica"
	"monica-proxy/internal/service"
	"monica-proxy/internal/types"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sashabaranov/go-openai"
)


// RegisterRoutes 注册 Echo 路由
func RegisterRoutes(e *echo.Echo, cfg *config.Config) {
	// 设置自定义错误处理器
	e.HTTPErrorHandler = middleware.ErrorHandler()

	// 添加中间件
	e.Use(middleware.BearerAuth(cfg))
	e.Use(middleware.RequestLogger(cfg))

	// 初始化服务实例
	chatService := service.NewChatService(cfg)
	modelService := service.NewModelService(cfg)
	imageService := service.NewImageService(cfg)

	// ChatGPT 风格的请求转发到 /v1/chat/completions
	e.POST("/v1/chat/completions", createChatCompletionHandler(chatService))
	// 获取支持的模型列表
	e.GET("/v1/models", createListModelsHandler(modelService))
	// DALL-E 风格的图片生成请求
	e.POST("/v1/images/generations", createImageGenerationHandler(imageService))
}

// createChatCompletionHandler 创建聊天完成处理器
func createChatCompletionHandler(chatService service.ChatService) echo.HandlerFunc {
	return func(c echo.Context) error {
	var req openai.ChatCompletionRequest
	if err := c.Bind(&req); err != nil {
		return errors.NewBadRequestError("无效的请求数据", err)
	}

	// 调用服务处理请求
	result, err := chatService.HandleChatCompletion(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	// 根据请求参数决定响应方式
	if req.Stream {
		// 对于流式请求，result是一个io.ReadCloser
		rawBody, ok := result.(io.Reader)
		if !ok {
			return errors.NewInternalError(nil)
		}

		// 确保关闭响应体
		closer, isCloser := rawBody.(io.Closer)
		if isCloser {
			defer closer.Close()
		}

		// 设置响应头
		c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
		c.Response().Header().Set("Cache-Control", "no-cache")
		c.Response().Header().Set("Transfer-Encoding", "chunked")
		c.Response().WriteHeader(http.StatusOK)

		// 流式处理响应
		if err := monica.StreamMonicaSSEToClient(req.Model, c.Response().Writer, rawBody); err != nil {
			return errors.NewInternalError(err)
		}
		return nil
	} else {
		// 对于非流式请求，直接返回JSON响应
		return c.JSON(http.StatusOK, result)
	}
	}
}

// createListModelsHandler 创建模型列表处理器
func createListModelsHandler(modelService service.ModelService) echo.HandlerFunc {
	return func(c echo.Context) error {
	// 调用服务获取模型列表
	models := modelService.GetSupportedModels()

	// 构造响应格式
	result := make(map[string][]struct {
		Id string `json:"id"`
	})

	result["data"] = make([]struct {
		Id string `json:"id"`
	}, 0)

	for _, model := range models {
		result["data"] = append(result["data"], struct {
			Id string `json:"id"`
		}{
			Id: model,
		})
	}
	return c.JSON(http.StatusOK, result)
	}
}

// createImageGenerationHandler 创建图片生成处理器
func createImageGenerationHandler(imageService service.ImageService) echo.HandlerFunc {
	return func(c echo.Context) error {
	// 解析请求
	var req types.ImageGenerationRequest
	if err := c.Bind(&req); err != nil {
		return errors.NewBadRequestError("无效的请求数据", err)
	}

	// 调用服务生成图片
	resp, err := imageService.GenerateImage(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	// 返回结果
	return c.JSON(http.StatusOK, resp)
	}
}
