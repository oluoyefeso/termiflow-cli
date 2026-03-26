package scheduler

import (
	"context"
	"time"

	"github.com/oluoyefeso/termiflow-cli/internal/db"
	"github.com/oluoyefeso/termiflow-cli/internal/intelligence"
	"github.com/oluoyefeso/termiflow-cli/internal/providers/llm"
	"github.com/oluoyefeso/termiflow-cli/internal/providers/search"
	"github.com/oluoyefeso/termiflow-cli/pkg/models"
)

type Scheduler struct {
	llmProvider    llm.Provider
	searchProvider search.Provider
	curator        *intelligence.Curator
}

func New(llmProvider llm.Provider, searchProvider search.Provider) *Scheduler {
	return &Scheduler{
		llmProvider:    llmProvider,
		searchProvider: searchProvider,
		curator:        intelligence.NewCurator(llmProvider),
	}
}

// RefreshSubscription fetches and processes new items for a subscription.
func (s *Scheduler) RefreshSubscription(ctx context.Context, sub *models.Subscription) ([]*models.FeedItem, error) {
	var allResults []search.SearchResult

	if s.searchProvider != nil && s.searchProvider.Available() {
		results, err := s.searchProvider.Search(ctx, search.SearchRequest{
			Query:      sub.Topic,
			MaxResults: 10,
			TimeRange:  sub.GetTimeRange(),
		})
		if err == nil {
			allResults = append(allResults, results...)
		}
	}

	allResults = deduplicateByURL(allResults)

	items, err := s.curator.CurateResults(ctx, sub.Topic, allResults)
	if err != nil {
		return nil, err
	}

	var newItems []*models.FeedItem
	for _, item := range items {
		item.SubscriptionID = sub.ID

		exists, _ := db.ItemExistsByURL(item.SourceURL)
		if !exists {
			if err := db.CreateFeedItem(item); err != nil {
				continue
			}
			newItems = append(newItems, item)
		}
	}

	if err := db.UpdateLastFetched(sub.ID); err != nil {
		return nil, err
	}

	return newItems, nil
}

// RefreshAllSubscriptions refreshes all active subscriptions.
func (s *Scheduler) RefreshAllSubscriptions(ctx context.Context) error {
	subs, err := db.GetActiveSubscriptions()
	if err != nil {
		return err
	}

	for _, sub := range subs {
		if shouldRefresh(sub) {
			_, _ = s.RefreshSubscription(ctx, sub)
		}
	}

	return nil
}

func shouldRefresh(sub *models.Subscription) bool {
	if sub.LastFetchedAt == nil {
		return true
	}

	now := time.Now()
	switch sub.Frequency {
	case "hourly":
		return now.Sub(*sub.LastFetchedAt) >= time.Hour
	case "daily":
		return now.Sub(*sub.LastFetchedAt) >= 24*time.Hour
	case "weekly":
		return now.Sub(*sub.LastFetchedAt) >= 7*24*time.Hour
	default:
		return now.Sub(*sub.LastFetchedAt) >= 24*time.Hour
	}
}

func deduplicateByURL(results []search.SearchResult) []search.SearchResult {
	seen := make(map[string]bool)
	var unique []search.SearchResult

	for _, r := range results {
		if !seen[r.URL] {
			seen[r.URL] = true
			unique = append(unique, r)
		}
	}

	return unique
}
