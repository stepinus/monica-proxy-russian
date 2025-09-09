package monica

import (
	"context"
	"monica-proxy/internal/config"
	"monica-proxy/internal/errors"
	"monica-proxy/internal/logger"
	"monica-proxy/internal/types"
	"monica-proxy/internal/utils"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

// SendMonicaRequest 发起对 Monica AI 的请求(使用 resty)
func SendMonicaRequest(ctx context.Context, cfg *config.Config, mReq *types.MonicaRequest) (*resty.Response, error) {
	// Логируем детали запроса
	logger.Info("Отправка HTTP запроса к Monica API",
		zap.String("url", types.BotChatURL),
		zap.String("method", "POST"),
		zap.String("task_uid", mReq.TaskUID),
		zap.String("bot_uid", mReq.BotUID),
		zap.String("language", mReq.Language),
		zap.String("task_type", mReq.TaskType),
	)

	// Создаем запрос с дополнительными заголовками
	req := utils.RestySSEClient.R().
		SetContext(ctx).
		SetHeader("cookie", cfg.Monica.Cookie).
		SetHeader("x-client-locale", "ru_RU"). // Явно устанавливаем русский язык
		SetBody(mReq)

	// Логируем все заголовки перед отправкой
	headers := make(map[string]string)
	for k, v := range req.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}
	logger.Info("Финальные HTTP заголовки перед отправкой в Monica",
		zap.Any("all_headers", headers),
		zap.String("x_client_locale_value", headers["X-Client-Locale"]),
		zap.Bool("has_correct_locale", headers["X-Client-Locale"] == "ru_RU"),
	)

	// 发起请求
	resp, err := req.Post(types.BotChatURL)

	if err != nil {
		logger.Error("Monica API请求失败", zap.Error(err))
		return nil, errors.NewRequestFailedError("Monica API调用失败", err)
	}

	// 如果需要在这里做更多判断，可自行补充
	return resp, nil
}

// SendCustomBotRequest 发送custom bot请求
func SendCustomBotRequest(ctx context.Context, cfg *config.Config, customBotReq *types.CustomBotRequest) (*resty.Response, error) {
	// Логируем детали запроса
	logger.Info("Отправка HTTP запроса к Monica Custom Bot API",
		zap.String("url", types.CustomBotChatURL),
		zap.String("method", "POST"),
		zap.String("task_uid", customBotReq.TaskUID),
		zap.String("bot_uid", customBotReq.BotUID),
		zap.String("language", customBotReq.Language),
		zap.String("locale", customBotReq.Locale),
		zap.String("ai_resp_language", customBotReq.AIRespLanguage),
		zap.String("task_type", customBotReq.TaskType),
	)

	// Создаем запрос с дополнительными заголовками
	req := utils.RestySSEClient.R().
		SetContext(ctx).
		SetHeader("cookie", cfg.Monica.Cookie).
		SetHeader("x-client-locale", "ru_RU"). // Явно устанавливаем русский язык
		SetBody(customBotReq)

	// Логируем все заголовки перед отправкой
	headers := make(map[string]string)
	for k, v := range req.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}
	logger.Info("Финальные HTTP заголовки Custom Bot перед отправкой в Monica",
		zap.Any("all_headers", headers),
		zap.String("x_client_locale_value", headers["X-Client-Locale"]),
		zap.Bool("has_correct_locale", headers["X-Client-Locale"] == "ru_RU"),
	)

	// 发起请求
	resp, err := req.Post(types.CustomBotChatURL)

	if err != nil {
		logger.Error("Custom Bot API请求失败", zap.Error(err))
		return nil, errors.NewRequestFailedError("Custom Bot API调用失败", err)
	}

	return resp, nil
}
