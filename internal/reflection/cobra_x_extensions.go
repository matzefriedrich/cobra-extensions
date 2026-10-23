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
	Aliases     []string
}

type CobraXFlagTag struct {
	CobraXTag
	Name         string
	Shorthand    string
	Usage        string
	DefaultValue string
	SettingKey   string
}

const (
	ErrorCobraXTagNotFound        = "cobra-x tag not found"
	ErrorCobraXLegacyTagsNotFound = "cobra-x legacy tags not found"
)

var (
	ErrCobraXCommandNotFound    = errors.New(ErrorCobraXTagNotFound)
	ErrCobraXLegacyTagsNotFound = errors.New(ErrorCobraXLegacyTagsNotFound)
)

func reflectCobraXCommand(field reflect.StructField, tagParser types.CobraXTagParser) (*CobraXCommandTag, error) {
	if field.Tag.Get(types.CobraXTagKey) == "" {
		return nil, types.NewCobraXError(ErrorCobraXTagNotFound)
	}
	use, attributes := tagParser.ParseField(field)
	help := attributes[types.HelpAttribute]
	description := attributes[types.DescriptionAttribute]
	aliases := splitCommandAliases(attributes[types.AliasAttribute])
	return &CobraXCommandTag{
		CobraXTag:   CobraXTag{Attributes: attributes},
		Use:         use,
		Help:        help,
		Description: description,
		Aliases:     aliases,
	}, nil
}

func splitCommandAliases(value string) []string {
	if value == "" {
		return nil
	}
	aliases := strings.Split(value, ",")
	for i := range aliases {
		aliases[i] = strings.TrimSpace(aliases[i])
	}
	return aliases
}

func reflectCobraXFlag(field reflect.StructField, tagParser types.CobraXTagParser) (*CobraXFlagTag, error) {
	if field.Tag.Get(types.CobraXTagKey) == "" {
		return nil, types.NewCobraXError(ErrorCobraXTagNotFound)
	}
	var name, shorthand, usage, defaultValue string
	nameExpr, attributes := tagParser.ParseField(field)
	if nameExpr != "" {
		name, shorthand = parseFlagNameExpression(nameExpr)
	}

	usage = resolveFlagUsage(attributes)
	defaultValue = attributes[types.DefaultValueAttribute]
	settingKey := attributes[types.SettingKeyAttribute]

	return &CobraXFlagTag{
		CobraXTag:    CobraXTag{Attributes: attributes},
		Name:         name,
		Shorthand:    shorthand,
		Usage:        usage,
		DefaultValue: defaultValue,
		SettingKey:   settingKey,
	}, nil
}

func resolveFlagUsage(attributes map[string]string) string {
	usage := attributes[types.HelpAttribute]
	if usage == "" {
		usage = attributes[types.DescriptionAttribute]
	}
	if usage == "" {
		usage = attributes[types.UsageAttribute]
	}
	return usage
}
