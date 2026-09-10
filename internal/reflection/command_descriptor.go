package reflection

import (
	"reflect"

	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// CommandDescriptor represents the metadata and configuration for a command, including its use, descriptions, flags, and arguments.
type commandDescriptor struct {
	use                  string
	short                string
	long                 string
	flags                []FlagDescriptor
	arguments            types.ArgumentsDescriptor
	defaultValueProvider types.DefaultValueProvider
}

var _ types.CommandDescriptor = (*commandDescriptor)(nil)

// NewCommandDescriptor creates a new CommandDescriptor with specified use, short and long descriptions, flags, and arguments.
func NewCommandDescriptor(use string, short string, long string, flags []FlagDescriptor, arguments types.ArgumentsDescriptor) types.CommandDescriptor {
	return &commandDescriptor{
		use:       use,
		short:     short,
		long:      long,
		flags:     flags,
		arguments: arguments,
	}
}

// BindFlags Binds the reflected flags configuration to the given *cobra.Command object.
func (d *commandDescriptor) BindFlags(target *cobra.Command) {
	if target == nil {
		return
	}
	targetFlags := target.Flags()
	for _, flag := range d.flags {
		bindFlag(targetFlags, flag)
	}
}

func bindFlag(targetFlags *pflag.FlagSet, flag FlagDescriptor) {
	flagName := flag.name
	flagUsage := flag.usage
	flagShorthand := flag.shorthand
	switch flag.kind {
	case reflect.String:
		bindStringFlag(targetFlags, flagName, flagShorthand, flag.AsString(), flagUsage)
	case reflect.Int, reflect.Int64:
		bindInt64Flag(targetFlags, flagName, flagShorthand, flag.AsInt64(), flagUsage)
	case reflect.Bool:
		bindBoolFlag(targetFlags, flagName, flagShorthand, flag.AsBool(), flagUsage)
	case reflect.Slice:
		bindSliceFlag(targetFlags, flag)
	}
}

func bindStringFlag(targetFlags *pflag.FlagSet, flagName string, flagShorthand string, defaultValue string, flagUsage string) {
	if flagShorthand != "" {
		targetFlags.StringP(flagName, flagShorthand, defaultValue, flagUsage)
	} else {
		targetFlags.String(flagName, defaultValue, flagUsage)
	}
}

func bindInt64Flag(targetFlags *pflag.FlagSet, flagName string, flagShorthand string, defaultValue int64, flagUsage string) {
	if flagShorthand != "" {
		targetFlags.Int64P(flagName, flagShorthand, defaultValue, flagUsage)
	} else {
		targetFlags.Int64(flagName, defaultValue, flagUsage)
	}
}

func bindBoolFlag(targetFlags *pflag.FlagSet, flagName string, flagShorthand string, defaultValue bool, flagUsage string) {
	if flagShorthand != "" {
		targetFlags.BoolP(flagName, flagShorthand, defaultValue, flagUsage)
	} else {
		targetFlags.Bool(flagName, defaultValue, flagUsage)
	}
}

func bindSliceFlag(targetFlags *pflag.FlagSet, flag FlagDescriptor) {
	flagName := flag.name
	flagUsage := flag.usage
	flagShorthand := flag.shorthand
	switch flag.elementKind {
	case reflect.String:
		bindStringSliceFlag(targetFlags, flagName, flagShorthand, flag.value.Interface().([]string), flagUsage)
	case reflect.Int:
		bindIntSliceFlag(targetFlags, flagName, flagShorthand, flag.value.Interface().([]int), flagUsage)
	case reflect.Int64:
		bindInt64SliceFlag(targetFlags, flagName, flagShorthand, flag.value.Interface().([]int64), flagUsage)
	case reflect.Bool:
		bindBoolSliceFlag(targetFlags, flagName, flagShorthand, flag.value.Interface().([]bool), flagUsage)
	}
}

func bindStringSliceFlag(targetFlags *pflag.FlagSet, flagName string, flagShorthand string, defaultValue []string, flagUsage string) {
	if flagShorthand != "" {
		targetFlags.StringSliceP(flagName, flagShorthand, defaultValue, flagUsage)
	} else {
		targetFlags.StringSlice(flagName, defaultValue, flagUsage)
	}
}

func bindIntSliceFlag(targetFlags *pflag.FlagSet, flagName string, flagShorthand string, defaultValue []int, flagUsage string) {
	if flagShorthand != "" {
		targetFlags.IntSliceP(flagName, flagShorthand, defaultValue, flagUsage)
	} else {
		targetFlags.IntSlice(flagName, defaultValue, flagUsage)
	}
}

func bindInt64SliceFlag(targetFlags *pflag.FlagSet, flagName string, flagShorthand string, defaultValue []int64, flagUsage string) {
	if flagShorthand != "" {
		targetFlags.Int64SliceP(flagName, flagShorthand, defaultValue, flagUsage)
	} else {
		targetFlags.Int64Slice(flagName, defaultValue, flagUsage)
	}
}

func bindBoolSliceFlag(targetFlags *pflag.FlagSet, flagName string, flagShorthand string, defaultValue []bool, flagUsage string) {
	if flagShorthand != "" {
		targetFlags.BoolSliceP(flagName, flagShorthand, defaultValue, flagUsage)
	} else {
		targetFlags.BoolSlice(flagName, defaultValue, flagUsage)
	}
}

// BindArguments Binds the reflected arguments configuration to the given *cobra.Command object.
func (d *commandDescriptor) BindArguments(target *cobra.Command) {
	if target == nil {
		return
	}

	target.Use = d.use
	target.Short = d.short
	target.Long = d.long

	d.arguments.BindArguments(target)
}

// UnmarshalArgumentValues deserializes the command argument values from a list of strings and binds them to the corresponding fields.
func (d *commandDescriptor) UnmarshalArgumentValues(args ...string) {
	d.arguments.BindArgumentValues(args...)
}

// UnmarshalFlagValues populates the CommandDescriptor's flags from the provided *cobra.Command object.
func (d *commandDescriptor) UnmarshalFlagValues(target *cobra.Command) {
	targetFlags := target.Flags()
	for _, flag := range d.flags {
		if d.applyProvidedDefaultValue(targetFlags, flag) {
			continue
		}
		readFlagValue(targetFlags, flag)
	}
}

func (d *commandDescriptor) applyProvidedDefaultValue(targetFlags *pflag.FlagSet, flag FlagDescriptor) bool {
	flagName := flag.name
	if targetFlags.Changed(flagName) || flag.settingKey == "" || d.defaultValueProvider == nil {
		return false
	}
	value, err := d.defaultValueProvider.GetValue(flag.settingKey)
	if err != nil || value == "" {
		return false
	}
	_ = flag.SetValueFromText(value)
	return true
}

func readFlagValue(targetFlags *pflag.FlagSet, flag FlagDescriptor) {
	flagName := flag.name
	switch flag.kind {
	case reflect.String:
		value, _ := targetFlags.GetString(flagName)
		_ = flag.SetValue(value)
	case reflect.Int, reflect.Int64:
		value, _ := targetFlags.GetInt64(flagName)
		_ = flag.SetValue(value)
	case reflect.Bool:
		value, _ := targetFlags.GetBool(flagName)
		_ = flag.SetValue(value)
	case reflect.Slice:
		readSliceFlagValue(targetFlags, flag)
	}
}

func readSliceFlagValue(targetFlags *pflag.FlagSet, flag FlagDescriptor) {
	flagName := flag.name
	switch flag.elementKind {
	case reflect.String:
		value, _ := targetFlags.GetStringSlice(flagName)
		_ = flag.SetValue(value)
	case reflect.Int:
		value, _ := targetFlags.GetIntSlice(flagName)
		_ = flag.SetValue(value)
	case reflect.Int64:
		value, _ := targetFlags.GetInt64Slice(flagName)
		_ = flag.SetValue(value)
	case reflect.Bool:
		value, _ := targetFlags.GetBoolSlice(flagName)
		_ = flag.SetValue(value)
	}
}

func (d *commandDescriptor) SetDefaultValueProvider(provider types.DefaultValueProvider) {
	d.defaultValueProvider = provider
}
