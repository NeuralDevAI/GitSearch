package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gitsearch",
	Short: "A CLI tool to search GitHub code and issues for AI coding agents",
	Long: `gitsearch is a high-performance CLI utility that fetches production-ready
code examples and bug-resolution threads directly from GitHub.

It outputs clean Markdown that can be easily parsed by AI coding agents
like Cursor, Claude Code, or Windsurf.

Requires GITHUB_TOKEN environment variable to be set.`,
	Version: "1.0.0",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
