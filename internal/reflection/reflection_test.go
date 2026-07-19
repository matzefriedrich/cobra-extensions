package reflection

import (
	"testing"

	"github.com/matzefriedrich/cobra-extensions/pkg/types"

	"github.com/stretchr/testify/assert"
)

type simpleCommand struct {
	Name types.CommandName `flag:"simple" usage:"short desc"`
	Flag string            `flag:"flag" shorthand:"f" usage:"flag usage"`
}

type embeddedCommand struct {
	simpleCommand
	Extra string `flag:"extra" usage:"extra usage"`
}

type interfaceCommand struct {
	Args types.CommandArgs `minimum:"1"`
	Arg1 any
}

type fullArgsCommand struct {
	Command struct {
		Args types.CommandArgs
		Arg1 string
		Arg2 int64
		Arg3 bool
	}
}

func Test_commandReflector_reflect_command_descriptor_correctly_reflects_different_struct_types(t *testing.T) {
	tests := []struct {
		name string
		cmd  any
	}{
		{name: "simple command", cmd: &simpleCommand{}},
		{name: "non-pointer command", cmd: simpleCommand{}},
		{name: "embedded command", cmd: &embeddedCommand{}},
		{name: "interface argument", cmd: &interfaceCommand{}},
		{name: "full arguments", cmd: &fullArgsCommand{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			reflector := NewCommandReflector[any]()

			// Act
			descriptor := reflector.ReflectCommandDescriptor(tt.cmd)

			// Assert
			assert.NotNil(t, descriptor)
		})
	}
}
