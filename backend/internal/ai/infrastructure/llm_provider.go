package infrastructure

import (
	"context"
	"io"
)

type ProviderType string

const (
	ProviderTypeOpenAI   ProviderType = "openai"
	ProviderTypeOllama   ProviderType = "ollama"
	ProviderTypeQwen     ProviderType = "qwen"
	ProviderTypeDeepSeek ProviderType = "deepseek"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
}

type ChatResponse struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int64                  `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChatChoice           `json:"choices"`
	Usage   map[string]interface{} `json:"usage"`
}

type ChatChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

type ChatChunk struct {
	ID      string        `json:"id"`
	Object  string        `json:"object"`
	Created int64         `json:"created"`
	Model   string        `json:"model"`
	Choices []ChunkChoice `json:"choices"`
}

type ChunkChoice struct {
	Index        int         `json:"index"`
	Delta        ChatMessage `json:"delta"`
	FinishReason string      `json:"finish_reason"`
}

type ProviderConfig struct {
	Type        ProviderType
	APIKey      string
	BaseURL     string
	Model       string
	Temperature float64
	MaxTokens   int
	Timeout     int
}

type LLMProvider interface {
	Name() ProviderType
	ChatCompletion(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
	StreamChatCompletion(ctx context.Context, req *ChatRequest) (<-chan *ChatChunk, error)
	IsAvailable(ctx context.Context) bool
	GetModel() string
}

type ProviderFactory interface {
	CreateProvider(config *ProviderConfig) (LLMProvider, error)
	GetAvailableProviders() []ProviderType
}

type providerFactory struct {
	providers map[ProviderType]LLMProvider
}

func NewProviderFactory() ProviderFactory {
	return &providerFactory{
		providers: make(map[ProviderType]LLMProvider),
	}
}

func (f *providerFactory) CreateProvider(config *ProviderConfig) (LLMProvider, error) {
	if provider, exists := f.providers[config.Type]; exists {
		return provider, nil
	}

	var provider LLMProvider
	var err error

	switch config.Type {
	case ProviderTypeOpenAI:
		provider, err = NewOpenAIProvider(config)
	case ProviderTypeOllama:
		provider, err = NewOllamaProvider(config)
	case ProviderTypeQwen:
		provider, err = NewQwenProvider(config)
	case ProviderTypeDeepSeek:
		provider, err = NewDeepSeekProvider(config)
	default:
		return nil, ErrUnsupportedProvider
	}

	if err != nil {
		return nil, err
	}

	f.providers[config.Type] = provider
	return provider, nil
}

func (f *providerFactory) GetAvailableProviders() []ProviderType {
	return []ProviderType{
		ProviderTypeOpenAI,
		ProviderTypeOllama,
		ProviderTypeQwen,
		ProviderTypeDeepSeek,
	}
}

type BaseProvider struct {
	config *ProviderConfig
}

func (p *BaseProvider) GetModel() string {
	return p.config.Model
}

var (
	ErrUnsupportedProvider = &ProviderError{Message: "unsupported provider type"}
	ErrProviderUnavailable = &ProviderError{Message: "provider is unavailable"}
	ErrInvalidConfig       = &ProviderError{Message: "invalid provider configuration"}
)

type ProviderError struct {
	Message string
	Cause   error
}

func (e *ProviderError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

func (e *ProviderError) Unwrap() error {
	return e.Cause
}

type StreamReader interface {
	Read() (*ChatChunk, error)
	Close() error
}

type streamReader struct {
	reader  io.Reader
	decoder interface{}
}

func NewStreamReader(reader io.Reader) StreamReader {
	return &streamReader{
		reader: reader,
	}
}

func (r *streamReader) Read() (*ChatChunk, error) {
	return nil, nil
}

func (r *streamReader) Close() error {
	return nil
}
