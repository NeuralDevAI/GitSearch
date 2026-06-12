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
	codeQuery    string
	codeLimit    int
	codeMinStars int
	codeSnippet  bool
	codeContext  int
)

var codeCmd = &cobra.Command{
	Use:   "code",
	Short: "Search GitHub for code files",
	Long: `Search GitHub repositories for code files matching specific keywords.
Returns raw file contents or hybrid snippets in structured Markdown format.

Example:
  gitsearch code -q "language:typescript stripe payment" -limit 2 -min-stars 100 --snippet --context 10`,
	RunE: runCodeSearch,
}

func init() {
	rootCmd.AddCommand(codeCmd)

	codeCmd.Flags().StringVarP(&codeQuery, "query", "q", "", "GitHub search query (required)")
	codeCmd.Flags().IntVarP(&codeLimit, "limit", "l", 3, "Maximum number of files to fetch")
	codeCmd.Flags().IntVar(&codeMinStars, "min-stars", 50, "Minimum repository stars")
	codeCmd.Flags().BoolVar(&codeSnippet, "snippet", false, "Extract relevant code snippets instead of full files")
	codeCmd.Flags().IntVar(&codeContext, "context", 10, "Number of context lines around snippet matches")

	codeCmd.MarkFlagRequired("query")
}

func runCodeSearch(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := github.NewClient(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	results, err := github.SearchCode(ctx, client, codeQuery, codeLimit, codeMinStars, codeSnippet, codeContext)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return err
	}

	if len(results) == 0 {
		fmt.Println("No results found.")
		return nil
	}

	for _, result := range results {
		fmt.Printf("# SOURCE: %s/%s/%s\n", result.Owner, result.Repo, result.Path)
		fmt.Printf("# STARS: %d | UPDATED: %s\n", result.Stars, result.UpdatedAt)
		fmt.Printf("# URL: %s\n\n", result.URL)

		if result.Language != "" {
			fmt.Printf("```%s\n", result.Language)
		} else {
			fmt.Println("```")
		}
		fmt.Println(result.Content)
		fmt.Println("```")
		fmt.Println("---")
		fmt.Println()
	}

	return nil
}
