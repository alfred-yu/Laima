package infrastructure

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type OpenAIProvider struct {
	*BaseProvider
	httpClient *http.Client
}

func NewOpenAIProvider(config *ProviderConfig) (*OpenAIProvider, error) {
	if config.APIKey == "" {
		return nil, ErrInvalidConfig
	}

	if config.BaseURL == "" {
		config.BaseURL = "https://api.openai.com/v1"
	}

	if config.Model == "" {
		config.Model = "gpt-4o"
	}

	if config.Timeout == 0 {
		config.Timeout = 60
	}

	return &OpenAIProvider{
		BaseProvider: &BaseProvider{config: config},
		httpClient: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
	}, nil
}

func (p *OpenAIProvider) Name() ProviderType {
	return ProviderTypeOpenAI
}

func (p *OpenAIProvider) ChatCompletion(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	if req.Temperature == 0 {
		req.Temperature = p.config.Temperature
	}
	if req.MaxTokens == 0 {
		req.MaxTokens = p.config.MaxTokens
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, &ProviderError{Message: "failed to marshal request", Cause: err}
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.config.BaseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, &ProviderError{Message: "failed to create request", Cause: err}
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.config.APIKey)

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

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, &ProviderError{Message: "failed to decode response", Cause: err}
	}

	return &chatResp, nil
}

func (p *OpenAIProvider) StreamChatCompletion(ctx context.Context, req *ChatRequest) (<-chan *ChatChunk, error) {
	if req.Temperature == 0 {
		req.Temperature = p.config.Temperature
	}
	if req.MaxTokens == 0 {
		req.MaxTokens = p.config.MaxTokens
	}
	req.Stream = true

	body, err := json.Marshal(req)
	if err != nil {
		return nil, &ProviderError{Message: "failed to marshal request", Cause: err}
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.config.BaseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, &ProviderError{Message: "failed to create request", Cause: err}
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.config.APIKey)

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

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				return
			}

			var chunk ChatChunk
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				continue
			}

			select {
			case chunkChan <- &chunk:
			case <-ctx.Done():
				return
			}
		}
	}()

	return chunkChan, nil
}

func (p *OpenAIProvider) IsAvailable(ctx context.Context) bool {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := p.ChatCompletion(ctx, &ChatRequest{
		Model:    p.config.Model,
		Messages: []ChatMessage{{Role: "user", Content: "ping"}},
		MaxTokens: 1,
	})
	return err == nil
}
