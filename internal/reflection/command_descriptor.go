package reflection

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"github.com/spf13/cobra"
)

// CommandDescriptor represents the metadata and configuration for a command, including its use, descriptions, flags, and arguments.
type commandDescriptor struct {
	key       string
	use       string
	short     string
	long      string
	flags     []FlagDescriptor
	arguments types.ArgumentsDescriptor
}

var _ types.CommandDescriptor = (*commandDescriptor)(nil)

// Key returns the key string associated with the CommandDescriptor.
func (d *commandDescriptor) Key() string {
	return d.key
}

// NewCommandDescriptor creates a new CommandDescriptor with specified use, short and long descriptions, flags, and arguments.
func NewCommandDescriptor(use string, short string, long string, flags []FlagDescriptor, arguments types.ArgumentsDescriptor) types.CommandDescriptor {
	key := makeCommandKey(use)
	return &commandDescriptor{
		key:       key,
		use:       use,
		short:     short,
		long:      long,
		flags:     flags,
		arguments: arguments,
	}
}

func makeCommandKey(use string) string {
	id := globalUidSequence.Next()
	idString := strconv.FormatUint(id, 10)
	key := fmt.Sprintf("%s-%s", use, idString)
	return key
}

// BindFlags Binds the reflected flags configuration to the given *cobra.Command object.
func (d *commandDescriptor) BindFlags(target *cobra.Command) {
	if target == nil {
		return
	}
	for _, f := range d.flags {
		targetFlags := target.Flags()
		name := f.name
		usage := f.usage
		switch f.kind {
		case reflect.String:
			targetFlags.String(name, f.AsString(), usage)
		case reflect.Int, reflect.Int64:
			targetFlags.Int64(name, f.AsInt64(), usage)
		case reflect.Bool:
			targetFlags.Bool(name, f.AsBool(), usage)
		}
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
	flags := target.Flags()
	for _, f := range d.flags {
		flagName := f.name
		switch f.kind {
		case reflect.String:
			value, _ := flags.GetString(flagName)
			_ = f.SetValue(value)
		case reflect.Int, reflect.Int64:
			value, _ := flags.GetInt64(flagName)
			_ = f.SetValue(value)
		case reflect.Bool:
			value, _ := flags.GetBool(flagName)
			_ = f.SetValue(value)
		}
	}
}
