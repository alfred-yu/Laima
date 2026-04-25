package infrastructure

import (
	"context"
)

type DeepSeekProvider struct {
	*OpenAIProvider
}

func NewDeepSeekProvider(config *ProviderConfig) (*DeepSeekProvider, error) {
	if config.APIKey == "" {
		return nil, ErrInvalidConfig
	}

	if config.BaseURL == "" {
		config.BaseURL = "https://api.deepseek.com/v1"
	}

	if config.Model == "" {
		config.Model = "deepseek-coder"
	}

	openaiProvider, err := NewOpenAIProvider(config)
	if err != nil {
		return nil, err
	}

	return &DeepSeekProvider{
		OpenAIProvider: openaiProvider,
	}, nil
}

func (p *DeepSeekProvider) Name() ProviderType {
	return ProviderTypeDeepSeek
}

func (p *DeepSeekProvider) IsAvailable(ctx context.Context) bool {
	return p.OpenAIProvider.IsAvailable(ctx)
}
