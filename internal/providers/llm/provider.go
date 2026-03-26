package llm

import (
	"context"
	"fmt"
	"os"

	"github.com/oluoyefeso/termiflow-cli/internal/config"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CompletionRequest struct {
	Messages    []Message
	MaxTokens   int
	Temperature float64
	Stream      bool
}

type CompletionResponse struct {
	Content      string
	FinishReason string
	Usage        Usage
}

type Usage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

type StreamChunk struct {
	Content string
	Done    bool
	Error   error
}

type Provider interface {
	Name() string
	Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error)
	Stream(ctx context.Context, req CompletionRequest) (<-chan StreamChunk, error)
	Available() bool
}

func GetProvider(_ string, cfg *config.Config) (Provider, error) {
	if os.Getenv("TERMIFLOW_MOCK") == "true" {
		p := NewMockProvider()
		if p == nil {
			return nil, fmt.Errorf("TERMIFLOW_MOCK=true requires building with -tags mock: go build -tags mock ./...")
		}
		return p, nil
	}
	return NewManagedProvider(
		cfg.Providers.Managed.APIKey,
		cfg.Providers.Managed.BaseURL,
	), nil
}
