package reflection

import (
	"errors"
	"reflect"

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
	SettingKey   string
}

const (
	cobraXTag                     = "cobra-x"
	ErrorCobraXTagNotFound        = "cobra-x tag not found"
	ErrorCobraXLegacyTagsNotFound = "cobra-x legacy tags not found"
	cobraXHelpTag                 = "help"
	cobraXDescriptionTag          = "description"
	cobraXUsageTag                = "usage"
	cobraXDefaultValueTag         = "default"
	cobraXSettingKeyTag           = "setting-key"
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
	cobraXTagValue := field.Tag.Get(cobraXTag)
	if cobraXTagValue == "" {
		return nil, types.NewCobraXError(ErrorCobraXTagNotFound)
	}
	var name, shorthand, usage, defaultValue string
	nameExpr, attributes := parseCobraX(cobraXTagValue)
	if nameExpr != "" {
		name, shorthand = parseFlagNameExpression(nameExpr)
	}

	usage = resolveFlagUsage(attributes)
	defaultValue = attributes[cobraXDefaultValueTag]
	settingKey := attributes[cobraXSettingKeyTag]

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
	usage := attributes[cobraXHelpTag]
	if usage == "" {
		usage = attributes[cobraXDescriptionTag]
	}
	if usage == "" {
		usage = attributes[cobraXUsageTag]
	}
	return usage
}
