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
