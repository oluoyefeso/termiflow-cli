package main

import (
	"os"

	"github.com/oluoyefeso/termiflow-cli/internal/cli"
	"github.com/oluoyefeso/termiflow-cli/internal/providers/llm"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func main() {
	cli.SetVersionInfo(Version, Commit, Date)
	llm.CLIVersion = Version
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
