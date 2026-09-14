package types

import (
	"reflect"
	"strings"
)

// yamlFlowAttributeKeys maps the parser-specific mapping keys to the canonical attribute names. The name key is special.
var yamlFlowAttributeKeys = map[string]string{
	"name":        "",
	"shorthand":   SettingKeyAttribute,
	"help":        HelpAttribute,
	"description": DescriptionAttribute,
	"usage":       UsageAttribute,
	"default":     DefaultValueAttribute,
}

// NewYamlFlowTagParser Creates a parser for cobra-x tags written as a flow-style YAML mapping,
// e.g. {name: '--name', shorthand: 'n', help: 'Prints a greeting.', default: 'World'}.
func NewYamlFlowTagParser() CobraXTagParser {
	return &yamlFlowTagParser{}
}

type yamlFlowTagParser struct{}

// ParseField parses the flow-style YAML mapping of the given struct field.
func (p *yamlFlowTagParser) ParseField(field reflect.StructField) (string, map[string]string) {
	mapping := strings.TrimSpace(field.Tag.Get(CobraXTagKey))
	if mapping == "" {
		return "", nil
	}
	if strings.HasPrefix(mapping, "{") && strings.HasSuffix(mapping, "}") {
		mapping = strings.TrimSpace(mapping[1 : len(mapping)-1])
	}
	attributes := make(map[string]string)
	var name string
	for _, entry := range splitYamlFlowEntries(mapping) {
		key, value, valid := splitYamlFlowEntry(entry)
		if !valid {
			return "", nil
		}
		target, known := yamlFlowAttributeKeys[key]
		if !known {
			continue
		}
		if key == "name" {
			name = value
		} else {
			attributes[target] = value
		}
	}
	if name == "" {
		return "", nil
	}
	return name, attributes
}

// splitYamlFlowEntries splits a flow mapping into its entries at the commas that appear outside single quotes.
func splitYamlFlowEntries(mapping string) []string {
	entries := make([]string, 0, 4)
	var entry strings.Builder
	inQuotes := false
	flush := func() {
		trimmed := strings.TrimSpace(entry.String())
		if trimmed != "" {
			entries = append(entries, trimmed)
		}
		entry.Reset()
	}
	for index := 0; index < len(mapping); index++ {
		current := mapping[index]
		switch current {
		case '\'':
			inQuotes = !inQuotes
			_ = entry.WriteByte(current)
		case ',':
			if inQuotes {
				_ = entry.WriteByte(current)
			} else {
				flush()
			}
		default:
			_ = entry.WriteByte(current)
		}
	}
	flush()
	return entries
}

// splitYamlFlowEntry splits a single entry at the first colon into its key and value parts.
func splitYamlFlowEntry(entry string) (string, string, bool) {
	separator := strings.Index(entry, ":")
	if separator < 0 {
		return "", "", false
	}
	key := strings.TrimSpace(entry[:separator])
	value, valid := unquoteTagScalar(strings.TrimSpace(entry[separator+1:]))
	return key, value, valid
}

// unquoteTagScalar removes one pair of surrounding single quotes, decoding doubled quotes as a literal apostrophe.
func unquoteTagScalar(value string) (string, bool) {
	if value == "" {
		return "", true
	}
	if value[0] == '\'' {
		if len(value) < 2 || value[len(value)-1] != '\'' {
			return "", false
		}
		return strings.ReplaceAll(value[1:len(value)-1], "''", "'"), true
	}
	return value, true
}
