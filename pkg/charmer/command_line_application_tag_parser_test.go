package charmer

import (
	"context"
	"testing"

	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"github.com/stretchr/testify/assert"
)

type standardTagHandler struct {
	types.BaseCommand `cobra-x:"standard-test" cobra-x-description:"Standard tag test command"`
	Name              string `cobra-x:"--name" cobra-x-shorthand:"n" cobra-x-usage:"Name to greet" cobra-x-default:"World"`
}

func (h *standardTagHandler) Execute(_ context.Context) {
}

type dslTagHandler struct {
	types.BaseCommand `cobra-x:"dsl-test"`
	Celebrant         string `cobra-x:"-c|--celebrant, help='Who to greet'"`
}

func (h *dslTagHandler) Execute(_ context.Context) {
}

func Test_CommandLineApplication_AddTypedCommand_applies_injected_standard_tag_parser(t *testing.T) {
	// Arrange
	handler := &standardTagHandler{}
	sut := NewCommandLineApplication("app", "app desc").
		WithTagParser(types.NewStandardTagParser()).
		AddTypedCommand(handler)
	sut.root.SetArgs([]string{"standard-test", "--name", "Ada"})

	// Act
	err := sut.Execute(t.Context())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "Ada", handler.Name)
}

func Test_CommandLineApplication_AddTypedCommand_uses_dsl_parser_by_default(t *testing.T) {
	// Arrange
	handler := &dslTagHandler{}
	sut := NewCommandLineApplication("app", "app desc").
		AddTypedCommand(handler)
	sut.root.SetArgs([]string{"dsl-test", "--celebrant", "Grace"})

	// Act
	err := sut.Execute(t.Context())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "Grace", handler.Celebrant)
}
