package reflection

import (
	"strings"
)

// parseFlagNameExpression parses a flag name expression (e.g., "name|-n" or "--name|-n").
func parseFlagNameExpression(expression string) (name string, shortHand string) {
	parts := strings.Split(expression, "|")
	name = ""
	shortHand = ""
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if longName, ok := strings.CutPrefix(part, "--"); ok {
			name = longName
		} else if shortName, ok := strings.CutPrefix(part, "-"); ok {
			shortHand = shortName
		} else {
			if name == "" {
				name = part
			}
		}
	}
	if name == "" && shortHand != "" {
		name = shortHand
	}
	return
}
