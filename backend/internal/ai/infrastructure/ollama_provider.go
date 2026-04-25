package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OllamaProvider struct {
	*BaseProvider
	httpClient *http.Client
}

type OllamaRequest struct {
	Model    string                 `json:"model"`
	Messages []ChatMessage          `json:"messages"`
	Stream   bool                   `json:"stream"`
	Options  map[string]interface{} `json:"options,omitempty"`
}

type OllamaResponse struct {
	Model     string      `json:"model"`
	CreatedAt time.Time   `json:"created_at"`
	Message   ChatMessage `json:"message"`
	Done      bool        `json:"done"`
}

type OllamaStreamResponse struct {
	Model     string      `json:"model"`
	CreatedAt time.Time   `json:"created_at"`
	Message   ChatMessage `json:"message"`
	Done      bool        `json:"done"`
}

func NewOllamaProvider(config *ProviderConfig) (*OllamaProvider, error) {
	if config.BaseURL == "" {
		config.BaseURL = "http://localhost:11434"
	}

	if config.Model == "" {
		config.Model = "qwen2.5-coder:32b"
	}

	if config.Timeout == 0 {
		config.Timeout = 120
	}

	return &OllamaProvider{
		BaseProvider: &BaseProvider{config: config},
		httpClient: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
	}, nil
}

func (p *OllamaProvider) Name() ProviderType {
	return ProviderTypeOllama
}

func (p *OllamaProvider) ChatCompletion(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	ollamaReq := OllamaRequest{
		Model:    p.config.Model,
		Messages: req.Messages,
		Stream:   false,
		Options: map[string]interface{}{
			"temperature": p.config.Temperature,
			"num_predict": p.config.MaxTokens,
		},
	}

	body, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, &ProviderError{Message: "failed to marshal request", Cause: err}
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.config.BaseURL+"/api/chat", bytes.NewReader(body))
	if err != nil {
		return nil, &ProviderError{Message: "failed to create request", Cause: err}
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, &ProviderError{Message: "failed to send request", Cause: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, &ProviderError{
			Message: fmt.Sprintf("API request failed with status %d", resp.StatusCode),
			Cause:   fmt.Errorf("%s", string(bodyBytes)),
		}
	}

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, &ProviderError{Message: "failed to decode response", Cause: err}
	}

	chatResp := &ChatResponse{
		ID:      fmt.Sprintf("ollama-%d", time.Now().Unix()),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   ollamaResp.Model,
		Choices: []ChatChoice{
			{
				Index:        0,
				Message:      ollamaResp.Message,
				FinishReason: "stop",
			},
		},
	}

	return chatResp, nil
}

func (p *OllamaProvider) StreamChatCompletion(ctx context.Context, req *ChatRequest) (<-chan *ChatChunk, error) {
	ollamaReq := OllamaRequest{
		Model:    p.config.Model,
		Messages: req.Messages,
		Stream:   true,
		Options: map[string]interface{}{
			"temperature": p.config.Temperature,
			"num_predict": p.config.MaxTokens,
		},
	}

	body, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, &ProviderError{Message: "failed to marshal request", Cause: err}
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.config.BaseURL+"/api/chat", bytes.NewReader(body))
	if err != nil {
		return nil, &ProviderError{Message: "failed to create request", Cause: err}
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, &ProviderError{Message: "failed to send request", Cause: err}
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, &ProviderError{
			Message: fmt.Sprintf("API request failed with status %d", resp.StatusCode),
			Cause:   fmt.Errorf("%s", string(bodyBytes)),
		}
	}

	chunkChan := make(chan *ChatChunk, 100)

	go func() {
		defer close(chunkChan)
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		for {
			var streamResp OllamaStreamResponse
			if err := decoder.Decode(&streamResp); err != nil {
				if err == io.EOF {
					return
				}
				continue
			}

			chunk := &ChatChunk{
				ID:      fmt.Sprintf("ollama-%d", time.Now().Unix()),
				Object:  "chat.completion.chunk",
				Created: time.Now().Unix(),
				Model:   streamResp.Model,
				Choices: []ChunkChoice{
					{
						Index:        0,
						Delta:        streamResp.Message,
						FinishReason: "",
					},
				},
			}

			if streamResp.Done {
				chunk.Choices[0].FinishReason = "stop"
			}

			select {
			case chunkChan <- chunk:
			case <-ctx.Done():
				return
			}

			if streamResp.Done {
				return
			}
		}
	}()

	return chunkChan, nil
}

func (p *OllamaProvider) IsAvailable(ctx context.Context) bool {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(ctx, "GET", p.config.BaseURL+"/api/tags", nil)
	if err != nil {
		return false
	}

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
