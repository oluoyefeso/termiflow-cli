package search

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/oluoyefeso/termiflow-cli/internal/config"
)

type SearchResult struct {
	Title       string
	URL         string
	Snippet     string
	Content     string
	PublishedAt time.Time
	Source      string
}

type SearchRequest struct {
	Query      string
	MaxResults int
	TimeRange  string // "day", "week", "month", "year"
}

type Provider interface {
	Name() string
	Search(ctx context.Context, req SearchRequest) ([]SearchResult, error)
	Available() bool
}

// GetSearchProvider returns the managed search provider.
func GetSearchProvider(cfg *config.Config) (Provider, error) {
	if os.Getenv("TERMIFLOW_MOCK") == "true" {
		p := NewMockSearchProvider()
		if p == nil {
			return nil, fmt.Errorf("TERMIFLOW_MOCK=true requires building with -tags mock: go build -tags mock ./...")
		}
		return p, nil
	}
	return NewManagedSearchProvider(
		cfg.Providers.Managed.APIKey,
		cfg.Providers.Managed.BaseURL,
	), nil
}
