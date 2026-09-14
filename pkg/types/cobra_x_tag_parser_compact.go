package types

import (
	"reflect"
	"strings"
)

type compactTagParser struct{}

// NewCompactTagParser Creates a new CobraXTagParser that interprets the compact cobra-x DSL:
// cobra-x:"name-expr, key=value, ..."
func NewCompactTagParser() CobraXTagParser {
	return &compactTagParser{}
}

// ParseField parses the value of the cobra-x tag in the compact DSL format.
func (p *compactTagParser) ParseField(field reflect.StructField) (string, map[string]string) {
	return parseCobraX(field.Tag.Get(CobraXTagKey))
}

// parseCobraX parses the value of a cobra-x tag.
// It returns the first part as a name expression and the rest as attributes.
func parseCobraX(tagValue string) (string, map[string]string) {
	attributes := make(map[string]string)
	parts := splitTagValue(tagValue)

	nameExpr := ""
	if len(parts) > 0 {
		nameExpr = strings.TrimSpace(parts[0])
		for i := 1; i < len(parts); i++ {
			if key, val, ok := parseAttribute(parts[i]); ok {
				attributes[key] = val
			}
		}
	}
	return nameExpr, attributes
}

// splitTagValue splits a tag value by commas, respecting single and double quotes.
func splitTagValue(tagValue string) []string {
	var parts []string
	var current strings.Builder
	inQuotes := false
	var quoteChar rune

	for _, r := range tagValue {
		if (r == '\'' || r == '"') && (current.Len() == 0 || current.String()[current.Len()-1] != '\\') {
			if !inQuotes {
				inQuotes = true
				quoteChar = r
			} else if r == quoteChar {
				inQuotes = false
			}
		}

		if r == ',' && !inQuotes {
			parts = append(parts, current.String())
			current.Reset()
		} else {
			current.WriteRune(r)
		}
	}
	parts = append(parts, current.String())
	return parts
}

// parseAttribute parses a key=value pair, stripping quotes from the value if present.
func parseAttribute(part string) (string, string, bool) {
	part = strings.TrimSpace(part)
	kv := strings.SplitN(part, "=", 2)
	if len(kv) != 2 {
		return "", "", false
	}
	key := strings.TrimSpace(kv[0])
	val := strings.TrimSpace(kv[1])
	if len(val) >= 2 && ((val[0] == '\'' && val[len(val)-1] == '\'') || (val[0] == '"' && val[len(val)-1] == '"')) {
		val = val[1 : len(val)-1]
	}
	return key, val, true
}
