package types

import (
	"reflect"
	"strings"
)

// fluentAttributeKeys maps the parser-specific chain keys to the canonical attribute names. The name key is special.
var fluentAttributeKeys = map[string]string{
	"name":        "",
	"shorthand":   SettingKeyAttribute,
	"help":        HelpAttribute,
	"description": DescriptionAttribute,
	"usage":       UsageAttribute,
	"default":     DefaultValueAttribute,
}

var fluentAttributeKeyOrder = []string{"name", "shorthand", "help", "description", "usage", "default"}

// NewFluentTagParser Creates a parser for cobra-x tags written as a fluent builder chain,
// e.g. --name.shorthand(n).help('Name to greet, with commas').default(host:8080).
func NewFluentTagParser() CobraXTagParser {
	return &fluentTagParser{}
}

type fluentTagParser struct{}

// ParseField parses the fluent builder chain of the given struct field.
func (p *fluentTagParser) ParseField(field reflect.StructField) (string, map[string]string) {
	tagValue := strings.TrimSpace(field.Tag.Get(CobraXTagKey))
	if tagValue == "" {
		return "", nil
	}
	nameEnd := len(tagValue)
	attributes := make(map[string]string)
	var name string
	nameFromSegment := false
	for index := 0; index < len(tagValue); index++ {
		if tagValue[index] != '.' {
			continue
		}
		key, openParen, matched := matchFluentSegment(tagValue, index)
		if !matched {
			continue
		}
		if nameEnd == len(tagValue) {
			nameEnd = index
		}
		value, closeParen, valid := scanFluentValue(tagValue, openParen)
		if !valid {
			return "", nil
		}
		target := fluentAttributeKeys[key]
		if target == "" {
			name = value
			nameFromSegment = true
		} else {
			attributes[target] = value
		}
		index = closeParen
	}
	if !nameFromSegment {
		name = strings.TrimSpace(tagValue[:nameEnd])
	}
	if name == "" {
		return "", nil
	}
	return name, attributes
}

// matchFluentSegment detects a keyed segment at the given dot, returning its key and the position of its opening parenthesis.
func matchFluentSegment(tagValue string, dotIndex int) (string, int, bool) {
	for _, key := range fluentAttributeKeyOrder {
		openParen := dotIndex + 1 + len(key)
		if openParen < len(tagValue) &&
			tagValue[openParen] == '(' &&
			strings.HasPrefix(tagValue[dotIndex+1:], key) {
			return key, openParen, true
		}
	}
	return "", 0, false
}

// scanFluentValue scans the value between the parenthesis pair, tracking nested parentheses and quoted regions.
func scanFluentValue(tagValue string, openParen int) (string, int, bool) {
	depth := 1
	inQuotes := false
	var value strings.Builder
	for index := openParen + 1; index < len(tagValue); index++ {
		current := tagValue[index]
		if current == '\'' {
			inQuotes = !inQuotes
			_ = value.WriteByte(current)
			continue
		}
		if !inQuotes {
			switch current {
			case '(':
				depth++
			case ')':
				depth--
				if depth == 0 {
					unquoted, valid := unquoteTagScalar(value.String())
					return unquoted, index, valid
				}
			}
		}
		_ = value.WriteByte(current)
	}
	return "", 0, false
}
