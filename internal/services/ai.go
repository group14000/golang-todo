package services

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

const openRouterURL = "https://openrouter.ai/api/v1/chat/completions"

// AIService handles interaction with OpenRouter chat completions API.
type AIService struct {
	apiKey string
	client *http.Client
}

type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AIChatRequest struct {
	Model    string      `json:"model"`
	Messages []AIMessage `json:"messages"`
	Stream   bool        `json:"stream,omitempty"`
}

type AIChoice struct {
	Message AIMessage `json:"message"`
}

type AIChatResponse struct {
	Choices []AIChoice `json:"choices"`
	Error   any        `json:"error"`
}

func NewAIService(apiKey string) *AIService {
	return &AIService{
		apiKey: apiKey,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (s *AIService) Chat(ctx context.Context, messages []AIMessage) (string, error) {
	if s.apiKey == "" {
		return "", errors.New("AI service not configured")
	}
	payload := AIChatRequest{Model: "deepseek/deepseek-chat-v3.1:free", Messages: messages}
	b, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, openRouterURL, bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", errors.New("AI provider returned status " + resp.Status)
	}

	var parsed AIChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", err
	}
	if len(parsed.Choices) == 0 {
		return "", errors.New("no choices returned")
	}
	return parsed.Choices[0].Message.Content, nil
}

// StreamChunk represents a piece of streamed output or an error.
type StreamChunk struct {
	Text string
	Err  error
}

// ChatStream performs a streaming chat completion.
func (s *AIService) ChatStream(ctx context.Context, messages []AIMessage) (<-chan StreamChunk, error) {
	if s.apiKey == "" {
		return nil, errors.New("AI service not configured")
	}
	payload := AIChatRequest{Model: "deepseek/deepseek-chat-v3.1:free", Messages: messages, Stream: true}
	b, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, openRouterURL, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	ch := make(chan StreamChunk)
	go func() {
		defer resp.Body.Close()
		defer close(ch)
		if resp.StatusCode >= 400 {
			ch <- StreamChunk{Err: errors.New("AI provider returned status " + resp.Status)}
			return
		}
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "data: ") {
				data := strings.TrimPrefix(line, "data: ")
				if data == "[DONE]" {
					return
				}
				// OpenRouter streams partial JSON objects similar to OpenAI format: {"choices":[{"delta":{"content":"..."}}]}
				var raw map[string]any
				if err := json.Unmarshal([]byte(data), &raw); err != nil {
					continue
				}
				choices, _ := raw["choices"].([]any)
				if len(choices) == 0 {
					continue
				}
				first, _ := choices[0].(map[string]any)
				delta, _ := first["delta"].(map[string]any)
				content, _ := delta["content"].(string)
				if content != "" {
					select {
					case ch <- StreamChunk{Text: content}:
					case <-ctx.Done():
						return
					}
				}
			}
		}
		if err := scanner.Err(); err != nil {
			ch <- StreamChunk{Err: err}
		}
	}()

	return ch, nil
}
