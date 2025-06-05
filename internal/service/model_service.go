package service

import (
	"monica-proxy/internal/logger"
	"monica-proxy/internal/types"

	"go.uber.org/zap"
)

// ModelService 模型服务接口
type ModelService interface {
	// GetSupportedModels 获取支持的模型列表
	GetSupportedModels() []string
}

// modelService 模型服务实现
type modelService struct{}

// NewModelService 创建模型服务实例
func NewModelService() ModelService {
	return &modelService{}
}

// GetSupportedModels 获取支持的模型列表
func (s *modelService) GetSupportedModels() []string {
	models := types.GetSupportedModels()
	
	logger.Info("获取支持的模型列表",
		zap.Int("model_count", len(models)),
	)
	
	return models
}
