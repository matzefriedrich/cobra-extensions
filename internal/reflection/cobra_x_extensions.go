package reflection

import (
	"reflect"
	"strings"
)

func reflectCobraXCommand(field reflect.StructField) (use, short, long string, ok bool) {
	cobraX := field.Tag.Get("cobra-x")
	if cobraX == "" {
		return "", "", "", false
	}
	use, attributes := parseCobraX(cobraX)

	short = attributes["help"]
	long = attributes["description"]

	return use, short, long, true
}

func reflectCobraXFlag(field reflect.StructField) (name, shorthand, usage, defaultValue string, ok bool) {
	cobraX := field.Tag.Get("cobra-x")
	if cobraX == "" {
		return "", "", "", "", false
	}
	nameExpr, attributes := parseCobraX(cobraX)
	if nameExpr != "" {
		name, shorthand = parseFlagNameExpression(nameExpr)
	}

	usage = attributes["help"]
	if usage == "" {
		usage = attributes["description"]
	}
	if usage == "" {
		usage = attributes["usage"]
	}

	defaultValue = attributes["default"]

	return name, shorthand, usage, defaultValue, true
}

func parseFlagNameExpression(s string) (string, string) {
	parts := strings.Split(s, "|")
	name := ""
	shorthand := ""
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "--") {
			name = strings.TrimPrefix(part, "--")
		} else if strings.HasPrefix(part, "-") {
			shorthand = strings.TrimPrefix(part, "-")
		} else {
			if name == "" {
				name = part
			}
		}
	}
	if name == "" && shorthand != "" {
		name = shorthand
	}
	return name, shorthand
}

func parseCobraX(tagValue string) (string, map[string]string) {
	attributes := make(map[string]string)

	var parts []string
	var current strings.Builder
	inQuotes := false
	var quoteChar rune

	for _, r := range tagValue {
		if (r == '\'' || r == '"') && (len(current.String()) == 0 || current.String()[len(current.String())-1] != '\\') {
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

	nameExpr := ""
	if len(parts) > 0 {
		nameExpr = strings.TrimSpace(parts[0])
		for i := 1; i < len(parts); i++ {
			part := strings.TrimSpace(parts[i])
			kv := strings.SplitN(part, "=", 2)
			if len(kv) == 2 {
				key := strings.TrimSpace(kv[0])
				val := strings.TrimSpace(kv[1])
				if len(val) >= 2 && ((val[0] == '\'' && val[len(val)-1] == '\'') || (val[0] == '"' && val[len(val)-1] == '"')) {
					val = val[1 : len(val)-1]
				}
				attributes[key] = val
			}
		}
	}
	return nameExpr, attributes
}
