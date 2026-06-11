# gitsearch

A high-performance, zero-LLM Go CLI utility that fetches production-ready code examples and bug-resolution threads directly from GitHub. Designed to be executed by AI coding agents (Cursor, Claude Code, Windsurf) or developers in the terminal.

## Features

- **Pure GitHub API Integration**: Direct REST/GraphQL API calls, no LLM dependencies
- **Clean Markdown Output**: Structured output that's easily parseable by AI agents
- **Fast & Lightweight**: Single binary, minimal dependencies
- **Production-Grade Results**: Filter by repository stars to ensure quality code examples
- **Two Search Modes**:
  - `code`: Search for code files with raw contents
  - `issues`: Search for issues/PRs with top-voted solutions

## Installation

### Prerequisites

- Go 1.21 or higher
- GitHub Personal Access Token

### Build from Source

```bash
git clone <repository-url>
cd gitsearch
go build -o gitsearch
```

### Cross-Platform Builds

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o gitsearch-linux

# macOS
GOOS=darwin GOARCH=amd64 go build -o gitsearch-darwin

# Windows
GOOS=windows GOARCH=amd64 go build -o gitsearch.exe
```

## Setup

### 1. Create GitHub Token

Create a GitHub Personal Access Token at: https://github.com/settings/tokens

**Required Scopes**: `public_repo` (or `repo` for private repositories)

### 2. Set Environment Variable

**Linux/macOS:**
```bash
export GITHUB_TOKEN=your_github_token_here
```

**Windows (PowerShell):**
```powershell
$env:GITHUB_TOKEN="your_github_token_here"
```

**Windows (CMD):**
```cmd
set GITHUB_TOKEN=your_github_token_here
```

## Usage

### Code Search

Search GitHub repositories for code files matching specific keywords.

```bash
gitsearch code -q "language:typescript stripe payment" --limit 3 --min-stars 100
```

**Flags:**
- `-q, --query` (required): GitHub search query
- `--limit` (default: 3): Maximum number of files to fetch
- `--min-stars` (default: 50): Minimum repository stars

**Output Format:**
```markdown
# SOURCE: owner/repo/path/to/file.ts
# STARS: 420 | UPDATED: 2026-05-10
# URL: https://github.com/owner/repo/blob/main/path/to/file.ts

```typescript
[Raw file contents here]
```
---
```

### Issues Search

Search GitHub issues and PRs for error messages or bug solutions.

```bash
gitsearch issues -q "docker EADDRINUSE" --state closed --limit 2
```

**Flags:**
- `-q, --query` (required): GitHub search query
- `-s, --state` (default: "closed"): Issue state (open, closed, or all)
- `--limit` (default: 3): Maximum number of issues to fetch

**Output Format:**
```markdown
# ISSUE: [Title of the Issue] (owner/repo #123)
# STATE: Closed | URL: https://github.com/owner/repo/issues/123

## Description:
[Brief body of the issue]

## Top Solution (Most Reacted Comment):
[Contents of the comment with the highest positive reactions]
---
```

## Query Examples

### Code Search Queries

```bash
# Search for Go chi router implementations
gitsearch code -q "language:go chi router" --limit 2

# Find React hooks examples in popular repos
gitsearch code -q "language:javascript useEffect cleanup" --min-stars 200

# Search for Python FastAPI authentication
gitsearch code -q "language:python fastapi jwt auth" --limit 3
```

### Issue Search Queries

```bash
# Find solutions for Docker port conflicts
gitsearch issues -q "docker EADDRINUSE" --state closed

# Search for React rendering issues
gitsearch issues -q "react 'Cannot read property of undefined'" --state closed

# Find database migration problems
gitsearch issues -q "prisma migration failed" --state closed --limit 5
```

## GitHub Search Query Syntax

gitsearch uses GitHub's search syntax. Common qualifiers:

- `language:go` - Filter by programming language
- `stars:>100` - Repositories with more than 100 stars
- `in:file` - Search file contents
- `path:src/` - Search specific paths
- `extension:ts` - Filter by file extension

Full syntax: https://docs.github.com/en/search-github/searching-on-github/searching-code

## Rate Limits

- **Authenticated**: 5,000 requests/hour (requires GITHUB_TOKEN)
- **Unauthenticated**: 60 requests/hour
- **Code Search**: 30 requests/minute (authenticated)

The tool automatically checks rate limits before executing searches.

## Error Handling

- Missing `GITHUB_TOKEN`: Clear error message with setup instructions
- Rate limit exceeded: Reports reset time
- Network errors: Displays error details to stderr
- Invalid queries: Shows usage help

## Use Cases

### For AI Coding Agents

AI agents can invoke gitsearch to:
1. Find production-ready code examples for implementing features
2. Search for bug solutions and error resolutions
3. Discover best practices from popular repositories
4. Fetch working implementations as reference code

### For Developers

Developers can use gitsearch to:
1. Quickly find code examples without leaving the terminal
2. Research how popular projects solve specific problems
3. Find issue resolutions for error messages
4. Analyze implementation patterns across repositories

## Architecture

- **Language**: Go 1.21+
- **CLI Framework**: cobra
- **GitHub API**: go-github/v62
- **Authentication**: oauth2

## Project Structure

```
gitsearch/
├── main.go              # Entry point
├── cmd/
│   ├── root.go          # Root command
│   ├── code.go          # Code search subcommand
│   └── issues.go        # Issues search subcommand
└── pkg/
    └── github/
        ├── client.go    # GitHub API client
        ├── code.go      # Code search logic
        └── issues.go    # Issue search logic
```

## License

MIT

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.
