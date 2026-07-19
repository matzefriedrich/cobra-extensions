package reflection

import (
	"testing"

	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"github.com/stretchr/testify/assert"
)

type cobraXFlagCommand struct {
	types.BaseCommand `cobra-x:"greet"`
	Name              string `cobra-x:"-n|--name, help='The name to greet'"`
	Force             bool   `cobra-x:"-f, help='Force greeting'"`
	Quiet             bool   `cobra-x:"--quiet, description='Quiet mode'"`
	DefaultValue      string `cobra-x:"--msg, default='Hello'"`
}

func Test_ReflectCommandDescriptor_with_cobra_x_tags(t *testing.T) {
	// Arrange
	reflector := NewCommandReflector[*cobraXFlagCommand]()
	cmd := &cobraXFlagCommand{DefaultValue: "Hello"}

	// Act
	descriptor := reflector.ReflectCommandDescriptor(cmd).(*commandDescriptor)

	// Assert
	assert.Equal(t, "greet", descriptor.use)

	var nameFlag, forceFlag, quietFlag, defaultFlag *FlagDescriptor
	for i := range descriptor.flags {
		f := &descriptor.flags[i]
		switch f.name {
		case "name":
			nameFlag = f
		case "f":
			forceFlag = f
		case "quiet":
			quietFlag = f
		case "msg":
			defaultFlag = f
		}
	}

	assert.NotNil(t, nameFlag)
	assert.Equal(t, "n", nameFlag.shorthand)
	assert.Equal(t, "The name to greet", nameFlag.usage)

	assert.NotNil(t, forceFlag)
	assert.Equal(t, "f", forceFlag.shorthand)
	assert.Equal(t, "Force greeting", forceFlag.usage)

	assert.NotNil(t, quietFlag)
	assert.Equal(t, "quiet", quietFlag.name)
	assert.Equal(t, "Quiet mode", quietFlag.usage)

	assert.NotNil(t, defaultFlag)
	assert.Equal(t, "msg", defaultFlag.name)
	assert.Equal(t, "Hello", defaultFlag.AsString())
}

func Test_parseCobraX_handles_commas_in_quotes(t *testing.T) {
	// Arrange
	tag := "-n|--name, help='Hello, World', default='Value'"

	// Act
	nameExpr, attrs := parseCobraX(tag)

	// Assert
	assert.Equal(t, "-n|--name", nameExpr)
	assert.Equal(t, "Hello, World", attrs["help"])
	assert.Equal(t, "Value", attrs["default"])
}

type cobraXFallbackCommand struct {
	types.BaseCommand `cobra-x:"legacy"`
}

type cobraXCommandWithHelp struct {
	types.BaseCommand `cobra-x:"run, help='Start the service', description='This command starts the background service.'"`
}

func Test_ReflectCommandDescriptor_with_cobra_fallback_tags(t *testing.T) {
	// Arrange
	reflector := NewCommandReflector[*cobraXFallbackCommand]()
	cmd := &cobraXFallbackCommand{}

	// Act
	descriptor := reflector.ReflectCommandDescriptor(cmd).(*commandDescriptor)

	// Assert
	assert.Equal(t, "legacy", descriptor.use)
}

func Test_ReflectCommandDescriptor_with_cobra_x_command_help(t *testing.T) {
	// Arrange
	reflector := NewCommandReflector[*cobraXCommandWithHelp]()
	cmd := &cobraXCommandWithHelp{}

	// Act
	descriptor := reflector.ReflectCommandDescriptor(cmd).(*commandDescriptor)

	// Assert
	assert.Equal(t, "run", descriptor.use)
	assert.Equal(t, "Start the service", descriptor.short)
	assert.Equal(t, "This command starts the background service.", descriptor.long)
}
