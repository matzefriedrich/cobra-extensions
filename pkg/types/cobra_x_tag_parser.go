package types

import (
	"reflect"
)

const (
	CobraXTagKey          = "cobra-x"
	CobraXShorthand       = "cobra-x-shorthand"
	HelpAttribute         = "help"
	DescriptionAttribute  = "description"
	UsageAttribute        = "usage"
	DefaultValueAttribute = "default"
	SettingKeyAttribute   = "setting-key"
)

// CobraXTagParser parses the cobra-x tags of a struct field into a name expression and an attributes map.
type CobraXTagParser interface {
	// ParseField parses the cobra-x tags of the given struct field.
	// It returns the flag/command name expression and the attribute map.
	ParseField(field reflect.StructField) (name string, attributes map[string]string)
}
