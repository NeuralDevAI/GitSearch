package github

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v62/github"
)

type CodeResult struct {
	Owner      string
	Repo       string
	Path       string
	Stars      int
	UpdatedAt  string
	URL        string
	Content    string
	Language   string
}

func SearchCode(ctx context.Context, client *github.Client, query string, limit, minStars int) ([]CodeResult, error) {
	if err := CheckRateLimit(ctx, client); err != nil {
		return nil, err
	}

	// Code Search API doesn't support stars filter directly
	// We'll fetch more results and filter by stars ourselves
	opts := &github.SearchOptions{
		Sort:  "indexed",
		Order: "desc",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	result, _, err := client.Search.Code(ctx, query, opts)
	if err != nil {
		return nil, fmt.Errorf("code search failed: %w", err)
	}

	if result.GetTotal() == 0 {
		return []CodeResult{}, nil
	}

	var results []CodeResult
	repoCache := make(map[string]*github.Repository)
	checked := 0
	for _, item := range result.CodeResults {
		if len(results) >= limit {
			break
		}

		checked++
		owner := item.Repository.GetOwner().GetLogin()
		repoName := item.Repository.GetName()
		repoKey := fmt.Sprintf("%s/%s", owner, repoName)

		// Check cache first
		fullRepo, ok := repoCache[repoKey]
		if !ok {
			// Fetch full repository info to get accurate star count
			var err error
			fullRepo, _, err = client.Repositories.Get(ctx, owner, repoName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to fetch repo info for %s: %v\n", repoKey, err)
				continue
			}
			repoCache[repoKey] = fullRepo
		}

		stars := fullRepo.GetStargazersCount()

		if stars < minStars {
			continue
		}

		path := item.GetPath()

		fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repoName, path, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to fetch content for %s/%s/%s: %v\n", owner, repoName, path, err)
			continue
		}

		// Use DownloadContents which handles decoding automatically
		content, err := fileContent.GetContent()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to get content for %s/%s/%s: %v\n", owner, repoName, path, err)
			continue
		}

		language := detectLanguage(path)

		results = append(results, CodeResult{
			Owner:      owner,
			Repo:       repoName,
			Path:       path,
			Stars:      stars,
			UpdatedAt:  fullRepo.GetUpdatedAt().Format("2006-01-02"),
			URL:        fileContent.GetHTMLURL(),
			Content:    content,
			Language:   language,
		})

		if len(results) >= limit {
			break
		}
	}

	return results, nil
}

func detectLanguage(path string) string {
	ext := strings.ToLower(filepath.Ext(path))

	langMap := map[string]string{
		".go":   "go",
		".js":   "javascript",
		".ts":   "typescript",
		".jsx":  "jsx",
		".tsx":  "tsx",
		".py":   "python",
		".java": "java",
		".c":    "c",
		".cpp":  "cpp",
		".cs":   "csharp",
		".rb":   "ruby",
		".php":  "php",
		".rs":   "rust",
		".kt":   "kotlin",
		".swift": "swift",
		".sql":  "sql",
		".sh":   "bash",
		".yaml": "yaml",
		".yml":  "yaml",
		".json": "json",
		".xml":  "xml",
		".html": "html",
		".css":  "css",
		".scss": "scss",
		".md":   "markdown",
	}

	if lang, ok := langMap[ext]; ok {
		return lang
	}
	return ""
}
