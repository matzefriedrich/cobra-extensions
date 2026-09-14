package types

import (
	"reflect"
	"strings"
)

type standardTagParser struct{}

// NewStandardTagParser Creates a new CobraXTagParser that interprets conventional
// space-separated struct tags: cobra-x:"name" cobra-x-shorthand:"x" cobra-x-usage:"..."
func NewStandardTagParser() CobraXTagParser {
	return &standardTagParser{}
}

func (p *standardTagParser) ParseField(field reflect.StructField) (string, map[string]string) {
	name := standardFieldName(field)
	if name == "" {
		return "", nil
	}
	attributes := make(map[string]string)
	collectStandardAttributes(field, attributes)
	return name, attributes
}

func standardFieldName(field reflect.StructField) string {
	name := strings.TrimSpace(field.Tag.Get(CobraXTagKey))
	if shorthand := strings.TrimSpace(field.Tag.Get(CobraXShorthand)); shorthand != "" {
		name = "-" + shorthand + "|" + name
	}
	return name
}

func collectStandardAttributes(field reflect.StructField, attributes map[string]string) {
	for tagName, attribute := range standardAttributeTags {
		collectAttribute(field, tagName, attribute, attributes)
	}
}

var standardAttributeTags = map[string]string{
	"cobra-x-help":        HelpAttribute,
	"cobra-x-description": DescriptionAttribute,
	"cobra-x-usage":       UsageAttribute,
	"cobra-x-default":     DefaultValueAttribute,
	"cobra-x-setting-key": SettingKeyAttribute,
}

func collectAttribute(field reflect.StructField, tagKey string, attributeKey string, attributes map[string]string) {
	value := strings.TrimSpace(field.Tag.Get(tagKey))
	if value != "" {
		attributes[attributeKey] = value
	}
}
