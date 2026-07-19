package reflection

import (
	"reflect"
)

func reflectLegacyCommand(field reflect.StructField) (use, short, long string, ok bool) {
	use = field.Tag.Get("flag")
	short = field.Tag.Get("usage")

	if use == "" && short == "" {
		return "", "", "", false
	}
	return use, short, "", true
}

func reflectLegacyFlag(field reflect.StructField) (name, shorthand, usage, defaultValue string, ok bool) {
	usage = field.Tag.Get("usage")
	shorthand = field.Tag.Get("shorthand")
	name = field.Tag.Get("flag")
	defaultValue = field.Tag.Get("default")

	if name == "" && shorthand == "" && usage == "" && defaultValue == "" {
		return "", "", "", "", false
	}
	return name, shorthand, usage, defaultValue, true
}
