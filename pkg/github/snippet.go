package github

import (
	"fmt"
	"strings"
)

// ExtractKeywords parses a GitHub query string and extracts raw search keywords
// ignoring search qualifiers like "language:go" or "stars:>50".
func ExtractKeywords(query string) []string {
	words := strings.Fields(query)
	var keywords []string
	for _, w := range words {
		if strings.Contains(w, ":") {
			continue
		}
		w = strings.Trim(w, `"'`)
		if len(w) > 1 {
			keywords = append(keywords, strings.ToLower(w))
		}
	}
	return keywords
}

// getCommentPrefix returns the appropriate single-line comment prefix for the language.
func getCommentPrefix(language string) string {
	lang := strings.ToLower(language)
	switch lang {
	case "python", "ruby", "bash", "shell", "yaml", "yml", "dockerfile":
		return "#"
	case "sql", "lua":
		return "--"
	default:
		return "//"
	}
}

// findHeaderEnd returns the 0-indexed line number up to which the file header (packages, imports) extends.
// It scans the first 50 lines.
func findHeaderEnd(lines []string) int {
	lastHeaderIdx := -1
	headerKeywords := []string{"package ", "import ", "import (", "require(", "#include", "using ", "from ", "import\t"}
	limit := 50
	if len(lines) < limit {
		limit = len(lines)
	}

	for i := 0; i < limit; i++ {
		trimmed := strings.TrimSpace(lines[i])
		isHeader := false
		for _, kw := range headerKeywords {
			if strings.HasPrefix(trimmed, kw) {
				isHeader = true
				break
			}
		}
		if isHeader {
			lastHeaderIdx = i
		}
	}

	// Handle block imports closing parenthesis
	if lastHeaderIdx != -1 && lastHeaderIdx < len(lines)-1 {
		// If the last header line starts a block, try to find the closing paren within 15 lines
		if strings.Contains(lines[lastHeaderIdx], "(") && !strings.Contains(lines[lastHeaderIdx], ")") {
			for j := lastHeaderIdx + 1; j < len(lines) && j < lastHeaderIdx+15; j++ {
				if strings.Contains(lines[j], ")") {
					lastHeaderIdx = j
					break
				}
			}
		}
	}

	return lastHeaderIdx
}

type rangeWindow struct {
	start int
	end   int
}

// ExtractSnippet pulls relevant parts of code matching the query with a context window.
func ExtractSnippet(content string, query string, language string, contextLines int) string {
	keywords := ExtractKeywords(query)
	if len(keywords) == 0 {
		return content // Nothing to match, return full content
	}

	// Normalize CRLF to LF
	content = strings.ReplaceAll(content, "\r\n", "\n")
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return ""
	}

	commentPrefix := getCommentPrefix(language)
	headerEnd := findHeaderEnd(lines)

	// Identify keyword matching lines
	var matchedLines []int
	for i := headerEnd + 1; i < len(lines); i++ {
		lineLower := strings.ToLower(lines[i])
		matched := false
		for _, kw := range keywords {
			if strings.Contains(lineLower, kw) {
				matched = true
				break
			}
		}
		if matched {
			matchedLines = append(matchedLines, i)
		}
	}

	// If no match is found, return the header plus the first 20 lines of code
	if len(matchedLines) == 0 {
		endIdx := headerEnd + 1 + 20
		if endIdx > len(lines) {
			endIdx = len(lines)
		}
		var builder strings.Builder
		for i := 0; i < endIdx; i++ {
			builder.WriteString(lines[i] + "\n")
		}
		if endIdx < len(lines) {
			builder.WriteString(fmt.Sprintf("%s ... [Lines %d-%d omitted] ...\n", commentPrefix, endIdx+1, len(lines)))
		}
		return builder.String()
	}

	// Create ranges around matched lines
	var rawRanges []rangeWindow
	for _, mLine := range matchedLines {
		start := mLine - contextLines
		if start < headerEnd+1 {
			start = headerEnd + 1
		}
		end := mLine + contextLines
		if end >= len(lines) {
			end = len(lines) - 1
		}
		rawRanges = append(rawRanges, rangeWindow{start: start, end: end})
	}

	// Merge overlapping ranges
	var mergedRanges []rangeWindow
	if len(rawRanges) > 0 {
		current := rawRanges[0]
		for i := 1; i < len(rawRanges); i++ {
			next := rawRanges[i]
			// If ranges overlap or touch, merge them
			if next.start <= current.end+1 {
				if next.end > current.end {
					current.end = next.end
				}
			} else {
				mergedRanges = append(mergedRanges, current)
				current = next
			}
		}
		mergedRanges = append(mergedRanges, current)
	}

	// Construct final output
	var builder strings.Builder

	// 1. Output Header (Package/Imports)
	if headerEnd >= 0 {
		for i := 0; i <= headerEnd; i++ {
			builder.WriteString(lines[i] + "\n")
		}
	}

	// 2. Output Ranges with omission markers
	lastPrintedLine := headerEnd

	for _, r := range mergedRanges {
		// Output omission marker if there's a gap
		if r.start > lastPrintedLine+1 {
			builder.WriteString(fmt.Sprintf("%s ... [Lines %d-%d omitted] ...\n", commentPrefix, lastPrintedLine+2, r.start))
		}

		// Output matching lines
		for i := r.start; i <= r.end; i++ {
			builder.WriteString(lines[i] + "\n")
		}
		lastPrintedLine = r.end
	}

	// 3. Output final omission marker if necessary
	if lastPrintedLine < len(lines)-1 {
		builder.WriteString(fmt.Sprintf("%s ... [Lines %d-%d omitted] ...\n", commentPrefix, lastPrintedLine+2, len(lines)))
	}

	return builder.String()
}
