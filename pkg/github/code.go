package github

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

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

type codeJob struct {
	index int
	item  *github.CodeResult
}

type jobResult struct {
	index int
	res   *CodeResult
	err   error
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

	numItems := len(result.CodeResults)
	// Process up to 50 items to keep rate limit usage reasonable
	if numItems > 50 {
		numItems = 50
	}

	var repoCache sync.Map // key: owner/repoName, value: *github.Repository
	jobs := make(chan codeJob, numItems)
	resultsChan := make(chan jobResult, numItems)

	workerCtx, cancelWorkers := context.WithCancel(ctx)
	defer cancelWorkers()

	numWorkers := 5
	if numWorkers > numItems {
		numWorkers = numItems
	}

	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-workerCtx.Done():
					return
				case job, ok := <-jobs:
					if !ok {
						return
					}

					owner := job.item.Repository.GetOwner().GetLogin()
					repoName := job.item.Repository.GetName()
					repoKey := fmt.Sprintf("%s/%s", owner, repoName)

					var fullRepo *github.Repository
					if val, ok := repoCache.Load(repoKey); ok {
						fullRepo = val.(*github.Repository)
					} else {
						var err error
						fullRepo, _, err = client.Repositories.Get(workerCtx, owner, repoName)
						if err != nil {
							resultsChan <- jobResult{index: job.index, err: err}
							continue
						}
						repoCache.Store(repoKey, fullRepo)
					}

					stars := fullRepo.GetStargazersCount()
					if stars < minStars {
						resultsChan <- jobResult{index: job.index} // Filtered out
						continue
					}

					path := job.item.GetPath()
					fileContent, _, _, err := client.Repositories.GetContents(workerCtx, owner, repoName, path, nil)
					if err != nil {
						resultsChan <- jobResult{index: job.index, err: err}
						continue
					}

					content, err := fileContent.GetContent()
					if err != nil {
						resultsChan <- jobResult{index: job.index, err: err}
						continue
					}

					language := detectLanguage(path)
					resultsChan <- jobResult{
						index: job.index,
						res: &CodeResult{
							Owner:     owner,
							Repo:      repoName,
							Path:      path,
							Stars:     stars,
							UpdatedAt: fullRepo.GetUpdatedAt().Format("2006-01-02"),
							URL:       fileContent.GetHTMLURL(),
							Content:   content,
							Language:  language,
						},
					}
				}
			}
		}()
	}

	// Send all jobs
	for i := 0; i < numItems; i++ {
		jobs <- codeJob{
			index: i,
			item:  result.CodeResults[i],
		}
	}
	close(jobs)

	// Monitor workers completion
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	collected := make([]*CodeResult, numItems)
	finished := make([]bool, numItems)

	for jr := range resultsChan {
		finished[jr.index] = true
		if jr.err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to fetch details for %d: %v\n", jr.index, jr.err)
			continue
		}
		if jr.res != nil {
			collected[jr.index] = jr.res
		}

		// Check if we can stop early while preserving the top ranked results
		validCount := 0
		canStop := false
		for idx, item := range collected {
			if !finished[idx] {
				break // Must wait for this higher-ranked result to finish
			}
			if item != nil {
				validCount++
				if validCount == limit {
					canStop = true
					break
				}
			}
		}

		if canStop {
			cancelWorkers()
			break
		}
	}

	// Filter out nil elements and build final ordered results
	var results []CodeResult
	for _, item := range collected {
		if item != nil {
			results = append(results, *item)
			if len(results) == limit {
				break
			}
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
