package reflection

import (
	"errors"
	"reflect"
	"strings"

	"github.com/matzefriedrich/cobra-extensions/pkg/types"
)

type CobraXTag struct {
	Attributes map[string]string
}

type CobraXCommandTag struct {
	CobraXTag
	Use         string
	Help        string
	Description string
}

type CobraXFlagTag struct {
	CobraXTag
	Name         string
	Shorthand    string
	Usage        string
	DefaultValue string
}

const (
	cobraXTag                     = "cobra-x"
	ErrorCobraXTagNotFound        = "cobra-x tag not found"
	ErrorCobraXLegacyTagsNotFound = "cobra-x legacy tags not found"
	cobraXHelpTag                 = "help"
	cobraXDescriptionTag          = "description"
	cobraXUsageTag                = "usage"
	cobraXDefaultValueTag         = "default"
)

var (
	ErrCobraXCommandNotFound    = errors.New(ErrorCobraXTagNotFound)
	ErrCobraXLegacyTagsNotFound = errors.New(ErrorCobraXLegacyTagsNotFound)
)

func reflectCobraXCommand(field reflect.StructField) (*CobraXCommandTag, error) {
	cobraX := field.Tag.Get(cobraXTag)
	if cobraX == "" {
		return nil, types.NewCobraXError(ErrorCobraXTagNotFound)
	}
	use, attributes := parseCobraX(cobraX)
	help := attributes[cobraXHelpTag]
	description := attributes[cobraXDescriptionTag]
	return &CobraXCommandTag{
		CobraXTag:   CobraXTag{Attributes: attributes},
		Use:         use,
		Help:        help,
		Description: description,
	}, nil
}

func reflectCobraXFlag(field reflect.StructField) (*CobraXFlagTag, error) {
	cobraX := field.Tag.Get(cobraXTag)
	if cobraX == "" {
		return nil, types.NewCobraXError(ErrorCobraXTagNotFound)
	}
	var name, shorthand, usageOrDescription, defaultValue string
	nameExpr, attributes := parseCobraX(cobraX)
	if nameExpr != "" {
		name, shorthand = parseFlagNameExpression(nameExpr)
	}

	usageOrDescription = attributes[cobraXHelpTag]
	if usageOrDescription == "" {
		usageOrDescription = attributes[cobraXDescriptionTag]
	}
	if usageOrDescription == "" {
		usageOrDescription = attributes[cobraXUsageTag]
	}

	defaultValue = attributes[cobraXDefaultValueTag]

	return &CobraXFlagTag{
		CobraXTag:    CobraXTag{Attributes: attributes},
		Name:         name,
		Shorthand:    shorthand,
		Usage:        usageOrDescription,
		DefaultValue: defaultValue,
	}, nil
}

func parseFlagNameExpression(s string) (name string, shortHand string) {
	parts := strings.Split(s, "|")
	name = ""
	shortHand = ""
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if after, ok := strings.CutPrefix(part, "--"); ok {
			name = after
		} else if after0, ok0 := strings.CutPrefix(part, "-"); ok0 {
			shortHand = after0
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
