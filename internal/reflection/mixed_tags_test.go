package reflection

import (
	"testing"

	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"github.com/stretchr/testify/assert"
)

type mixedTagsCommand struct {
	_    types.BaseCommand `flag:"mixed" usage:"short description" json:"-"`
	Name string            `flag:"name" shorthand:"n" usage:"name usage" json:"name"`
}

func Test_ReflectCommandDescriptor_with_mixed_tags(t *testing.T) {
	// Arrange
	reflector := NewCommandReflector[*mixedTagsCommand]()
	cmd := &mixedTagsCommand{Name: "test"}

	// Act
	descriptor := reflector.ReflectCommandDescriptor(cmd).(*commandDescriptor)

	// Assert
	assert.Equal(t, "mixed", descriptor.use)
	assert.Equal(t, "short description", descriptor.short)

	var nameFlag *FlagDescriptor
	for _, f := range descriptor.flags {
		if f.name == "name" {
			nameFlag = &f
			break
		}
	}
	assert.NotNil(t, nameFlag)
	assert.Equal(t, "n", nameFlag.shorthand)
	assert.Equal(t, "name usage", nameFlag.usage)
}

type mixedCobraXTagsCommand struct {
	types.BaseCommand `cobra-x:"mixed, help='short description'" json:"-"`
	Name              string `cobra-x:"-n|--name, help='name usage'" json:"name"`
}

func Test_ReflectCommandDescriptor_with_mixed_cobra_x_tags(t *testing.T) {
	// Arrange
	reflector := NewCommandReflector[*mixedCobraXTagsCommand]()
	cmd := &mixedCobraXTagsCommand{Name: "test"}

	// Act
	descriptor := reflector.ReflectCommandDescriptor(cmd).(*commandDescriptor)

	// Assert
	assert.Equal(t, "mixed", descriptor.use)
	assert.Equal(t, "short description", descriptor.short)

	var nameFlag *FlagDescriptor
	for _, f := range descriptor.flags {
		if f.name == "name" {
			nameFlag = &f
			break
		}
	}
	assert.NotNil(t, nameFlag)
	assert.Equal(t, "n", nameFlag.shorthand)
	assert.Equal(t, "name usage", nameFlag.usage)
}

type cobraXStrictIsolationCommand struct {
	types.BaseCommand `cobra-x:"new" usage:"should be ignored"`
	Name              string `cobra-x:"--name" usage:"should be ignored"`
}

func Test_ReflectCommandDescriptor_with_cobra_x_strict_isolation(t *testing.T) {
	// Arrange
	reflector := NewCommandReflector[*cobraXStrictIsolationCommand]()
	cmd := &cobraXStrictIsolationCommand{}

	// Act
	descriptor := reflector.ReflectCommandDescriptor(cmd).(*commandDescriptor)

	// Assert
	assert.Equal(t, "new", descriptor.use)
	assert.Equal(t, "", descriptor.short) // Legacy usage:"should be ignored" must be ignored

	var nameFlag *FlagDescriptor
	for _, f := range descriptor.flags {
		if f.name == "name" {
			nameFlag = &f
			break
		}
	}
	assert.NotNil(t, nameFlag)
	assert.Equal(t, "", nameFlag.usage) // Legacy usage:"should be ignored" must be ignored
}
