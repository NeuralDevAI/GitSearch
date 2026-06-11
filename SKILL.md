# gitsearch - GitHub Code & Issue Search Tool

AI agent skill for fetching production-ready code examples and bug solutions from GitHub.

## Tool Overview

`gitsearch` is a CLI utility that searches GitHub and returns clean Markdown output with:
- **Code search**: Raw file contents from high-quality repositories
- **Issue search**: Problem descriptions + top-voted solution comments

## Prerequisites

- `gitsearch` binary must be in PATH or current directory
- `GITHUB_TOKEN` environment variable must be set
- If token missing, instruct user to set it: `export GITHUB_TOKEN=your_token`

## When to Use gitsearch

### Code Search Triggers

Invoke `gitsearch code` when:

1. **Unknown Implementation Patterns**
   - User asks: "How do I implement X in Y framework?"
   - You lack confidence in the implementation approach
   - Need real-world examples of a specific library/API usage

2. **Framework-Specific Features**
   - Authentication flows (OAuth, JWT, session management)
   - Payment integrations (Stripe, PayPal)
   - Real-time features (WebSockets, SSE, WebRTC)
   - Complex UI components (drag-and-drop, virtualized lists)
   - Database patterns (migrations, transactions, ORMs)

3. **Advanced/Niche Techniques**
   - Performance optimizations
   - Security implementations
   - Complex CSS animations or layouts
   - Advanced React patterns (Suspense, concurrent features)
   - GraphQL schema design or resolvers

4. **Language-Specific Idioms**
   - Go concurrency patterns
   - Rust lifetime management
   - TypeScript advanced types
   - Python async/await patterns

### Issue Search Triggers

Invoke `gitsearch issues` when:

1. **Error Messages**
   - User reports a specific error message
   - Stack trace with unclear root cause
   - Build/compilation failures

2. **Known Bug Patterns**
   - "Why does X fail when Y?"
   - "How to fix [specific error]?"
   - Common framework bugs (CORS, port conflicts, dependency issues)

3. **Configuration Problems**
   - Docker/K8s setup issues
   - Build tool configuration (webpack, vite, rollup)
   - CI/CD pipeline failures

## Query Construction Guidelines

### Code Search Queries

**Template**: `language:<lang> [library/framework] [specific-feature] [context-keywords]`

**Examples:**

```bash
# Bad: Too generic
gitsearch code -q "authentication"

# Good: Specific language, library, and use case
gitsearch code -q "language:typescript nextjs middleware authentication jwt" --limit 2 --min-stars 200

# Bad: Missing language qualifier
gitsearch code -q "stripe payment"

# Good: Language + library + feature
gitsearch code -q "language:javascript stripe payment intent webhook" --min-stars 100
```

**Best Practices:**
- Always include `language:<lang>` to filter by programming language
- Add framework/library name (e.g., `react`, `fastapi`, `django`, `express`)
- Include 2-4 specific keywords describing the feature
- Use `--min-stars` to ensure production quality (50-200 depending on need)
- Keep `--limit` low (1-3) to avoid token waste

**Common Language Tags:**
- `language:typescript`, `language:javascript`, `language:python`
- `language:go`, `language:java`, `language:rust`
- `language:c++`, `language:csharp`, `language:swift`

### Issue Search Queries

**Template**: `"<error-message>" [technology] [context]`

**Examples:**

```bash
# Error message search
gitsearch issues -q "EADDRINUSE" --state closed --limit 2

# Framework-specific error
gitsearch issues -q "react Cannot read property of undefined setState" --state closed

# Build error
gitsearch issues -q "webpack Module not found" --state closed --limit 3
```

**Best Practices:**
- Quote exact error messages for precision
- Use `--state closed` (default) to find resolved issues
- Include framework/tool name for context
- Keep `--limit` at 2-3 for most relevant results

## Execution Pattern

### 1. Determine If gitsearch Is Needed

```
IF user asks for implementation AND (
   you lack confidence OR
   it's a niche/advanced feature OR
   it involves external APIs/libraries
) THEN use gitsearch code

IF user reports error message OR asks "how to fix X"
THEN use gitsearch issues
```

### 2. Compile the Query

Extract key terms from user's request:
- Programming language (from context or explicit)
- Framework/library mentioned
- Core feature/functionality
- Specific keywords

### 3. Execute Search

```bash
# Code search example
gitsearch code -q "language:go chi router middleware auth" --limit 2 --min-stars 100

# Issue search example
gitsearch issues -q "prisma migration failed column" --state closed --limit 2
```

### 4. Parse and Apply Results

