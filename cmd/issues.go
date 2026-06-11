package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"gitsearch/pkg/github"

	"github.com/spf13/cobra"
)

var (
	issuesQuery string
	issuesState string
	issuesLimit int
)

var issuesCmd = &cobra.Command{
	Use:   "issues",
	Short: "Search GitHub for issues and pull requests",
	Long: `Search GitHub issues and PRs for error messages or bugs.
Returns issue details with the top-voted solution comment.

Example:
  gitsearch issues -q "docker EADDRINUSE" -state closed -limit 2`,
	RunE: runIssuesSearch,
}

func init() {
	rootCmd.AddCommand(issuesCmd)

	issuesCmd.Flags().StringVarP(&issuesQuery, "query", "q", "", "GitHub search query (required)")
	issuesCmd.Flags().StringVarP(&issuesState, "state", "s", "closed", "Issue state: open, closed, or all")
	issuesCmd.Flags().IntVarP(&issuesLimit, "limit", "l", 3, "Maximum number of issues to fetch")

	issuesCmd.MarkFlagRequired("query")
}

func runIssuesSearch(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := github.NewClient(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	results, err := github.SearchIssues(ctx, client, issuesQuery, issuesState, issuesLimit)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return err
	}

	if len(results) == 0 {
		fmt.Println("No results found.")
		return nil
	}

	for _, result := range results {
		fmt.Printf("# ISSUE: %s (%s/%s #%d)\n", result.Title, result.Owner, result.Repo, result.Number)
		fmt.Printf("# STATE: %s | URL: %s\n\n", result.State, result.URL)

		fmt.Println("## Description:")
		if result.Body != "" {
			fmt.Println(result.Body)
		} else {
			fmt.Println("(No description provided)")
		}
		fmt.Println()

		if result.TopSolution != "" {
			fmt.Println("## Top Solution (Most Reacted Comment):")
			fmt.Println(result.TopSolution)
		} else {
			fmt.Println("## Top Solution:")
			fmt.Println("(No comments found)")
		}

		fmt.Println("---")
		fmt.Println()
	}

	return nil
}
