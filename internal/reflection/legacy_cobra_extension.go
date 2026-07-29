package reflection

import (
	"errors"
	"reflect"

	"github.com/matzefriedrich/cobra-extensions/pkg/types"
)

func reflectLegacyCommand(field reflect.StructField) (*CobraXCommandTag, error) {
	use := field.Tag.Get("flag")
	short := field.Tag.Get("usage")
	if use == "" && short == "" {
		return nil, errors.New(ErrorCobraXLegacyTagsNotFound)
	}
	return &CobraXCommandTag{
		Use:         use,
		Help:        short,
		Description: short,
	}, nil
}

func reflectLegacyFlag(field reflect.StructField) (*CobraXFlagTag, error) {
	usage := field.Tag.Get("usage")
	shorthand := field.Tag.Get("shorthand")
	name := field.Tag.Get("flag")
	defaultValue := field.Tag.Get("default")

	if name == "" && shorthand == "" && usage == "" && defaultValue == "" {
		return nil, types.NewCobraXError(ErrorCobraXLegacyTagsNotFound)
	}
	return &CobraXFlagTag{
		Name:         name,
		Shorthand:    shorthand,
		Usage:        usage,
		DefaultValue: defaultValue,
	}, nil
}
