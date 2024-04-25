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

func (d *CommandDescriptor) Arguments() ArgumentsDescriptor {
	return d.arguments
}

func (d *CommandDescriptor) Key() string {
	return d.key
}

func (d *CommandDescriptor) LongDescriptionText() string {
	return d.long
}

func (d *CommandDescriptor) ShortDescriptionText() string {
	return d.short
}

func (d *CommandDescriptor) Use() string {
	return d.use
}

func (d *CommandDescriptor) Flags() []FlagDescriptor {
	return d.flags
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

func (d *CommandDescriptor) BindFlags(target *cobra.Command) {
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
