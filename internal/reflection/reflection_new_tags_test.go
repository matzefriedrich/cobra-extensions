package reflection

import (
	"testing"

	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"github.com/stretchr/testify/assert"
)

type newTagsCommand struct {
	use   types.CommandName `cobra-x:"newcmd, help='new short', description='new long'"`
	Name  string            `cobra-x:"-n|--name, help='name usage'"`
	Force bool              `cobra-x:"--force, help='force usage'"`
}

type baseCommandWithTags struct {
	types.BaseCommand `cobra-x:"basecmd, help='base short', description='base long'"`
	Flag              string `cobra-x:"--flag, help='flag usage'"`
}

func Test_ReflectCommandDescriptor_with_new_tags(t *testing.T) {
	// Arrange
	reflector := NewCommandReflector[*newTagsCommand]()
	cmd := &newTagsCommand{}

	// Act
	descriptor := reflector.ReflectCommandDescriptor(cmd).(*commandDescriptor)

	// Assert
	assert.Equal(t, "newcmd", descriptor.use)
	assert.Equal(t, "new short", descriptor.short)
	assert.Equal(t, "new long", descriptor.long)

	assert.Len(t, descriptor.flags, 2)
	
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

func Test_ReflectCommandDescriptor_with_BaseCommand_tags(t *testing.T) {
	// Arrange
	reflector := NewCommandReflector[*baseCommandWithTags]()
	cmd := &baseCommandWithTags{}

	// Act
	descriptor := reflector.ReflectCommandDescriptor(cmd).(*commandDescriptor)

	// Assert
	assert.Equal(t, "basecmd", descriptor.use)
	assert.Equal(t, "base short", descriptor.short)
	assert.Equal(t, "base long", descriptor.long)
}

type unexportedUseCommand struct {
	_    types.CommandName `cobra-x:"unexported"`
	Flag string            `cobra-x:"--flag"`
}

func Test_ReflectCommandDescriptor_with_unexported_use_field(t *testing.T) {
	// Arrange
	reflector := NewCommandReflector[*unexportedUseCommand]()
	cmd := &unexportedUseCommand{}

	// Act
	descriptor := reflector.ReflectCommandDescriptor(cmd).(*commandDescriptor)

	// Assert
	assert.Equal(t, "unexported", descriptor.use)
}

type descriptionTagCommand struct {
	_ types.BaseCommand `cobra-x:"desccmd, help='short text', description='long text from description'"`
}

func Test_ReflectCommandDescriptor_with_description_tag(t *testing.T) {
	// Arrange
	reflector := NewCommandReflector[*descriptionTagCommand]()
	cmd := &descriptionTagCommand{}

	// Act
	descriptor := reflector.ReflectCommandDescriptor(cmd).(*commandDescriptor)

	// Assert
	assert.Equal(t, "desccmd", descriptor.use)
	assert.Equal(t, "short text", descriptor.short)
	assert.Equal(t, "long text from description", descriptor.long)
}
