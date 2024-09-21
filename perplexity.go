package perplexity

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Main URL
const ppxtyURL = "https://api.perplexity.ai/chat/completions"

// Constants for available models
const (
	// Perplexity Sonar Models - Online
	ModelLlama31SonarSmall128kOnline = "llama-3.1-sonar-small-128k-online"
	ModelLlama31SonarLarge128kOnline = "llama-3.1-sonar-large-128k-online"
	ModelLlama31SonarHuge128kOnline  = "llama-3.1-sonar-huge-128k-online"

	// Perplexity Chat Models
	ModelLlama31SonarSmall128kChat = "llama-3.1-sonar-small-128k-chat"
	ModelLlama31SonarLarge128kChat = "llama-3.1-sonar-large-128k-chat"

	// Open-Source Models
	ModelLlama31_8BInstruct  = "llama-3.1-8b-instruct"
	ModelLlama31_70BInstruct = "llama-3.1-70b-instruct"
)

// Constants for message roles
const (
	RoleSystem    = "system"
	RoleUser      = "user"
	RoleAssistant = "assistant"
)

// Message represents a message in the chat completion request
type Message struct {
	Role    string `json:"role"`    // Role of the message sender (e.g., "system", "user", "assistant")
	Content string `json:"content"` // Content of the message
}

// ChatCompletionRequest represents the request body for the chat completion API
type ChatCompletionRequest struct {
	Model            string    `json:"model"`                       // Model to use for the completion
	Messages         []Message `json:"messages"`                    // List of messages in the conversation
	MaxTokens        int       `json:"max_tokens,omitempty"`        // Maximum number of tokens to generate
	Temperature      float64   `json:"temperature,omitempty"`       // Sampling temperature
	TopP             float64   `json:"top_p,omitempty"`             // Nucleus sampling probability
	TopK             int       `json:"top_k,omitempty"`             // Top-K sampling
	FrequencyPenalty float64   `json:"frequency_penalty,omitempty"` // Frequency penalty
	PresencePenalty  float64   `json:"presence_penalty,omitempty"`  // Presence penalty
}

// ChatCompletionResponse represents the response from the chat completion API
type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Message      Message `json:"message"`
		FinishReason string  `json:"finish_reason"`
		Index        int     `json:"index"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// ClientOptions represents additional options to adjust the client behavior
type ClientOptions struct {
	RequestTimeoutInSeconds time.Duration
}

func (r ChatCompletionResponse) isSingle() bool {
	return len(r.Choices) == 1
}

func (r ChatCompletionResponse) isComplete() bool {
	for _, key := range r.Choices {
		if key.FinishReason == "stop" {
			return true
		}
	}
	return false
}

func (r ChatCompletionResponse) GetCompleteSingleMessage() (string, error) {
	switch {
	case !r.isSingle():
		return "", fmt.Errorf("there more than 1 choice in response")
	case !r.isComplete():
		return "", fmt.Errorf("choice is not complete")
	}
	return r.Choices[0].Message.Content, nil
}

// ValidationError represents the structure of a validation error response
type ValidationError struct {
	Detail []struct {
		Loc  []interface{} `json:"loc"`
		Msg  string        `json:"msg"`
		Type string        `json:"type"`
	} `json:"detail"`
}

// Client represents a client for the Perplexity API
type Client struct {
	APIKey     string
	Model      string
	HTTPClient *http.Client
}

var defaultClientOptions = &ClientOptions{
	RequestTimeoutInSeconds: 30,
}

// NewClient creates a new Perplexity API client
func NewClient(apiKey, model string, options *ClientOptions) *Client {
	if options == nil {
		options = defaultClientOptions
	}

	return &Client{
		APIKey: apiKey,
		Model:  model,
		HTTPClient: &http.Client{
			Timeout: options.RequestTimeoutInSeconds * time.Second,
		},
	}
}

// ChatCompletions sends a chat completion request to the Perplexity API
func (c *Client) ChatCompletions(ctx context.Context, request ChatCompletionRequest) (*ChatCompletionResponse, error) {
	// Set the model if not already set in the request
	if request.Model == "" {
		request.Model = c.Model
	}

	payload, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", ppxtyURL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		if res.StatusCode == http.StatusUnprocessableEntity {
			var validationError ValidationError
			if err := json.Unmarshal(body, &validationError); err != nil {
				return nil, fmt.Errorf("request failed with status %d: %s", res.StatusCode, string(body))
			}
			return nil, fmt.Errorf("validation error: %v", validationError)
		}
		return nil, fmt.Errorf("error: %s", body)
	}

	var response ChatCompletionResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
