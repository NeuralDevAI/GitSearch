package github

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/v62/github"
	"golang.org/x/oauth2"
)

func NewClient(ctx context.Context) (*github.Client, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN environment variable is not set.\n\nPlease set it with:\n  export GITHUB_TOKEN=your_github_token\n\nYou can create a token at: https://github.com/settings/tokens")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	tc.Timeout = 15 * time.Second

	client := github.NewClient(tc)
	return client, nil
}

func CheckRateLimit(ctx context.Context, client *github.Client) error {
	rateLimits, _, err := client.RateLimits(ctx)
	if err != nil {
		return fmt.Errorf("failed to check rate limits: %w", err)
	}

	if rateLimits.Search.Remaining == 0 {
		resetTime := rateLimits.Search.Reset.Time
		return fmt.Errorf("GitHub API rate limit exceeded. Resets at: %s", resetTime.Format(time.RFC3339))
	}

	return nil
}
