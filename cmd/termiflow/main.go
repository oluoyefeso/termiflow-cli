package main

import (
	"os"

	"github.com/oluoyefeso/termiflow-cli/internal/cli"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func main() {
	cli.SetVersionInfo(Version, Commit, Date)
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
