package infrastructure

import (
	"context"
)

type QwenProvider struct {
	*OpenAIProvider
}

func NewQwenProvider(config *ProviderConfig) (*QwenProvider, error) {
	if config.APIKey == "" {
		return nil, ErrInvalidConfig
	}

	if config.BaseURL == "" {
		config.BaseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	}

	if config.Model == "" {
		config.Model = "qwen-coder-plus"
	}

	openaiProvider, err := NewOpenAIProvider(config)
	if err != nil {
		return nil, err
	}

	return &QwenProvider{
		OpenAIProvider: openaiProvider,
	}, nil
}

func (p *QwenProvider) Name() ProviderType {
	return ProviderTypeQwen
}

func (p *QwenProvider) IsAvailable(ctx context.Context) bool {
	return p.OpenAIProvider.IsAvailable(ctx)
}
