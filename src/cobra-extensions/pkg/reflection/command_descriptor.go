package reflection

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"reflect"
)

type CommandDescriptor struct {
	key       string
	use       string
	short     string
	long      string
	flags     []FlagDescriptor
	arguments ArgumentsDescriptor
}

func (d *CommandDescriptor) Key() string {
	return d.key
}

func NewCommandDescriptor(use string, short string, long string, flags []FlagDescriptor, arguments ArgumentsDescriptor) CommandDescriptor {
	key := makeCommandKey(use)
	return CommandDescriptor{
		key:       key,
		use:       use,
		short:     short,
		long:      long,
		flags:     flags,
		arguments: arguments,
	}
}

func makeCommandKey(use string) string {
	id := uuid.New()
	idString := id.String()[0:8]
	key := fmt.Sprintf("%s-%s", use, idString)
	return key
}

// BindFlags Binds the reflected flags configuration to the given *cobra.Command object.
func (d *CommandDescriptor) BindFlags(target *cobra.Command) {
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
func (d *CommandDescriptor) BindArguments(target *cobra.Command) {
	if target == nil {
		return
	}

	target.Use = d.use
	target.Short = d.short
	target.Long = d.long

	d.arguments.BindArguments(target)
}

func (d *CommandDescriptor) UnmarshalArgumentValues(args ...string) {
	d.arguments.BindArgumentValues(args...)
}

func (d *CommandDescriptor) UnmarshalFlagValues(target *cobra.Command) {
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
