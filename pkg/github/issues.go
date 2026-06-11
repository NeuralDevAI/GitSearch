package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v62/github"
)

type IssueResult struct {
	Title       string
	Number      int
	Owner       string
	Repo        string
	State       string
	URL         string
	Body        string
	TopSolution string
}

func SearchIssues(ctx context.Context, client *github.Client, query string, state string, limit int) ([]IssueResult, error) {
	if err := CheckRateLimit(ctx, client); err != nil {
		return nil, err
	}

	fullQuery := query
	if state != "" && state != "all" {
		fullQuery = fmt.Sprintf("%s state:%s", query, state)
	}

	opts := &github.SearchOptions{
		Sort:  "reactions",
		Order: "desc",
		ListOptions: github.ListOptions{
			PerPage: limit,
		},
	}

	result, _, err := client.Search.Issues(ctx, fullQuery, opts)
	if err != nil {
		return nil, fmt.Errorf("issue search failed: %w", err)
	}

	if result.GetTotal() == 0 {
		return []IssueResult{}, nil
	}

	var results []IssueResult
	for i, issue := range result.Issues {
		if i >= limit {
			break
		}

		repoFullName := issue.GetRepository().GetFullName()
		parts := splitRepoFullName(repoFullName)
		if len(parts) != 2 {
			continue
		}
		owner := parts[0]
		repoName := parts[1]

		topSolution := ""
		if issue.GetComments() > 0 {
			topSolution, _ = findTopSolution(ctx, client, owner, repoName, issue.GetNumber())
		}

		results = append(results, IssueResult{
			Title:       issue.GetTitle(),
			Number:      issue.GetNumber(),
			Owner:       owner,
			Repo:        repoName,
			State:       issue.GetState(),
			URL:         issue.GetHTMLURL(),
			Body:        issue.GetBody(),
			TopSolution: topSolution,
		})

		if len(results) >= limit {
			break
		}
	}

	return results, nil
}

func findTopSolution(ctx context.Context, client *github.Client, owner, repo string, issueNumber int) (string, error) {
	opts := &github.IssueListCommentsOptions{
		Sort:      github.String("created"),
		Direction: github.String("asc"),
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	var allComments []*github.IssueComment
	for {
		comments, resp, err := client.Issues.ListComments(ctx, owner, repo, issueNumber, opts)
		if err != nil {
			return "", fmt.Errorf("failed to list comments: %w", err)
		}
		allComments = append(allComments, comments...)

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	if len(allComments) == 0 {
		return "", nil
	}

	var topComment *github.IssueComment
	maxReactions := 0

	for _, comment := range allComments {
		reactions := comment.GetReactions()
		totalPositive := reactions.GetPlusOne() + reactions.GetHeart() + reactions.GetHooray() + reactions.GetRocket()

		if totalPositive > maxReactions {
			maxReactions = totalPositive
			topComment = comment
		}
	}

	if topComment != nil {
		return topComment.GetBody(), nil
	}

	return allComments[0].GetBody(), nil
}

func splitRepoFullName(fullName string) []string {
	parts := make([]string, 0, 2)
	for i, part := range []rune(fullName) {
		if part == '/' {
			parts = append(parts, fullName[:i])
			parts = append(parts, fullName[i+1:])
			break
		}
	}
	return parts
}
