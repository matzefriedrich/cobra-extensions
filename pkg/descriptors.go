package pkg

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type FlagDescriptor struct {
	name  string
	kind  reflect.Kind
	value reflect.Value
	usage string
}

func (d *FlagDescriptor) AsString() string {
	return d.value.String()
}

func (d *FlagDescriptor) AsInt64() int64 {
	return d.value.Int()
}

func (d *FlagDescriptor) AsBool() bool {
	return d.value.Bool()
}

func NewFlagDescriptor(name string, usage string, t reflect.Kind, v reflect.Value) FlagDescriptor {
	return FlagDescriptor{
		name:  name,
		usage: usage,
		kind:  t,
		value: v,
	}
}

func (d *FlagDescriptor) SetValue(value interface{}) error {
	switch d.kind {
	case reflect.String:
		s, ok := value.(string)
		if ok {
			d.value.SetString(s)
			return nil
		}
		return errors.New("the specified value does not match the flag type")
	case reflect.Int, reflect.Int64:
		n, ok := value.(int64)
		if ok {
			d.value.SetInt(n)
			return nil
		}
		return errors.New("the specified value does not match the flag type")
	case reflect.Bool:
		b, ok := value.(bool)
		if ok {
			d.value.SetBool(b)
			return nil
		}
		return errors.New("the specified value does not match the flag type")
	}
	return errors.New("type not supported. flag must be of type string, int, or bool")
}

type CommandDescriptor struct {
	key   string
	use   string
	short string
	long  string
	flags []FlagDescriptor
}

func NewCommandDescriptor(use string, short string, long string, flags []FlagDescriptor) CommandDescriptor {
	key := makeCommandKey(use)
	return CommandDescriptor{
		key:   key,
		use:   use,
		short: short,
		long:  long,
		flags: flags,
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
