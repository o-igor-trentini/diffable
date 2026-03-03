package openai

import (
	"strings"

	tiktoken "github.com/pkoukk/tiktoken-go"
)

func CountTokens(text, model string) int {
	enc, err := tiktoken.EncodingForModel(model)
	if err != nil {
		enc, _ = tiktoken.GetEncoding("cl100k_base")
	}
	return len(enc.Encode(text, nil, nil))
}

var filteredPatterns = []string{
	"_test.go",
	".spec.ts",
	".spec.tsx",
	".test.ts",
	".test.tsx",
	"__snapshots__",
	"package-lock.json",
	"go.sum",
	"yarn.lock",
	"pnpm-lock.yaml",
	"vendor/",
	"dist/",
	"build/",
	"node_modules/",
}

func PreprocessDiff(rawDiff string) string {
	return PreprocessDiffForLevel(rawDiff, "functional")
}

func PreprocessDiffForLevel(rawDiff, level string) string {
	if rawDiff == "" {
		return ""
	}

	sections := splitDiffSections(rawDiff)
	var kept []string

	for _, section := range sections {
		if shouldFilterSection(section) {
			continue
		}
		if isBinarySection(section) {
			continue
		}
		if level == "qa_detailed" {
			kept = append(kept, reduceContextWithLines(section, 5))
		} else {
			kept = append(kept, reduceContext(section))
		}
	}

	return strings.Join(kept, "\n")
}

func ChunkDiff(diff string, maxTokens int) []string {
	if diff == "" {
		return nil
	}

	sections := splitDiffSections(diff)
	if len(sections) == 0 {
		return []string{diff}
	}

	var chunks []string
	var current strings.Builder
	currentTokens := 0

	for _, section := range sections {
		sectionTokens := CountTokens(section, "gpt-4o-mini")

		if sectionTokens > maxTokens {
			if current.Len() > 0 {
				chunks = append(chunks, current.String())
				current.Reset()
				currentTokens = 0
			}
			chunks = append(chunks, section)
			continue
		}

		if currentTokens+sectionTokens > maxTokens && current.Len() > 0 {
			chunks = append(chunks, current.String())
			current.Reset()
			currentTokens = 0
		}

		if current.Len() > 0 {
			current.WriteString("\n")
		}
		current.WriteString(section)
		currentTokens += sectionTokens
	}

	if current.Len() > 0 {
		chunks = append(chunks, current.String())
	}

	return chunks
}

func splitDiffSections(diff string) []string {
	lines := strings.Split(diff, "\n")
	var sections []string
	var current strings.Builder

	for _, line := range lines {
		if strings.HasPrefix(line, "diff --git") && current.Len() > 0 {
			sections = append(sections, current.String())
			current.Reset()
		}
		if current.Len() > 0 {
			current.WriteString("\n")
		}
		current.WriteString(line)
	}

	if current.Len() > 0 {
		sections = append(sections, current.String())
	}

	return sections
}

func shouldFilterSection(section string) bool {
	firstLine := section
	if idx := strings.Index(section, "\n"); idx != -1 {
		firstLine = section[:idx]
	}

	lower := strings.ToLower(firstLine)
	for _, pattern := range filteredPatterns {
		if strings.Contains(lower, strings.ToLower(pattern)) {
			return true
		}
	}
	return false
}

func isBinarySection(section string) bool {
	return strings.Contains(section, "Binary files") && strings.Contains(section, "differ")
}

func reduceContextWithLines(section string, contextLines int) string {
	lines := strings.Split(section, "\n")

	changeIndices := make(map[int]bool)
	for i, line := range lines {
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-") {
			if !strings.HasPrefix(line, "---") && !strings.HasPrefix(line, "+++") {
				changeIndices[i] = true
			}
		}
	}

	keepIndices := make(map[int]bool)
	for i, line := range lines {
		if strings.HasPrefix(line, "diff --git") ||
			strings.HasPrefix(line, "---") ||
			strings.HasPrefix(line, "+++") ||
			strings.HasPrefix(line, "@@") {
			keepIndices[i] = true
			continue
		}

		if changeIndices[i] {
			keepIndices[i] = true
			continue
		}
	}

	for idx := range changeIndices {
		for delta := 1; delta <= contextLines; delta++ {
			before := idx - delta
			after := idx + delta

			if before >= 0 && !changeIndices[before] {
				keepIndices[before] = true
			}
			if after < len(lines) && !changeIndices[after] {
				keepIndices[after] = true
			}
		}
	}

	var result []string
	for i, line := range lines {
		if keepIndices[i] {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

func reduceContext(section string) string {
	lines := strings.Split(section, "\n")
	var result []string

	for _, line := range lines {
		if line == "" {
			continue
		}
		switch {
		case strings.HasPrefix(line, "diff --git"):
			result = append(result, line)
		case strings.HasPrefix(line, "---"):
			result = append(result, line)
		case strings.HasPrefix(line, "+++"):
			result = append(result, line)
		case strings.HasPrefix(line, "@@"):
			result = append(result, line)
		case strings.HasPrefix(line, "+"):
			result = append(result, line)
		case strings.HasPrefix(line, "-"):
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}
