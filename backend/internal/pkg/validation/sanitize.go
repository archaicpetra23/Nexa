package validation

import (
	"html"
	"strings"
)

func SanitizeString(s string) string {
	return strings.TrimSpace(html.EscapeString(s))
}

func SanitizeSearchQuery(q string) string {
	escaped := strings.ReplaceAll(q, "\\", "\\\\")
	escaped = strings.ReplaceAll(escaped, "%", "\\%")
	escaped = strings.ReplaceAll(escaped, "_", "\\_")
	return strings.TrimSpace(escaped)
}

func SanitizeMultiple(inputs ...string) []string {
	result := make([]string, len(inputs))
	for i, input := range inputs {
		result[i] = SanitizeString(input)
	}
	return result
}
