package utils

import (
	"regexp"
	"strings"
)

// ExtractCommandUse extracts the command use string from given command name by matching different naming conventions.
func ExtractCommandUse(s string) string {

	regexes := []*regexp.Regexp{
		regexp.MustCompile(`^(.+)Command$`),     // CamelCase convention
		regexp.MustCompile(`^([^_]+)_command$`), // SnakeCase convention
		regexp.MustCompile(`^([^_]+)$`),         // Single-word command
	}

	for _, regex := range regexes {
		match := regex.FindStringSubmatch(s)
		if len(match) > 1 {
			return strings.ToLower(match[1])
		}
	}

	return s
}