**For Code Results:**
- Read the returned code blocks
- Identify the relevant pattern/implementation
- Adapt to user's specific context (don't copy-paste blindly)
- Explain what the code does and why it works
- Credit the source repository if you reference it directly

**For Issue Results:**
- Read the issue description to confirm it matches user's problem
- Extract the solution from "Top Solution" comment
- Translate solution to user's context
- Provide step-by-step instructions

### 5. Communicate to User

**Template:**
```
I found a production example from [repo-name] (⭐ stars) that shows how to [feature].

Here's the relevant approach:
[Your adaptation of the code/solution]

[Explanation of how it works]
```

## Token Efficiency Guidelines

1. **Limit Results**: Default to `--limit 1` or `--limit 2` unless you need to compare approaches
2. **Filter by Stars**: Use `--min-stars 100-200` for popular frameworks, `--min-stars 50` for niche libraries
3. **Precise Queries**: Specific queries = fewer, better results = less token usage
4. **Single Search**: If first result is good enough, don't make additional searches

## Error Handling

### GITHUB_TOKEN Missing
```
Error: GITHUB_TOKEN environment variable is not set.
```
**Response**: "You need to set your GitHub token. Run: `export GITHUB_TOKEN=your_token`. Get one at https://github.com/settings/tokens"

### Rate Limit Exceeded
```
Error: GitHub API rate limit exceeded. Resets at: 2026-06-10T06:30:00Z
```
**Response**: "GitHub rate limit hit. Resets at [time]. I'll proceed with my existing knowledge for now."

### No Results Found
```
No results found.
```
**Response**: Reformulate query with broader keywords or fewer constraints, OR proceed without search results.

## Advanced Query Techniques

### Path-Specific Searches
```bash
# Search only in specific directories
gitsearch code -q "language:typescript path:src/auth jwt" --limit 2

# Search by file extension
gitsearch code -q "language:python extension:py fastapi crud" --min-stars 150
```

### Repository-Specific Searches
```bash
# Not directly supported, but can filter by stars and language
gitsearch code -q "language:go stars:>500 grpc microservice" --limit 2
```

### Complex Boolean Queries
```bash
# Multiple keywords (implicit AND)
gitsearch code -q "language:rust async tokio http client" --limit 2

# Exact phrases (use quotes in the query string)
gitsearch code -q 'language:javascript "useEffect cleanup"' --limit 2
```

## Real-World Examples

### Example 1: User Needs OAuth Implementation

**User**: "How do I implement Google OAuth in my Next.js app?"

**Agent Decision**: Use gitsearch (specific external API + framework)

**Query**:
```bash
gitsearch code -q "language:typescript nextjs google oauth authentication" --limit 2 --min-stars 100
```

**After Results**: Adapt the pattern found, explain the flow, provide implementation.

### Example 2: User Has Docker Error

**User**: "I'm getting 'Bind for 0.0.0.0:3000 failed: port is already allocated'"

**Agent Decision**: Use gitsearch issues (specific error message)

**Query**:
```bash
gitsearch issues -q "docker port is already allocated" --state closed --limit 2
```

**After Results**: Extract solution (usually `docker ps`, kill process, or change port), provide commands.

### Example 3: User Asks Generic Question

**User**: "How do I center a div?"

**Agent Decision**: DON'T use gitsearch (trivial CSS, you have sufficient knowledge)

**Response**: Directly provide CSS solution.

### Example 4: Advanced React Pattern

**User**: "How do I implement infinite scroll with React Query and virtual scrolling?"

**Agent Decision**: Use gitsearch (advanced/niche combination)

**Query**:
```bash
gitsearch code -q "language:typescript react-query infinite scroll virtual" --limit 2 --min-stars 150
```

**After Results**: Show the pattern, explain virtualization + data fetching coordination.

## Integration Checklist

Before executing gitsearch:
- [ ] Is this a case where real production code would help?
- [ ] Have I formulated a specific query with language and key terms?
- [ ] Have I set appropriate `--limit` and `--min-stars`?
- [ ] Am I prepared to adapt the results (not just copy-paste)?

After receiving results:
- [ ] Did I find relevant code/solution?
- [ ] Have I adapted it to user's context?
- [ ] Have I explained what it does and why?
- [ ] Have I avoided token waste by keeping only relevant parts?

## Summary

**Use gitsearch when**: You need production examples or bug solutions from real repositories.

**Don't use when**: You have sufficient knowledge, it's trivial, or it's too generic to search effectively.

**Key principle**: Precise queries + low limits + careful adaptation = efficient, high-quality assistance.
