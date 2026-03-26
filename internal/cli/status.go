package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/oluoyefeso/termiflow-cli/internal/config"
	"github.com/oluoyefeso/termiflow-cli/internal/db"
	"github.com/oluoyefeso/termiflow-cli/internal/ui"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current configuration summary and database stats",
	RunE:  runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) error {
	cfg := config.Get()

	fmt.Println(ui.Header("termiflow status"))
	fmt.Println()

	// API key (masked)
	apiKey := cfg.Providers.Managed.APIKey
	if apiKey != "" {
		masked := apiKey[:6] + "..." + apiKey[len(apiKey)-4:]
		fmt.Print(ui.Info("API Key", masked))
	} else {
		fmt.Print(ui.Info("API Key", ui.ErrorStyle.Render("not configured")))
	}

	// Base URL
	baseURL := cfg.Providers.Managed.BaseURL
	if baseURL == "" {
		baseURL = "https://api.termiflow.com"
	}
	fmt.Print(ui.Info("Base URL", baseURL))

	// Provider
	fmt.Print(ui.Info("Provider", cfg.General.DefaultProvider))
	fmt.Println()

	// Subscriptions
	subs, err := db.GetActiveSubscriptions()
	if err == nil {
		fmt.Print(ui.Info("Subscriptions", fmt.Sprintf("%d active", len(subs))))
		for _, sub := range subs {
			lastFetched := "never"
			if sub.LastFetchedAt != nil {
				lastFetched = sub.LastFetchedAt.Format("2006-01-02 15:04")
			}
			fmt.Printf("    %s  (last: %s)\n", ui.TitleStyle.Render(sub.Topic), ui.MutedStyle.Render(lastFetched))
		}
	}
	fmt.Println()

	// Database size
	dbPath := filepath.Join(config.GetDataDir(), "termiflow.db")
	if info, err := os.Stat(dbPath); err == nil {
		sizeKB := info.Size() / 1024
		if sizeKB < 1024 {
			fmt.Print(ui.Info("Database", fmt.Sprintf("%d KB (%s)", sizeKB, dbPath)))
		} else {
			fmt.Print(ui.Info("Database", fmt.Sprintf("%.1f MB (%s)", float64(sizeKB)/1024, dbPath)))
		}
	}

	// Config path
	fmt.Print(ui.Info("Config", config.GetConfigPath()))
	fmt.Println()

	return nil
}
